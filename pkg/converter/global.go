package converter

import (
	"github.com/google/uuid"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

func ServicesChainID(config *midonet.Config) uuid.UUID {
	baseID := IDForTenant(config.Tenant)
	return SubID(baseID, "Services Chain")
}

func MainChainID(config *midonet.Config) uuid.UUID {
	return IDForTenant(config.Tenant)
}

func GlobalResources(config *midonet.Config) []midonet.APIResource {
	tenant := config.Tenant
	baseID := IDForTenant(tenant)
	mainChainID := baseID
	preChainID := SubID(baseID, "Pre Chain")
	servicesChainID := ServicesChainID(config)
	jumpToPreRuleID := SubID(baseID, "Jump To Pre")
	jumpToServicesRuleID := SubID(baseID, "Jump To Services")
	revSNATRuleID := SubID(baseID, "Reverse SNAT")
	revDNATRuleID := SubID(baseID, "Reverse DNAT")
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
		// REVISIT: Ensure the order of rules
		midonet.JumpRule(&jumpToServicesRuleID, &mainChainID, &servicesChainID),
		midonet.JumpRule(&jumpToPreRuleID, &mainChainID, &preChainID),
		// Reverse NAT rules for Endpoints.
		// REVISIT: Ensure the order of rules
		&midonet.Rule{
			Parent:     midonet.Parent{ID: &preChainID},
			ID:         &revDNATRuleID,
			Type:       "rev_dnat",
			FlowAction: "accept",
		},
		&midonet.Rule{
			Parent:     midonet.Parent{ID: &preChainID},
			ID:         &revSNATRuleID,
			Type:       "rev_snat",
			FlowAction: "continue",
		},
	}
}

func EnsureGlobalResources(config *midonet.Config) error {
	resources := GlobalResources(config)
	cli := midonet.NewClient(config)
	return cli.Push(resources)
}
