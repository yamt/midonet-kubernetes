package pod

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/informers"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
)

type Handler struct {
}

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset) *controller.Controller {
	queue := controller.AddHandler(si.Core().V1().Pods().Informer(), "Pod")
	return controller.NewController(si, queue, &Handler{})
}

func (h *Handler) Handle(key interface{}) error {
	log.WithField("key", key).Info("Processing")
	return nil
}
