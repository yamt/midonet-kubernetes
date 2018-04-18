// Copyright 2015 Tigera Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	cnitypes "github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/yamt/midonet-kubernetes/pkg/cni/midonet"
	"github.com/yamt/midonet-kubernetes/pkg/cni/types"
	"github.com/yamt/midonet-kubernetes/pkg/cni/utils"
	"github.com/yamt/midonet-kubernetes/pkg/converter/pod"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// CmdAddK8s performs the "ADD" operation on a kubernetes pod
// Having kubernetes code in its own file avoids polluting the mainline code. It's expected that the kubernetes case will
// more special casing than the mainline code.
func CmdAddK8s(args *skel.CmdArgs, conf types.NetConf, epIDs utils.WEPIdentifiers) (*current.Result, error) {
	var err error
	var result *current.Result

	utils.ConfigureLogging(conf.LogLevel)

	logger := logrus.WithFields(logrus.Fields{
		"ContainerID":      epIDs.ContainerID,
		"Pod":              epIDs.Pod,
		"Namespace":        epIDs.Namespace,
	})

	logger.Info("Extracted identifiers for CmdAddK8s")

	client, err := newK8sClient(conf, logger)
	if err != nil {
		return nil, err
	}
	logger.WithField("client", client).Debug("Created Kubernetes client")

	// Replace the actual value in the args.StdinData as that's what's passed to the IPAM plugin.
	fmt.Fprint(os.Stderr, "Calico CNI fetching podCidr from Kubernetes\n")
	var stdinData map[string]interface{}
	if err := json.Unmarshal(args.StdinData, &stdinData); err != nil {
		return nil, err
	}
	podCidr, err := getPodCidr(client, conf, epIDs.Namespace, epIDs.Pod)
	if err != nil {
		logger.Info("Failed to getPodCidr")
		return nil, err
	}
	logger.WithField("podCidr", podCidr).Info("Fetched podCidr")
	stdinData["ipam"].(map[string]interface{})["subnet"] = podCidr
	fmt.Fprintf(os.Stderr, "Calico CNI passing podCidr to host-local IPAM: %s\n", podCidr)
	args.StdinData, err = json.Marshal(stdinData)
	if err != nil {
		return nil, err
	}
	logger.WithField("stdin", string(args.StdinData)).Debug("Updated stdin data")

	logger.Debugf("Calling IPAM plugin %s", conf.IPAM.Type)
	ipamResult, err := ipam.ExecAdd(conf.IPAM.Type, args.StdinData)
	if err != nil {
		return nil, err
	}
	logger.Debugf("IPAM plugin returned: %+v", ipamResult)

	// Convert IPAM result into current Result.
	// IPAM result has a bunch of fields that are optional for an IPAM plugin
	// but required for a CNI plugin, so this is to populate those fields.
	// See CNI Spec doc for more details.
	result, err = current.NewResultFromResult(ipamResult)
	if err != nil {
		utils.ReleaseIPAllocation(logger, conf.IPAM.Type, args.StdinData)
		return nil, err
	}

	if len(result.IPs) == 0 {
		utils.ReleaseIPAllocation(logger, conf.IPAM.Type, args.StdinData)
		return nil, errors.New("IPAM plugin returned missing IP config")
	}

	// maybeReleaseIPAM cleans up any IPAM allocations if we were creating a new endpoint;
	// it is a no-op if this was a re-network of an existing endpoint.
	maybeReleaseIPAM := func() {
		logger.Debug("Checking if we need to clean up IPAM.")
//		logger := logger.WithField("IPs", endpoint.Spec.IPNetworks)
		logger.Info("Releasing IPAM allocation after failure")
		utils.ReleaseIPAllocation(logger, conf.IPAM.Type, args.StdinData)
	}

	// Whether the endpoint existed or not, the veth needs (re)creating.
	podKey := fmt.Sprintf("%s/%s", epIDs.Namespace, epIDs.Pod)
	hostVethName := pod.IFNameForKey(podKey)
	_, contVethMac, err := utils.DoNetworking(args, conf, result, logger, hostVethName)
	if err != nil {
		logger.WithError(err).Error("Error setting up networking")
		maybeReleaseIPAM()
		return nil, err
	}

	mac, err := net.ParseMAC(contVethMac)
	if err != nil {
		logger.WithError(err).WithField("mac", mac).Error("Error parsing container MAC")
		maybeReleaseIPAM()
		return nil, err
	}

	portID := pod.IDForKey(podKey)
	err = midonet.RunMmCtlBind(portID, hostVethName)
	if err != nil {
		logger.WithError(err).Error("RunMmCtlBind failed")
		maybeReleaseIPAM()
		return nil, err
	}

	// REVISIT(yamamoto): Feed midonet mac/ip info

	return result, nil
}

