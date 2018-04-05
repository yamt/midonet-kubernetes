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
	subs := make(midonet.SubResourceMap)
	if obj != nil {
		svc := obj.(*v1.Service)
		svcIP := svc.Spec.ClusterIP
		if svc.Spec.Type != v1.ServiceTypeClusterIP || svcIP == "" || svcIP == v1.ClusterIPNone {
			return resources, nil, nil
		}
		for _, p := range svc.Spec.Ports {
			// Note: portKey format should be consistent with the
			// endpoints converter so that it can find the right chain
			// to add rules. (portChainID)
			// NameSpace/Name/ServicePort.Name
			portKey := fmt.Sprintf("%s/%s", key, p.Name)
			portChainID := converter.IDForKey(portKey)
			resources = append(resources, &midonet.Chain{
				ID:       &portChainID,
				Name:     fmt.Sprintf("KUBE-SVC-%s", portKey),
				TenantID: config.Tenant,
			})

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
			// Use a separate key for sub resource so that those will
			// be deleted and re-created whenever they got changed.
			// Note that MidoNet Rules are not updateable.
			// The above KUBE-SVC- chain is not a part of this sub resource
			// because we want to avoid re-creating the chain itself
			// as it would remove rules in the chain.  (Those rules are
			// managed by a separate "endpoints" controller.)
			portSubKey := fmt.Sprintf("%s/serviceport/%s/%d/%d", portKey, svcIP, proto, port)
			subs[portSubKey] = &servicePort{portKey, svcIP, proto, port}
		}
	}
	return resources, subs, nil
}

type servicePort struct {
	portKey string
	ip      string
	proto   int
	port    int
}

func (s *servicePort) Convert(key string, config *midonet.Config) ([]midonet.APIResource, error) {
	svcsChainID := converter.ServicesChainID(config)
	jumpRuleID := converter.IDForKey(key)
	portChainID := converter.IDForKey(s.portKey)
	return []midonet.APIResource{
		&midonet.Rule{
			Parent:       midonet.Parent{ID: &svcsChainID},
			ID:           &jumpRuleID,
			DLType:       0x800,
			NWDstAddress: s.ip,
			NWSrcLength:  32,
			NWProto:      s.proto,
			TPDst:        &midonet.PortRange{Start: s.port, End: s.port},
			Type:         "jump",
			JumpChainID:  &portChainID,
		},
	}, nil
}
