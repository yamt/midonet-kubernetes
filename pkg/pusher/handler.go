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

package pusher

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type Handler struct {
	client *midonet.Client
	config *midonet.Config
}

func newHandler(config *midonet.Config) *Handler {
	client := midonet.NewClient(config)
	return &Handler{
		client: client,
		config: config,
	}
}

func (h *Handler) Update(key string, gvk schema.GroupVersionKind, obj interface{}) error {
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	clog.Info("pusher handler update")
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	clog.Info("pusher handler delete")
	return nil
}
