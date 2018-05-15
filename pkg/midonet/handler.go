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

package midonet

import (
	log "github.com/sirupsen/logrus"
)

// SubResource is a pseudo resource to represent a part of a k8s resource.
// For example, we represent a k8s service as a set of "ServicePort"
// sub resources.
type SubResource interface {
	Convert(key string, config *Config) ([]APIResource, error)
}

type SubResourceMap map[string]SubResource

type Updater interface {
	Update(key string, parentObj interface{}, resources map[string][]APIResource) error
	Delete(key string) error
}

type Handler struct {
	converter Converter
	updater   Updater
	config    *Config

	resolver *HostResolver
}

func NewHandler(converter Converter, updater Updater, config *Config) *Handler {
	client := NewClient(config)
	return &Handler{
		converter:         converter,
		updater:           updater,
		config:            config,
		resolver:          NewHostResolver(client),
	}
}

func (h *Handler) convertSubResources(key string, parentObj interface{}, added SubResourceMap, converted map[string][]APIResource, clog *log.Entry) error {
	for k, r := range added {
		v, err := r.Convert(k, h.config)
		if err != nil {
			clog.WithError(err).WithFields(log.Fields{
				"key":     key,
				"sub-key": k,
			}).Error("failed to convert a sub resource")
			return err
		}
		if len(v) > 0 {
			converted[k] = v
		}
	}
	return nil
}

func (h *Handler) Update(key string, obj interface{}) error {
	converted := make(map[string][]APIResource)
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	v, subResources, err := h.converter.Convert(key, obj, h.config, h.resolver)
	if err != nil {
		// REVISIT: this should not be fatal
		clog.WithError(err).Fatal("Failed to convert")
	}
	if len(v) > 0 {
		converted[key] = v
	}
	err = h.convertSubResources(key, obj, subResources, converted, clog)
	if err != nil {
		// REVISIT: this should not be fatal
		clog.WithError(err).Fatal("Failed to convert sub resources")
	}
	err = h.updater.Update(key, obj, converted)
	if err != nil {
		clog.WithError(err).Error("Failed to update")
		return err
	}
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	err := h.updater.Delete(key)
	if err != nil {
		clog.WithError(err).Error("Failed to delete")
		return err
	}
	return nil
}
