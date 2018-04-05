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

func (_ *serviceConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, midonet.SubResourceMap, error) {
	resources := make([]midonet.APIResource, 0)
	if obj != nil {
		svc := obj.(*v1.Service)
		svcIP := svc.Spec.ClusterIP
		if svc.Spec.Type != v1.ServiceTypeClusterIP || svcIP == "" || svcIP == v1.ClusterIPNone {
			return resources, nil, nil
		}
		for _, p := range svc.Spec.Ports {
			portKey := fmt.Sprintf("%s/%s", key, p.Name)
			portChainID := converter.IDForKey(portKey)
			resources = append(resources, &midonet.Chain{
				ID:       &portChainID,
				Name:     fmt.Sprintf("KUBE-SVC-%s", portKey),
				TenantID: config.Tenant,
			})
		}
	}
	return resources, nil, nil
}
