package midonet

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/types"
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
	var routerPortSubnet []*types.IPNet
	if obj != nil {
		node := obj.(*v1.Node)
		subnet, err := ParseCIDR(node.Spec.PodCIDR)
		if err != nil {
			log.WithField("node", node).Fatal("Failed to parse PodCIDR")
		}
		routerPortSubnet = []*types.IPNet{subnet}
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
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			"",
			fmt.Sprintf("/ports/%v/link", bridgePortID),
			"application/vnd.org.midonet.PortLink-v1+json",
			&PortLink{
				PortID: &bridgePortID,
				PeerID: &routerPortID,
				Type: "Bridge",
			},
		},
	}, nil
}
