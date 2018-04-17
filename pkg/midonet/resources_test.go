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

package midonet

import (
	"encoding/json"
	"net"
	"reflect"
	"testing"

	"github.com/containernetworking/cni/pkg/types"
)

func TestBridgeEmpty(t *testing.T) {
	obj := Bridge{}
	actual, err := json.Marshal(obj)
	if err != nil {
		t.Errorf("got error %v", err)
	}
	expected := "{}"
	if string(actual) != expected {
		t.Errorf("got %v\nwant %v", string(actual), expected)
	}
}

func TestBridgeUnknownField(t *testing.T) {
	// ignore unknown fields
	blob := []byte(`{ "foo": 1 }`)
	actual := Bridge{}
	err := json.Unmarshal(blob, &actual)
	if err != nil {
		t.Errorf("got error %v", err)
	}
	expected := Bridge{}
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestBridgeArray(t *testing.T) {
	blob := []byte(`[{ "name": "foo" }, { "name": "bar" }]`)
	var actual []Bridge
	err := json.Unmarshal(blob, &actual)
	if err != nil {
		t.Errorf("got error %v", err)
	}
	expected := []Bridge{
		{Name: "foo"},
		{Name: "bar"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestRouterPort(t *testing.T) {
	blob := []byte(`{ "portSubnet": ["192.168.1.1/24", "10.0.1.9/8"] }`)
	var actual Port
	err := json.Unmarshal(blob, &actual)
	if err != nil {
		t.Errorf("got error %v", err)
	}
	ip, err := ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Errorf("got error %v", err)
	}
	ip2, err := ParseCIDR("10.0.1.9/8")
	if err != nil {
		t.Errorf("got error %v", err)
	}
	expected := Port{
		PortSubnet: []*types.IPNet{
			ip,
			ip2,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestRoute(t *testing.T) {
	blob := []byte(`{ "dstNetworkAddr": "192.168.1.0", "dstNetworkLength": 24, "nextHopGateway": "10.1.1.2" }`)
	var actual Route
	err := json.Unmarshal(blob, &actual)
	if err != nil {
		t.Errorf("got error %v", err)
	}
	expected := Route{
		DstNetworkAddr:   net.ParseIP("192.168.1.0"),
		DstNetworkLength: 24,
		NextHopGateway:   net.ParseIP("10.1.1.2"),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
