package service

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type serviceConverter struct{}

func newServiceConverter() midonet.Converter {
	return &serviceConverter{}
}

func (_ *serviceConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, error) {
	resources := make([]midonet.APIResource, 0)
	if obj != nil {
		service := obj.(*v1.Service)
		spec := service.Spec
		for _, p := range spec.Ports {
			portKey := fmt.Sprintf("%s/%s", key, p.Name)
			portChainID := midonet.IDForKey(portKey)
			resources = append(resources, &midonet.Chain{
				ID:   &portChainID,
				Name: fmt.Sprintf("KUBE-SVC-%s", portKey),
			})
		}
	}
	return resources, nil
}
