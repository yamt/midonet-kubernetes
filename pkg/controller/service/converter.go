package service

import (
	"fmt"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type serviceConverter struct{}

func newServiceConverter() midonet.Converter {
	return &serviceConverter{}
}

func (_ *serviceConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]*midonet.APIResource, error) {
	serviceChainID := midonet.IDForKey(key)
	return []*midonet.APIResource{
		{
			fmt.Sprintf("/chains"),
			"",
			fmt.Sprintf("/chains/%v", serviceChainID),
			&midonet.Chain{
				ID:   &serviceChainID,
				Name: fmt.Sprintf("KUBE-SVC-%s", key),
			},
		},
	}, nil
}
