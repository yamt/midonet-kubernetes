package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
		svcsChainID := converter.ServicesChainID(config)
		svc := obj.(*v1.Service)
		svcIP := svc.Spec.ClusterIP
		if svc.Spec.Type != v1.ServiceTypeClusterIP || svcIP == "" || svcIP == v1.ClusterIPNone {
			return resources, nil, nil
		}
		for _, p := range svc.Spec.Ports {
			var proto int
			switch p.Protocol {
			case "TCP":
				proto = 6
			case "UDP":
				proto = 17
			default:
				log.WithField("protocol", p.Protocol).Fatal("Unknown protocol")
			}
			port := int(p.Port)
			portKey := fmt.Sprintf("%s/%s", key, p.Name)
			baseID := converter.IDForKey(portKey)
			portChainID := baseID
			jumpRuleID := converter.SubID(portChainID, "Jump to ServicePort")
			resources = append(resources, &midonet.Chain{
				ID:       &portChainID,
				Name:     fmt.Sprintf("KUBE-SVC-%s", portKey),
				TenantID: config.Tenant,
			}, &midonet.Rule{
				Parent:       midonet.Parent{ID: &svcsChainID},
				ID:           &jumpRuleID,
				DLType:       0x800,
				NWDstAddress: svcIP,
				NWSrcLength:  32,
				NWProto:      proto,
				TPDst:        &midonet.PortRange{Start: port, End: port},
				Type:         "jump",
				JumpChainID:  &portChainID,
			})
		}
	}
	return resources, nil, nil
}
