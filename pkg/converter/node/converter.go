package node

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/plugins/pkg/ip"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/yamt/midonet-kubernetes/pkg/converter"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type nodeConverter struct{}

func newNodeConverter() midonet.Converter {
	return &nodeConverter{}
}

func (c *nodeConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, midonet.SubResourceMap, error) {
	baseID := converter.IDForKey(key)
	routerPortMAC := converter.MACForKey(key)
	routerID := config.ClusterRouter
	bridgeID := baseID
	bridgePortID := converter.SubID(baseID, "Bridge Port")
	routerPortID := converter.SubID(baseID, "Router Port")
	subnetRouteID := converter.SubID(baseID, "Route")
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
	svcsChainID := converter.ServicesChainID(config)
	return []midonet.APIResource{
		&midonet.Bridge{
			ID:            &bridgeID,
			Name:          bridgeName,
			TenantID:      config.Tenant,
			InboundFilter: &svcsChainID,
		},
		&midonet.Port{
			Parent: midonet.Parent{ID: &bridgeID},
			ID:     &bridgePortID,
			Type:   "Bridge",
		},
		&midonet.Port{
			Parent:     midonet.Parent{ID: &routerID},
			ID:         &routerPortID,
			Type:       "Router",
			PortSubnet: routerPortSubnet,
			// MidoNet API automatically generates random portMac for POST.
			// On the other hand, it clears the portMac field for PUT.
			// I suspect the latter is a bug.  Use a deterministically
			// generated MAC address to avoid issues.
			// See https://midonet.atlassian.net/browse/MNA-1251
			PortMAC: midonet.HardwareAddr(routerPortMAC),
		},
		&midonet.Route{
			Parent:           midonet.Parent{ID: &routerID},
			ID:               &subnetRouteID,
			DstNetworkAddr:   subnetAddr,
			DstNetworkLength: subnetLen,
			SrcNetworkAddr:   net.ParseIP("0.0.0.0"),
			SrcNetworkLength: 0,
			NextHopPort:      &routerPortID,
			Type:             "Normal",
		},
		&midonet.PortLink{
			Parent: midonet.Parent{ID: &bridgePortID},
			// Do not specify portId to avoid a MidoNet bug.
			// See https://midonet.atlassian.net/browse/MNA-1249
			// PortID: &bridgePortID,
			PeerID: &routerPortID,
		},
	}, nil, nil
}
