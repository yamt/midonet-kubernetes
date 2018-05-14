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
	Push(key string, parentObj interface{}, resources []APIResource) error
	Delete(key string) error
}

type Handler struct {
	converter Converter
	updater   Updater
	config    *Config

	resolver *HostResolver

	// in-core cache of sub resources.
	// REVISIT: consider making this Kubernetes Custom Resources.
	knownSubResources map[string]SubResourceMap
}

func NewHandler(converter Converter, updater Updater, config *Config) *Handler {
	client := NewClient(config)
	return &Handler{
		converter:         converter,
		updater:           updater,
		config:            config,
		resolver:          NewHostResolver(client),
		knownSubResources: make(map[string]SubResourceMap),
	}
}

func (h *Handler) deletedSubResources(key string, rs SubResourceMap) SubResourceMap {
	known := h.knownSubResources[key]
	deleted := make(SubResourceMap)
	for k, r := range known {
		if _, ok := rs[k]; !ok {
			deleted[k] = r
		}
	}
	return deleted
}

func merge(a SubResourceMap, b SubResourceMap) SubResourceMap {
	c := make(SubResourceMap)
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

func (h *Handler) handleSubResources(key string, parentObj interface{}, added SubResourceMap, deleted SubResourceMap, updater Updater, clog *log.Entry) error {
	if _, ok := h.knownSubResources[key]; !ok {
		h.knownSubResources[key] = make(SubResourceMap)
	}
	for k, _ := range merge(h.deletedSubResources(key, added), deleted) {
		err := updater.Delete(k)
		if err != nil {
			clog.WithError(err).WithFields(log.Fields{
				"key":     key,
				"sub-key": k,
			}).Error("failed to delete a sub resource")
			return err
		}
		delete(h.knownSubResources[key], k)
	}
	for k, r := range added {
		convertedSub, err := r.Convert(k, h.config)
		err = updater.Push(k, parentObj, convertedSub)
		// Remember the resource regardless of err as we might have
		// partially pushed.
		h.knownSubResources[key][k] = r
		if err != nil {
			clog.WithError(err).WithFields(log.Fields{
				"key":     key,
				"sub-key": k,
			}).Error("failed to push a sub resource")
			return err
		}
	}
	if len(h.knownSubResources[key]) == 0 {
		delete(h.knownSubResources, key)
	}
	return nil
}

func (h *Handler) Update(key string, obj interface{}) error {
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	converted, subResources, err := h.converter.Convert(key, obj, h.config, h.resolver)
	if err != nil {
		// REVISIT: this should not be fatal
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	err = h.updater.Push(key, obj, converted)
	if err != nil {
		clog.WithError(err).Error("Failed to push")
		return err
	}
	err = h.handleSubResources(key, obj, subResources, nil, h.updater, clog)
	if err != nil {
		clog.WithError(err).Error("handleSubResources")
		return err
	}
	// TODO: annotate kubernetes obj
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	converted, subResources, err := h.converter.Convert(key, nil, h.config, h.resolver)
	if err != nil {
		// REVISIT: this should not be fatal
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	err = h.handleSubResources(key, nil, nil, subResources, h.updater, clog)
	if err != nil {
		clog.WithError(err).Error("handleSubResources")
		return err
	}
	err = h.updater.Delete(key)
	if err != nil {
		clog.WithError(err).Error("Failed to delete")
		return err
	}
	return nil
}
