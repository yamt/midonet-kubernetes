package node

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/plugins/pkg/ip"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type nodeConverter struct{}

func newNodeConverter() midonet.Converter {
	return &nodeConverter{}
}

func (c *nodeConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]*midonet.APIResource, error) {
	baseID := midonet.IdForKey(key)
	routerPortMac := midonet.MacForKey(key)
	routerID := config.ClusterRouter
	bridgeID := baseID
	bridgePortID := midonet.SubID(baseID, "Bridge Port")
	routerPortID := midonet.SubID(baseID, "Router Port")
	subnetRouteID := midonet.SubID(baseID, "Route")
	var routerPortSubnet []*types.IPNet
	var subnetAddr net.IP
	var subnetLen int
	var bridgeName string
	if obj != nil {
		node := obj.(*v1.Node)
		meta, err := meta.Accessor(obj)
		if err == nil {
			bridgeName = fmt.Sprintf("%s/%s", meta.GetNamespace(), meta.GetName())
		}
		addr, subnet, err := net.ParseCIDR(node.Spec.PodCIDR)
		if err != nil {
			log.WithField("node", node).Fatal("Failed to parse PodCIDR")
		}
		routerPortSubnet = []*types.IPNet{
			{ip.NextIP(addr), subnet.Mask},
		}
		subnetAddr = subnet.IP
		subnetLen, _ = subnet.Mask.Size()
	}
	return []*midonet.APIResource{
		{
			"/bridges",
			fmt.Sprintf("/bridges/%v", bridgeID),
			fmt.Sprintf("/bridges/%v", bridgeID),
			&midonet.Bridge{
				ID:   &bridgeID,
				Name: bridgeName,
			},
		},
		{
			fmt.Sprintf("/bridges/%v/ports", bridgeID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			&midonet.Port{
				ID:   &bridgePortID,
				Type: "Bridge",
			},
		},
		{
			fmt.Sprintf("/routers/%v/ports", routerID),
			fmt.Sprintf("/ports/%v", routerPortID),
			fmt.Sprintf("/ports/%v", routerPortID),
			&midonet.Port{
				ID:         &routerPortID,
				Type:       "Router",
				PortSubnet: routerPortSubnet,
				// MidoNet API automatically generates random portMac for POST.
				// On the other hand, it clears the portMac field for PUT.
				// I suspect the latter is a bug.  Use a deterministically
				// generated Mac address to avoid issues.
				// See https://midonet.atlassian.net/browse/MNA-1251
				PortMac: midonet.HardwareAddr(routerPortMac),
			},
		},
		{
			fmt.Sprintf("/routers/%v/routes", routerID),
			fmt.Sprintf("/routes/%v", subnetRouteID),
			fmt.Sprintf("/routes/%v", subnetRouteID),
			&midonet.Route{
				ID:               &subnetRouteID,
				DstNetworkAddr:   subnetAddr,
				DstNetworkLength: subnetLen,
				SrcNetworkAddr:   net.ParseIP("0.0.0.0"),
				SrcNetworkLength: 0,
				NextHopPort:      &routerPortID,
				Type:             "Normal",
			},
		},
		{
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			"",
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			&midonet.PortLink{
				// Do not specify portId to avoid a MidoNet bug.
				// See https://midonet.atlassian.net/browse/MNA-1249
				// PortID: &bridgePortID,
				PeerID: &routerPortID,
			},
		},
	}, nil
}
