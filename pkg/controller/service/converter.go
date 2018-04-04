package service

import (
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type serviceConverter struct{}

func newServiceConverter() midonet.Converter {
	return &serviceConverter{}
}

func (c *serviceConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]*midonet.APIResource, error) {
	return []*midonet.APIResource{
	}, nil
}
