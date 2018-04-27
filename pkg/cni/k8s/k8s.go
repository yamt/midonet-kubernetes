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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/yamt/midonet-kubernetes/pkg/cni/midonet"
	"github.com/yamt/midonet-kubernetes/pkg/cni/types"
	"github.com/yamt/midonet-kubernetes/pkg/cni/utils"
	"github.com/yamt/midonet-kubernetes/pkg/converter/node"
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

	subnetInfo, err := node.GetSubnetInfo(podCidr)
	if err != nil {
		return nil, err
	}
	gatewayIP := subnetInfo.GatewayIP
	nodeIP := subnetInfo.NodeIP

	stdinData["ipam"].(map[string]interface{})["subnet"] = podCidr
	stdinData["ipam"].(map[string]interface{})["gateway"] = gatewayIP.String()
	fmt.Fprintf(os.Stderr, "Calico CNI passing podCidr to host-local IPAM: %s\n", podCidr)
	args.StdinData, err = json.Marshal(stdinData)
	if err != nil {
		return nil, err
	}
	logger.WithField("stdin", string(args.StdinData)).Debug("Updated stdin data")

retry_ipam:
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

	for _, ip := range result.IPs {
		if bytes.Equal(ip.Address.IP, nodeIP.IP) {
			// Just leak it and retry
			goto retry_ipam
		}
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
	contVethMac, err := utils.DoNetworking(result.IPs, args.Netns, args.IfName, hostVethName, logger)
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
	return GetNodePodCIDR(client, nodename)
}

func GetNodePodCIDR(client *kubernetes.Clientset, nodename string) (string, error) {
	node, err := client.CoreV1().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if node.Spec.PodCIDR == "" {
		return "", fmt.Errorf("no podCidr for node %s", nodename)
	}
	return node.Spec.PodCIDR, nil
}
