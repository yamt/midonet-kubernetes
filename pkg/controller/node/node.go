package node

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/informers"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type Handler struct {
	config *midonet.Config
}

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset, config *midonet.Config) *controller.Controller {
	informer := si.Core().V1().Nodes().Informer()
	return controller.NewController("Node", informer, &Handler{config})
}

func (h *Handler) Update(key string, obj interface{}) error {
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	converted, err := midonet.ConvertNode(key, obj, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	err = midonet.Push(converted, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to push")
	}
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	converted, err := midonet.ConvertNode(key, nil, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	return nil
}
