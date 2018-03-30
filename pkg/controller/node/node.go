package node

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/informers"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
)

type PodHandler struct {
}

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset) *controller.Controller {
	queue := controller.AddHandler(si.Core().V1().Nodes().Informer(), "Node")
	return controller.NewController(si, queue, &PodHandler{})
}

func (p *PodHandler) Handle(key interface{}) error {
	log.WithField("key", key).Info("Processing")
	return nil
}
