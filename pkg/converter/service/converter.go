package service

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/yamt/midonet-kubernetes/pkg/converter"
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
		if spec.Type != v1.ServiceTypeClusterIP || spec.ClusterIP == "" {
			return resources, nil
		}
		// REVISIT: what to do for ClusterIPNone?
		// Note: ClusterIP can't be changed
		// https://github.com/kubernetes/kubernetes/blob/1102fd0dcbc4a408045e8d1bc42f056909e72322/staging/src/k8s.io/api/core/v1/types.go#L3468
		for _, p := range spec.Ports {
			portKey := fmt.Sprintf("%s/%s", key, p.Name)
			portChainID := converter.IDForKey(portKey)
			resources = append(resources, &midonet.Chain{
				ID:       &portChainID,
				Name:     fmt.Sprintf("KUBE-SVC-%s", portKey),
				TenantID: config.Tenant,
			})
		}
	}
	return resources, nil
}
