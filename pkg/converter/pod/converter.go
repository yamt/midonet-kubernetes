// Copyright (C) 2018 Midokura SARL.
// All rights reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package pod

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"

	"github.com/yamt/midonet-kubernetes/pkg/converter"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type podConverter struct{}

func newPodConverter() midonet.Converter {
	return &podConverter{}
}

func (c *podConverter) Convert(key string, obj interface{}, config *midonet.Config) ([]midonet.APIResource, midonet.SubResourceMap, error) {
	clog := log.WithField("key", key)
	baseID := converter.IDForKey("Pod", key)
	bridgePortID := baseID
	var bridgeID uuid.UUID
	if obj != nil {
		spec := obj.(*v1.Pod).Spec
		nodeName := spec.NodeName
		if nodeName == "" {
			clog.Info("NodeName is not set")
			return nil, nil, nil
		}
		bridgeID = converter.IDForKey("Node", nodeName)
	}
	return []midonet.APIResource{
		&midonet.Port{
			Parent: midonet.Parent{ID: &bridgeID},
			ID:     &bridgePortID,
			Type:   "Bridge",
		},
	}, nil, nil
}
