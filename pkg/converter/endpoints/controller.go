package endpoints

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"github.com/yamt/midonet-kubernetes/pkg/controller"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
)

func NewController(si informers.SharedInformerFactory, kc *kubernetes.Clientset, config *midonet.Config) *controller.Controller {
	informer := si.Core().V1().Endpoints().Informer()
	svcInformer := si.Core().V1().Services().Informer()
	handler := midonet.NewHandler(newEndpointsConverter(svcInformer), config)
	c := controller.NewController("Endpoints", informer, handler)
	// Kick the Endpoints controller when the corresponding Service is updated.
	svcInformer.AddEventHandler(controller.NewEventHandler("svc-eps", c.GetQueue()))
	return c
}
