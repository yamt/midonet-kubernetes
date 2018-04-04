package midonet

import (
	"fmt"
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

type Parent struct {
	ID *uuid.UUID `json:"-"`
}

type PortRange struct {
	// Can't specify 0 explicitly but it should be ok for our usage
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/bridge.html
type Bridge struct {
	ID       *uuid.UUID `json:"id,omitempty"`
	TenantID string     `json:"tenantId,omitempty"`
	Name     string     `json:"name,omitempty"`
}

func (_ *Bridge) MediaType() string {
	return "application/vnd.org.midonet.Bridge-v4+json"
}

func (res *Bridge) Path(op string) string {
	switch op {
	case "POST":
		return "/bridges"
	case "PUT", "DELETE":
		return fmt.Sprintf("/bridges/%s", res.ID)
	default:
		return ""
	}
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/port.html
type Port struct {
	Parent
	ID         *uuid.UUID     `json:"id,omitempty"`
	Type       string         `json:"type"`
	PortSubnet []*types.IPNet `json:"portSubnet,omitempty"`
	PortMAC    HardwareAddr   `json:"portMac,omitempty"`
}

func (_ *Port) MediaType() string {
	return "application/vnd.org.midonet.Port-v3+json"
}

func (res *Port) Path(op string) string {
	switch op {
	case "POST":
		var parentType string
		switch res.Type {
		case "Bridge":
			parentType = "bridges"
		case "Router":
			parentType = "routers"
		}
		return fmt.Sprintf("/%s/%s/ports", parentType, res.Parent.ID)
	case "PUT", "DELETE":
		return fmt.Sprintf("/ports/%s", res.ID)
	default:
		return ""
	}
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/port-link.html
type PortLink struct {
	Parent
	PortID *uuid.UUID `json:"portId"`
	PeerID *uuid.UUID `json:"peerId"`
}

func (_ *PortLink) MediaType() string {
	return "application/vnd.org.midonet.PortLink-v1+json"
}

func (res *PortLink) Path(op string) string {
	switch op {
	case "POST", "DELETE":
		return fmt.Sprintf("/ports/%s/link", res.Parent.ID)
	default:
		return ""
	}
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/route.html
type Route struct {
	Parent
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

func (res *Route) Path(op string) string {
	switch op {
	case "POST":
		return fmt.Sprintf("/routers/%s/routes/%s", res.Parent.ID, res.ID)
	case "DELETE":
		return fmt.Sprintf("/routes/%s", res.ID)
	default:
		return ""
	}
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/chain.html
type Chain struct {
	ID       *uuid.UUID `json:"id,omitempty"`
	TenantID string     `json:"tenantId,omitempty"`
	Name     string     `json:"name,omitempty"`
}

func (_ *Chain) MediaType() string {
	return "application/vnd.org.midonet.Chain-v1+json"
}

func (res *Chain) Path(op string) string {
	switch op {
	case "POST":
		return "/chains"
	case "DELETE":
		return fmt.Sprintf("/chains/%s", res.ID)
	default:
		return ""
	}
}

// https://docs.midonet.org/docs/v5.4/en/rest-api/content/rule.html
type Rule struct {
	Parent
	ID           *uuid.UUID `json:"id,omitempty"`
	Type         string     `json:"type"`
	DLType       int        `json:"dlType,omitempty"`
	NwDstAddress string     `json:"nwDstAddress,omitempty"`
	NwDstLength  int        `json:"nwDstLength,omitempty"`
	NwProto      int        `json:"nwProto,omitempty"`
	NwSrcAddress string     `json:"nwSrcAddress,omitempty"`
	NwSrcLength  int        `json:"nwSrcLength,omitempty"`
	TPDST        *PortRange `json:"tpDst,omitempty"`
	TPSRC        *PortRange `json:"tpSrc,omitempty"`

	// JUMP
	JumpChainID *uuid.UUID `json:"jumpChainId,omitempty"`
}

func (_ *Rule) MediaType() string {
	return "application/vnd.org.midonet.Rule-v2+json"
}

func (res *Rule) Path(op string) string {
	switch op {
	case "POST":
		return fmt.Sprintf("/chains/%s/rules", res.Parent.ID)
	case "DELETE":
		return fmt.Sprintf("/rules/%s", res.ID)
	default:
		return ""
	}
}
