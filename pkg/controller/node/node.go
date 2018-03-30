package node

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

type Handler struct {
}

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset, config *midonet.Config) *controller.Controller {
	informer := si.Core().V1().Nodes().Informer()
	return controller.NewController("Node", informer, &Handler{})
}

func (h *Handler) Handle(key string, informer cache.SharedIndexInformer) error {
	clog := log.WithField("key", key)
	clog.Info("Processing")
	obj, exists, err := informer.GetIndexer().GetByKey(key)
	if err != nil {
		clog.WithError(err).Fatal("GetBykey")
	}
	if !exists {
		clog.Info("Deleted")
		return nil
	}
	clog.WithField("obj", obj).Info("Updated")
	return nil
}
