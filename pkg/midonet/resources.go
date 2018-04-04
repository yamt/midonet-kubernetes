package midonet

import (
	"net"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/google/uuid"
)

func ParseCIDR(s string) (*types.IPNet, error) {
	tmp, err := types.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	ip := types.IPNet(*tmp)
	return &ip, nil
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/resource-models.html

type Bridge struct {
	ID       *uuid.UUID `json:"id,omitempty"`
	TenantID string     `json:"tenantId,omitempty"`
	Name     string     `json:"name,omitempty"`
}

func (_ *Bridge) MediaType() string {
	return "application/vnd.org.midonet.Bridge-v4+json"
}

type Port struct {
	ID         *uuid.UUID     `json:"id,omitempty"`
	Type       string         `json:"type"`
	PortSubnet []*types.IPNet `json:"portSubnet,omitempty"`
	PortMAC    HardwareAddr   `json:"portMac,omitempty"`
}

func (_ *Port) MediaType() string {
	return "application/vnd.org.midonet.Port-v3+json"
}

type PortLink struct {
	PortID *uuid.UUID `json:"portId"`
	PeerID *uuid.UUID `json:"peerId"`
}

func (_ *PortLink) MediaType() string {
	return "application/vnd.org.midonet.PortLink-v1+json"
}

type Route struct {
	ID               *uuid.UUID `json:"id,omitempty"`
	DstNetworkAddr   net.IP     `json:"dstNetworkAddr"`
	DstNetworkLength int        `json:"dstNetworkLength"`
	NextHopGateway   net.IP     `json:"nextHopGateway,omitempty"`
	NextHopPort      *uuid.UUID `json:"nextHopPort"`
	SrcNetworkAddr   net.IP     `json:"srcNetworkAddr"`
	SrcNetworkLength int        `json:"srcNetworkLength"`
	Type             string     `json:"type"`
}

func (_ *Route) MediaType() string {
	return "application/vnd.org.midonet.Route-v1+json"
}

type Chain struct {
	ID       *uuid.UUID `json:"id,omitempty"`
	TenantID string     `json:"tenantId,omitempty"`
	Name     string     `json:"name,omitempty"`
}

func (_ *Chain) MediaType() string {
	return "application/vnd.org.midonet.Chain-v1+json"
}

type PortRange struct {
	// Can't specify 0 explicitly but it should be ok for our usage
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type Rule struct {
	ID           *uuid.UUID `json:"id,omitempty"`
	Type         string     `json:"type"`
	DLType       int        `json:"dlType,omitempty"`
	NwDstAddress string     `json:"nwDstAddress,omitempty"`
	NwDstLength  int        `json:"nwDstLength,omitempty"`
	NwProto      int        `json:"nwProto,omitempty"`
	NwSrcAddress string     `json:"nwDstAddress,omitempty"`
	NwSrcLength  int        `json:"nwDstLength,omitempty"`
	TPDST        *PortRange `json:"tpDst",omitempty`
	TPSRC        *PortRange `json:"tpSrc",omitempty`
}

func (_ *Rule) MediaType() string {
	return "application/vnd.org.midonet.Rule-v2+json"
}