// CmdDelK8s performs the "DEL" operation on a kubernetes pod.
// The following logic only applies to kubernetes since it sends multiple DELs for the same
// endpoint. See: https://github.com/kubernetes/kubernetes/issues/44100
// REVISIT(yamamoto): we don't unbind port (mm-ctl --unbind-port) to
// avoid issues with multiple DELs.
func CmdDelK8s(epIDs utils.WEPIdentifiers, args *skel.CmdArgs, conf types.NetConf, logger *logrus.Entry) error {

	// Release the IP address by calling the configured IPAM plugin.
	ipamErr := utils.CleanUpIPAM(conf, args, logger)

	// Clean up namespace by removing the interfaces.
	err := utils.CleanUpNamespace(args, logger)
	if err != nil {
		return err
	}

	// Return the IPAM error if there was one. The IPAM error will be lost if there was also an error in cleaning up
	// the device or endpoint, but crucially, the user will know the overall operation failed.
	if ipamErr != nil {
		return ipamErr
	}

	logger.Info("Teardown processing complete.")

	return nil
}

// ipAddrsResult parses the ipAddrs annotation and calls the configured IPAM plugin for
// each IP passed to it by setting the IP field in CNI_ARGS, and returns the result of calling the IPAM plugin.
// Example annotation value string: "[\"10.0.0.1\", \"2001:db8::1\"]"
func ipAddrsResult(ipAddrs string, conf types.NetConf, args *skel.CmdArgs, logger *logrus.Entry) (*current.Result, error) {
	logger.Infof("Parsing annotation \"cni.projectcalico.org/ipAddrs\":%s", ipAddrs)

	// We need to make sure there is only one IPv4 and/or one IPv6
	// passed in, since CNI spec only supports one of each right now.
	ipList, err := validateAndExtractIPs(ipAddrs, "cni.projectcalico.org/ipAddrs", logger)
	if err != nil {
		return nil, err
	}

	result := current.Result{}

	// Go through all the IPs passed in as annotation value and call IPAM plugin
	// for each, and populate the result variable with IP4 and/or IP6 IPs returned
	// from the IPAM plugin.
	for _, ip := range ipList {
		// Call callIPAMWithIP with the ip address.
		r, err := callIPAMWithIP(ip, conf, args, logger)
		if err != nil {
			return nil, fmt.Errorf("error getting IP from IPAM: %s", err)
		}

		result.IPs = append(result.IPs, r.IPs[0])
		logger.Debugf("Adding IPv%s: %s to result", r.IPs[0].Version, ip.String())
	}

	return &result, nil
}

