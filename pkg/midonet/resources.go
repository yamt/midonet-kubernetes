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

type Port struct {
	ID         *uuid.UUID       `json:"id,omitempty"`
	Type       string           `json:"type"`
	PortSubnet []*types.IPNet   `json:"portSubnet,omitempty"`
	PortMac    net.HardwareAddr `json:"portMac,omitempty"`
}

type PortLink struct {
	PortID *uuid.UUID `json:"portId"`
	PeerID *uuid.UUID `json:"peerId"`
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
