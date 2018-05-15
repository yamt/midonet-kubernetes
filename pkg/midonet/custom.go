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
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/yamt/midonet-kubernetes/pkg/apis/midonet/v1"
	mncli "github.com/yamt/midonet-kubernetes/pkg/client/clientset/versioned"
)

type TranslationUpdater struct {
	client mncli.Interface
}

func NewTranslationUpdater(client mncli.Interface) *TranslationUpdater {
	return &TranslationUpdater{
		client: client,
	}
}

func (u *TranslationUpdater) Update(key string, parentKind schema.GroupVersionKind, parentObj interface{}, resources map[string][]APIResource) error {
	for k, res := range resources {
		err := u.updateOne(k, parentKind, parentObj, res)
		if err != nil {
			return err
		}
	}
	// TODO: remove stale Translations for the owner
	return nil
}

func (u *TranslationUpdater) updateOne(key string, parentKind schema.GroupVersionKind, parentObj interface{}, resources []APIResource) error {
	ns, name, err := extractNames(key)
	if err != nil {
		return err
	}
	clog := log.WithFields(log.Fields{
		"key": key,
		"namespace": ns,
		"name": name,
	})
	pmeta, err := meta.Accessor(parentObj)
	if err != nil {
		clog.WithError(err).Error("Accessor")
		return err
	}
	owners := []metav1.OwnerReference{
		{
			APIVersion: parentKind.GroupVersion().String(),
			Kind:       parentKind.Kind,
			Name:       pmeta.GetName(),
			UID:        pmeta.GetUID(),
		},
	}
	var v1rs []v1.APIResource
	for _, res := range resources {
		data, err := json.Marshal(res)
		if err != nil {
			clog.WithError(err).Error("Marshal")
			return err
		}
		r := v1.APIResource{
			Kind: TypeNameForObject(res),
			Body: string(data),
		}
		// Parent
		v1rs = append(v1rs, r)
	}
	obj := &v1.Translation{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Resources: v1rs,
	}
	meta, err := meta.Accessor(obj)
	if err != nil {
		clog.WithError(err).Error("Accessor(obj)")
		return err
	}
	meta.SetOwnerReferences(owners)
	clog = clog.WithField("obj", obj)
	clog = clog.WithField("obj-name", meta.GetName())
	newObj, err := u.client.MidonetV1().Translations(ns).Create(obj)
	if err != nil {
		clog.WithError(err).Error("Create")
		return err
	}
	clog.WithField("newObj", newObj).Info("Created CR")
	return nil
}

func (u *TranslationUpdater) Delete(key string) error {
	ns, name, err := extractNames(key)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	return u.client.MidonetV1().Translations(ns).Delete(name, &opts)
}

func extractNames(key string) (string, string, error) {
	sep := strings.SplitN(key, "/", 2)
	if len(sep) != 2 {
		return "", "", fmt.Errorf("Unrecognized key %s", key)
	}
	return sep[0], sep[1], nil
}
