package pod

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type Handler struct {
}

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset, config *midonet.Config) *controller.Controller {
	informer := si.Core().V1().Pods().Informer()
	return controller.NewController("Pod", informer, &Handler{})
}

func (h *Handler) Update(key string, obj interface{}) error {
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	converted, err := midonet.ConvertPod(key, obj, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	cli := midonet.NewClient(h.config)
	err = cli.Push(converted)
	if err != nil {
		clog.WithError(err).Fatal("Failed to push")
	}
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	clog.Info("On Delete")
	return nil
}
