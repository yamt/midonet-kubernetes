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
