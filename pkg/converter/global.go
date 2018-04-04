package converter

import (
	"github.com/google/uuid"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

func ServiceChainID(tenant string) uuid.UUID {
	baseID := IDForTenant(tenant)
	return SubID(baseID, "Services Chain")
}

func GlobalResources(tenant string) []midonet.APIResource {
	baseID := IDForTenant(tenant)
	mainChainID := baseID
	preChainID := SubID(baseID, "Pre Chain")
	servicesChainID := ServiceChainID(tenant)
	jumpToPreRuleID := SubID(baseID, "Jump To Pre")
	jumpToServicesRuleID := SubID(baseID, "Jump To Services")
	return []midonet.APIResource{
		&midonet.Chain{
			ID:   &mainChainID,
			Name: "KUBE-MAIN",
		},
		&midonet.Chain{
			ID:   &preChainID,
			Name: "KUBE-PRE",
		},
		&midonet.Chain{
			ID:   &servicesChainID,
			Name: "KUBE-SERVICES",
		},
		midonet.JumpRule(&jumpToPreRuleID, &mainChainID, &preChainID),
		midonet.JumpRule(&jumpToServicesRuleID, &mainChainID, &servicesChainID),
	}
}
