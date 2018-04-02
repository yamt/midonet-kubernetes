package midonet

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/plugins/pkg/ip"
	log "github.com/sirupsen/logrus"
	 "k8s.io/api/core/v1"
)

type APIResource struct {
	PathForPost		string
	PathForPut		string
	PathForDelete	string
	MediaType		string
	Body			interface{}
}

// Common notes about ConvertXXX functions.
// - if nil obj is given, only PathForDelete fields for the
//   APIResource returned are valid.

func ConvertNode(key string, obj interface{}, config *Config) ([]*APIResource, error) {
	baseID := idForKey(key)
	routerID := config.ClusterRouter
	bridgeID := baseID
	bridgePortID := subID(baseID, "Bridge Port")
	routerPortID := subID(baseID, "Router Port")
	subnetRouteID := subID(baseID, "Route")
	var routerPortSubnet []*types.IPNet
	var subnet types.IPNet
	var subnetLen int
	if obj != nil {
		node := obj.(*v1.Node)
		addr, subnet, err := net.ParseCIDR(node.Spec.PodCIDR)
		if err != nil {
			log.WithField("node", node).Fatal("Failed to parse PodCIDR")
		}
		routerPortSubnet = []*types.IPNet{
			{ip.NextIP(addr), subnet.Mask},
		}
		subnetLen, _ = subnet.Mask.Size()
	}
	return []*APIResource{
		{
			"/bridges",
			fmt.Sprintf("/bridges/%v", bridgeID),
			fmt.Sprintf("/bridges/%v", bridgeID),
			"application/vnd.org.midonet.Bridge-v4+json",
			&Bridge{
				ID: &bridgeID,
			},
		},
		{
			fmt.Sprintf("/bridges/%v/ports", bridgeID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			"application/vnd.org.midonet.Port-v3+json",
			&Port{
				ID: &bridgePortID,
				Type: "Bridge",
			},
		},
		{
			fmt.Sprintf("/routers/%v/ports", routerID),
			fmt.Sprintf("/ports/%v", routerPortID),
			fmt.Sprintf("/ports/%v", routerPortID),
			"application/vnd.org.midonet.Port-v3+json",
			&Port{
				ID: &routerPortID,
				Type: "Router",
				PortSubnet: routerPortSubnet,
			},
		},
		{
			fmt.Sprintf("/routers/%v/routes", routerID),
			fmt.Sprintf("/routes/%v", subnetRouteID),
			fmt.Sprintf("/routes/%v", subnetRouteID),
			"application/vnd.org.midonet.Route-v1+json",
			&Route{
				ID: &subnetRouteID,
				DstNetworkAddr: subnet.IP,
				DstNetworkLength: subnetLen,
				SrcNetworkAddr: net.ParseIP("0.0.0.0"),
				SrcNetworkLength: 0,
				NextHopPort: &routerPortID,
				Type: "Normal",
			},
		},
		{
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			"",
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			"application/vnd.org.midonet.PortLink-v1+json",
			&PortLink{
				// bug
				// PortID: &bridgePortID,
				PeerID: &routerPortID,
			},
		},
	}, nil
}