// callIPAMWithIP sets CNI_ARGS with the IP and calls the IPAM plugin with it
// to get current.Result and then it unsets the IP field from CNI_ARGS ENV var,
// so it doesn't pollute the subsequent requests.
func callIPAMWithIP(ip net.IP, conf types.NetConf, args *skel.CmdArgs, logger *logrus.Entry) (*current.Result, error) {

	// Save the original value of the CNI_ARGS ENV var for backup.
	originalArgs := os.Getenv("CNI_ARGS")
	logger.Debugf("Original CNI_ARGS=%s", originalArgs)

	ipamArgs := struct {
		cnitypes.CommonArgs
		IP net.IP `json:"ip,omitempty"`
	}{}

	if err := cnitypes.LoadArgs(args.Args, &ipamArgs); err != nil {
		return nil, err
	}

	if ipamArgs.IP != nil {
		logger.Errorf("'IP' variable already set in CNI_ARGS environment variable.")
	}

	// Request the provided IP address using the IP CNI_ARG.
	// See: https://github.com/containernetworking/cni/blob/master/CONVENTIONS.md#cni_args for more info.
	newArgs := originalArgs + ";IP=" + ip.String()
	logger.Debugf("New CNI_ARGS=%s", newArgs)

	// Set CNI_ARGS to the new value.
	err := os.Setenv("CNI_ARGS", newArgs)
	if err != nil {
		return nil, fmt.Errorf("error setting CNI_ARGS environment variable: %v", err)
	}

	// Run the IPAM plugin.
	logger.Debugf("Calling IPAM plugin %s", conf.IPAM.Type)
	r, err := ipam.ExecAdd(conf.IPAM.Type, args.StdinData)
	if err != nil {
		// Restore the CNI_ARGS ENV var to it's original value,
		// so the subsequent calls don't get polluted by the old IP value.
		if err := os.Setenv("CNI_ARGS", originalArgs); err != nil {
			logger.Errorf("Error setting CNI_ARGS environment variable: %v", err)
		}
		return nil, err
	}
	logger.Debugf("IPAM plugin returned: %+v", r)

	// Restore the CNI_ARGS ENV var to it's original value,
	// so the subsequent calls don't get polluted by the old IP value.
	if err := os.Setenv("CNI_ARGS", originalArgs); err != nil {
		// Need to clean up IP allocation if this step doesn't succeed.
		utils.ReleaseIPAllocation(logger, conf.IPAM.Type, args.StdinData)
		logger.Errorf("Error setting CNI_ARGS environment variable: %v", err)
		return nil, err
	}

	// Convert IPAM result into current Result.
	// IPAM result has a bunch of fields that are optional for an IPAM plugin
	// but required for a CNI plugin, so this is to populate those fields.
	// See CNI Spec doc for more details.
	ipamResult, err := current.NewResultFromResult(r)
	if err != nil {
		return nil, err
	}

	if len(ipamResult.IPs) == 0 {
		return nil, errors.New("IPAM plugin returned missing IP config")
	}

	return ipamResult, nil
}

// overrideIPAMResult generates current.Result like the one produced by IPAM plugin,
// but sets IP field manually since IPAM is bypassed with this annotation.
// Example annotation value string: "[\"10.0.0.1\", \"2001:db8::1\"]"
func overrideIPAMResult(ipAddrsNoIpam string, logger *logrus.Entry) (*current.Result, error) {
	logger.Infof("Parsing annotation \"cni.projectcalico.org/ipAddrsNoIpam\":%s", ipAddrsNoIpam)

	// We need to make sure there is only one IPv4 and/or one IPv6
	// passed in, since CNI spec only supports one of each right now.
	ipList, err := validateAndExtractIPs(ipAddrsNoIpam, "cni.projectcalico.org/ipAddrsNoIpam", logger)
	if err != nil {
		return nil, err
	}

	result := current.Result{}

	// Go through all the IPs passed in as annotation value and populate
	// the result variable with IP4 and/or IP6 IPs.
	for _, ip := range ipList {
		var version string
		var mask net.IPMask

		if ip.To4() != nil {
			version = "4"
			mask = net.CIDRMask(32, 32)

		} else {
			version = "6"
			mask = net.CIDRMask(128, 128)
		}

		ipConf := &current.IPConfig{
			Version: version,
			Address: net.IPNet{
				IP:   ip,
				Mask: mask,
			},
		}
		result.IPs = append(result.IPs, ipConf)
		logger.Debugf("Adding IPv%s: %s to result", ipConf.Version, ip.String())
	}

	return &result, nil
}

