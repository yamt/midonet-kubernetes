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
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"github.com/yamt/midonet-kubernetes/pkg/apis/midonet/v1"
	mncli "github.com/yamt/midonet-kubernetes/pkg/client/clientset/versioned"
	mninformers "github.com/yamt/midonet-kubernetes/pkg/client/informers/externalversions"
	"github.com/yamt/midonet-kubernetes/pkg/controller"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

func NewController(si informers.SharedInformerFactory, msi mninformers.SharedInformerFactory, kc *kubernetes.Clientset, mc *mncli.Clientset, config *midonet.Config) *controller.Controller {
	informer := msi.Midonet().V1().Translations().Informer()
	handler := newHandler(config)
	gvk := v1.SchemeGroupVersion.WithKind("Translation")
	return controller.NewController(gvk, informer, handler)
}