package pod

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset, config *midonet.Config) *controller.Controller {
	informer := si.Core().V1().Pods().Informer()
	handler := midonet.NewHandler(midonet.NewPodConverter(), config)
	return controller.NewController("Pod", informer, handler)
}
