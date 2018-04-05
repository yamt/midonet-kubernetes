package endpoints

import (
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/yamt/midonet-kubernetes/pkg/converter"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type endpointsConverter struct {
	svcInformer cache.SharedIndexInformer
}

func newEndpointsConverter(svcInformer cache.SharedIndexInformer) midonet.Converter {
	return &endpointsConverter{svcInformer}
}

type endpoint struct {
	svcIP    string
	ip       string
	port     int
	protocol v1.Protocol
}

func portKeyFromEPKey(epKey string) string {
	// a epKey looks like Namespace/Name/EndpointPort.Name/...
	sep := strings.Split(epKey, "/")
	return strings.Join(sep[:3], "/")
}

func (ep *endpoint) Convert(epKey string, config *midonet.Config) ([]midonet.APIResource, error) {
	// REVISIT: An assumption here is that, if ServicePort.Name is empty,
	// the corresponding EndpointPort.Name is also empty.  It isn't clear
	// to me (yamamoto) from the documentation.
	portKey := portKeyFromEPKey(epKey)
	portChainID := converter.IDForKey(portKey)
	baseID := converter.IDForKey(epKey)
	epChainID := baseID
	epJumpRuleID := converter.SubID(baseID, "Jump to Endpoint")
	epDNATRuleID := converter.SubID(baseID, "DNAT")
	epSNATRuleID := converter.SubID(baseID, "SNAT")
	return []midonet.APIResource{
		&midonet.Chain{
			ID:       &epChainID,
			Name:     fmt.Sprintf("KUBE-SEP-%s", epKey),
			TenantID: config.Tenant,
		},
		// REVISIT: kube-proxy implements load-balancing with its
		// equivalent of this rule, using iptables probabilistic
		// match.  We can probably implement something similar
		// here if the backend has the following functionalities.
		//
		//   1. probabilistic match
		//   2. applyIfExists equivalent
		//
		// For now, we just install a normal 100% matching rule.
		// It means that the endpoint which happens to have its
		// jump rule the earliest in the chain handles 100% of
		// traffic.
		&midonet.Rule{
			Parent:      midonet.Parent{ID: &portChainID},
			ID:          &epJumpRuleID,
			Type:        "jump",
			JumpChainID: &epChainID,
		},
		&midonet.Rule{
			Parent: midonet.Parent{ID: &epChainID},
			ID:     &epDNATRuleID,
			Type:   "dnat",
			NatTargets: &[]midonet.NatTarget{
				{
					AddressFrom: ep.ip,
					AddressTo:   ep.ip,
					PortFrom:    ep.port,
					PortTo:      ep.port,
				},
			},
			FlowAction: "accept",
		},
		// SNAT traffic from the endpoint itself. Otherwise,
		// the return traffic doesn't work.
		// Note: Endpoint IP might or might not belong to the cluster ip
		// range.  It can be external.
		//
		// The source IP to use for this purpose is somewhat arbitrary
		// and doesn't seem consistent among networking implementations.
		// We use the ClusterIP of the corresponding Service.
		// For example, kube-proxy uses iptables MASQUERADE target for
		// this purpose.  It means that the source IP of the outgoing
		// interface is chosen after an L3 routing decision.  With flannel,
		// it would be the address of the cni0 interface on the node.
		&midonet.Rule{
			Parent:       midonet.Parent{ID: &epChainID},
			ID:           &epSNATRuleID,
			Type:         "snat",
			DLType:       0x800,
			NwSrcAddress: ep.ip,
			NwSrcLength:  32,
			NatTargets: &[]midonet.NatTarget{
				{
					AddressFrom: ep.svcIP,
					AddressTo:   ep.svcIP,
					// REVISIT: arbitrary port range
					PortFrom: 30000,
					PortTo:   60000,
				},
			},
			FlowAction: "continue",
		},
	}, nil
}

func endpoints(svcIP string, subsets []v1.EndpointSubset) map[string][]endpoint {
	m := make(map[string][]endpoint, 0)
	for _, s := range subsets {
		for _, a := range s.Addresses {
			for _, p := range s.Ports {
				ep := endpoint{svcIP, a.IP, int(p.Port), p.Protocol}
				l := m[p.Name]
				l = append(l, ep)
				m[p.Name] = l
			}
		}
	}
	return m
}

func (c *endpointsConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, midonet.SubResourceMap, error) {
	// Just return each endpoints as SubResource.
	resources := make([]midonet.APIResource, 0)
	subs := make(midonet.SubResourceMap)
	if obj != nil {
		svcObj, exists, err := c.svcInformer.GetIndexer().GetByKey(key)
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			// REVISIT: This might not be a transient state.  E.g. when
			// a user created Endpoints manually.
			return nil, nil, fmt.Errorf("No corresponding service")
		}
		svc := svcObj.(*v1.Service)
		svcIP := svc.Spec.ClusterIP
		if svc.Spec.Type != v1.ServiceTypeClusterIP || svcIP == "" || svcIP == v1.ClusterIPNone {
			return nil, nil, nil
		}
		endpoint := obj.(*v1.Endpoints)
		for portName, eps := range endpoints(svcIP, endpoint.Subsets) {
			portKey := fmt.Sprintf("%s/%s", key, portName)
			for _, ep := range eps {
				// We include almost everything in the key so that a modified
				// endpoint is treated as another resource for the
				// MidoNet side.  Note that MidoNet Chains and Rules are not
				// updateable.
				epKey := fmt.Sprintf("%s/%s/%s/%d/%s", portKey, svcIP, ep.ip, ep.port, ep.protocol)
				subs[epKey] = &ep
			}
		}
	}
	return resources, subs, nil
}
