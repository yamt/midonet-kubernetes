package endpoints

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type endpointsConverter struct{}

func newEndpointsConverter() midonet.Converter {
	return &endpointsConverter{}
}

type endpoint struct {
	ip       string
	port     int
	protocol v1.Protocol
}

func endpoints(subsets []v1.EndpointSubset) map[string][]endpoint {
	m := make(map[string][]endpoint, 0)
	for _, s := range subsets {
		for _, a := range s.Addresses {
			for _, p := range s.Ports {
				ep := endpoint{a.IP, int(p.Port), p.Protocol}
				l := m[p.Name]
				l = append(l, ep)
				m[p.Name] = l
			}
		}
	}
	return m
}

func (c *endpointsConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, error) {
	resources := make([]midonet.APIResource, 0)
	if obj != nil {
		endpoint := obj.(*v1.Endpoints)
		for k, eps := range endpoints(endpoint.Subsets) {
			for _, ep := range eps {
				portKey := fmt.Sprintf("%s/%s", key, k)
				portChainID := midonet.IDForKey(portKey)
				epKey := fmt.Sprintf("%s:%s:%d:%s", portKey, ep.ip, ep.port, ep.protocol)
				baseID := midonet.IDForKey(epKey)
				epChainID := baseID
				epJumpRuleID := midonet.SubID(baseID, "Jump to Endpoint")
				resources = append(resources, &midonet.Chain{
					ID:   &epChainID,
					Name: fmt.Sprintf("KUBE-SEP-%s", epKey),
				})
				resources = append(resources, &midonet.Rule{
					Parent:      midonet.Parent{ID: &portChainID},
					ID:          &epJumpRuleID,
					Type:        "Jump",
					JumpChainID: &epChainID,
				})
			}
		}
	}
	return resources, nil
}
