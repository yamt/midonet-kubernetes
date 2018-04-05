package converter

import (
	"github.com/google/uuid"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

func ServicesChainID(config *midonet.Config) uuid.UUID {
	baseID := IDForTenant(config.Tenant)
	return SubID(baseID, "Services Chain")
}

func GlobalResources(config *midonet.Config) []midonet.APIResource {
	tenant := config.Tenant
	baseID := IDForTenant(tenant)
	mainChainID := baseID
	preChainID := SubID(baseID, "Pre Chain")
	servicesChainID := ServicesChainID(config)
	jumpToPreRuleID := SubID(baseID, "Jump To Pre")
	jumpToServicesRuleID := SubID(baseID, "Jump To Services")
	return []midonet.APIResource{
		&midonet.Chain{
			ID:       &mainChainID,
			Name:     "KUBE-MAIN",
			TenantID: tenant,
		},
		&midonet.Chain{
			ID:       &preChainID,
			Name:     "KUBE-PRE",
			TenantID: tenant,
		},
		&midonet.Chain{
			ID:       &servicesChainID,
			Name:     "KUBE-SERVICES",
			TenantID: tenant,
		},
		midonet.JumpRule(&jumpToServicesRuleID, &mainChainID, &servicesChainID),
		midonet.JumpRule(&jumpToPreRuleID, &mainChainID, &preChainID),
	}
}

func EnsureGlobalResources(config *midonet.Config) error {
	resources := GlobalResources(config)
	cli := midonet.NewClient(config)
	return cli.Push(resources)
}
