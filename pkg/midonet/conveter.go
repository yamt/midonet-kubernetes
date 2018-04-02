package midonet

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/types"
	log "github.com/sirupsen/logrus"
	 "k8s.io/api/core/v1"
)

type MidoNetAPIResource struct {
	PathForPost		string
	PathForPut		string
	PathForDelete	string
	Body			interface{}
}

// Common notes about ConvertXXX functions.
// - if nil obj is given, only PathForDelete fields for the
//   MidoNetAPIResource returned are valid.

func ConvertNode(key string, obj interface{}, config *Config) ([]MidoNetAPIResource, error) {
	baseID := idForKey(key)
	routerID := config.ClusterRouter
	bridgeID := baseID
	bridgePortID := subID(baseID, "Bridge Port")
	routerPortID := subID(baseID, "Router Port")
	var routerPortSubnet []*types.IPNet
	if obj != nil {
		node := obj.(*v1.Node)
		subnet, err := ParseCIDR(node.Spec.PodCIDR)
		if err != nil {
			log.WithField("node", node).Fatal("Failed to parse PodCIDR")
		}
		routerPortSubnet = []*types.IPNet{subnet}
	}
	return []MidoNetAPIResource{
		{
			"/bridges",
			fmt.Sprintf("/bridges/%v", bridgeID),
			fmt.Sprintf("/bridges/%v", bridgeID),
			&Bridge{
				ID: &bridgeID,
			},
		},
		{
			fmt.Sprintf("/bridges/%v/ports", bridgeID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			&Port{
				ID: &bridgePortID,
				Type: "Bridge",
			},
		},
		{
			fmt.Sprintf("/routers/%v/ports", routerID),
			fmt.Sprintf("/ports/%v", routerPortID),
			fmt.Sprintf("/ports/%v", routerPortID),
			&Port{
				ID: &routerPortID,
				Type: "Router",
				PortSubnet: routerPortSubnet,
			},
		},
		{
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			"",
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			&PortLink{
				PortID: &bridgePortID,
				PeerID: &routerPortID,
			},
		},
	}, nil
}
