package pod

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/api/core/v1"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type podConverter struct{}

func newPodConverter() midonet.Converter {
	return &podConverter{}
}

func (c *podConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, error) {
	baseID := midonet.IDForKey(key)
	bridgePortID := baseID
	var bridgeID uuid.UUID
	if obj != nil {
		pod := obj.(*v1.Pod)
		nodeName := pod.Spec.NodeName
		if nodeName == "" {
			return nil, fmt.Errorf("NodeName is not set")
		}
		bridgeID = midonet.IDForKey(nodeName)
	}
	return []midonet.APIResource{
		&midonet.Port{
			Parent: midonet.Parent{ID: &bridgeID},
			ID:     &bridgePortID,
			Type:   "Bridge",
		},
	}, nil
}
