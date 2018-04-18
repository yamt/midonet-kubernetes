// Copyright (c) 2016 Tigera, Inc. All rights reserved.

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

package types

import (
	"net"

	"github.com/containernetworking/cni/pkg/types"
)

// Kubernetes a K8s specific struct to hold config
type Kubernetes struct {
	K8sAPIRoot string `json:"k8s_api_root"`
	Kubeconfig string `json:"kubeconfig"`
	NodeName   string `json:"node_name"`
}

type NetworkInfo struct {
	Name   string `json:"name"`
	Labels struct {
		Labels []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"labels,omitempty"`
	} `json:"labels,omitempty"`
}

// NetConf stores the common network config for this CNI plugin
type NetConf struct {
	CNIVersion string `json:"cniVersion,omitempty"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	IPAM       struct {
		Name       string
		Type       string   `json:"type"`
		Subnet     string   `json:"subnet"`
		AssignIpv4 *string  `json:"assign_ipv4"`
		AssignIpv6 *string  `json:"assign_ipv6"`
		IPv4Pools  []string `json:"ipv4_pools,omitempty"`
		IPv6Pools  []string `json:"ipv6_pools,omitempty"`
	} `json:"ipam,omitempty"`
	MTU                  int        `json:"mtu"`
	LogLevel             string     `json:"log_level"`
	Kubernetes           Kubernetes `json:"kubernetes"`
}

// K8sArgs is the valid CNI_ARGS used for Kubernetes
type K8sArgs struct {
	types.CommonArgs
	IP                         net.IP
	K8S_POD_NAME               types.UnmarshallableString
	K8S_POD_NAMESPACE          types.UnmarshallableString
	K8S_POD_INFRA_CONTAINER_ID types.UnmarshallableString
}
