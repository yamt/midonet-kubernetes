package midonet

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/api/core/v1"
)

type PodConverter struct{}

func (c *PodConverter) Convert(key string, obj interface{}, config *Config) ([]*APIResource, error) {
	baseID := idForKey(key)
	bridgePortID := baseID
	var bridgeID uuid.UUID
	if obj != nil {
		pod := obj.(*v1.Pod)
		nodeName := pod.Spec.NodeName
		if nodeName == "" {
			return nil, fmt.Errorf("NodeName is not set")
		}
		bridgeID = idForKey(nodeName)
	}
	return []*APIResource{
		{
			fmt.Sprintf("/bridges/%v/ports", bridgeID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			fmt.Sprintf("/ports/%v", bridgePortID),
			"application/vnd.org.midonet.Port-v3+json",
			&Port{
				ID:   &bridgePortID,
				Type: "Bridge",
			},
		},
	}, nil
}
