// Copyright (C) 2018 Midokura SARL.
// All rights reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"io"
	"text/template"
)

const (
	cniConfigTemplate = `
{
  "name": "midonet-pod-network",
  "type": "midonet-kube-cni",
  "ipam": {
    "type": "host-local"
  },
  "kubernetes": {
    "podcidr": "{{ .PodCIDR }}"
  }
}`
)

type cniConfigData struct {
	PodCIDR string
}

func generateCNIConfig(writer io.Writer, podCIDR string) error {
	tmpl, err := template.New("cniconfig").Parse(cniConfigTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(writer, &cniConfigData{PodCIDR: podCIDR})
}