// validateAndExtractIPs is a utility function that validates the passed IP list to make sure
// there is one IPv4 and/or one IPv6 and then returns the slice of IPs.
func validateAndExtractIPs(ipAddrs string, annotation string, logger *logrus.Entry) ([]net.IP, error) {
	// Parse IPs from JSON.
	ips, err := parseIPAddrs(ipAddrs, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPs %s for annotation \"%s\": %s", ipAddrs, annotation, err)
	}

	// annotation value can't be empty.
	if len(ips) == 0 {
		return nil, fmt.Errorf("annotation \"%s\" specified but empty", annotation)
	}

	var hasIPv4, hasIPv6 bool
	var ipList []net.IP

	// We need to make sure there is only one IPv4 and/or one IPv6
	// passed in, since CNI spec only supports one of each right now.
	for _, ip := range ips {
		ipAddr := net.ParseIP(ip)
		if ipAddr == nil {
			logger.WithField("IP", ip).Error("Invalid IP format")
			return nil, fmt.Errorf("invalid IP format: %s", ip)
		}

		if ipAddr.To4() != nil {
			if hasIPv4 {
				// Check if there is already has been an IPv4 in the list, as we only support one IPv4 and/or one IPv6 per interface for now.
				return nil, fmt.Errorf("cannot have more than one IPv4 address for \"%s\" annotation", annotation)
			}
			hasIPv4 = true
		} else {
			if hasIPv6 {
				// Check if there is already has been an IPv6 in the list, as we only support one IPv4 and/or one IPv6 per interface for now.
				return nil, fmt.Errorf("cannot have more than one IPv6 address for \"%s\" annotation", annotation)
			}
			hasIPv6 = true
		}

		// Append the IP to ipList slice.
		ipList = append(ipList, ipAddr)
	}

	return ipList, nil
}

// parseIPAddrs is a utility function that parses string of IPs in json format that are
// passed in as a string and returns a slice of string with IPs.
// It also makes sure the slice isn't empty.
func parseIPAddrs(ipAddrsStr string, logger *logrus.Entry) ([]string, error) {
	var ips []string

	err := json.Unmarshal([]byte(ipAddrsStr), &ips)
	if err != nil {
		return nil, fmt.Errorf("failed to parse '%s' as JSON: %s", ipAddrsStr, err)
	}

	logger.Debugf("IPs parsed: %v", ips)

	return ips, nil
}

func newK8sClient(conf types.NetConf, logger *logrus.Entry) (*kubernetes.Clientset, error) {
	// Some config can be passed in a kubeconfig file
	kubeconfig := conf.Kubernetes.Kubeconfig

	// Config can be overridden by config passed in explicitly in the network config.
	configOverrides := &clientcmd.ConfigOverrides{}

	// Also allow the K8sAPIRoot to appear under the "kubernetes" block in the network config.
	if conf.Kubernetes.K8sAPIRoot != "" {
		configOverrides.ClusterInfo.Server = conf.Kubernetes.K8sAPIRoot
	}

	// Use the kubernetes client code to load the kubeconfig file and combine it with the overrides.
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	logger.Debugf("Kubernetes config %v", config)

	// Create the clientset
	return kubernetes.NewForConfig(config)
}

func getPodCidr(client *kubernetes.Clientset, conf types.NetConf, podNamespace string, podName string) (string, error) {
	// Pull the node name out of the config if it's set. Defaults to nodename
	var nodename string
	if conf.Kubernetes.NodeName != "" {
		nodename = conf.Kubernetes.NodeName
	} else {
		pod, err := client.CoreV1().Pods(podNamespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		nodename = pod.Spec.NodeName
	}

	node, err := client.CoreV1().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if node.Spec.PodCIDR == "" {
		return "", fmt.Errorf("no podCidr for node %s", nodename)
	}
	return node.Spec.PodCIDR, nil
}
