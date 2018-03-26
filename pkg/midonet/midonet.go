package midonet

import (
	log "github.com/sirupsen/logrus"
	kapi "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/yamt/midonet-kubernetes/pkg/factory"
	"github.com/yamt/midonet-kubernetes/pkg/kube"
)

type Controller struct {
	Kube           kube.Interface
	watchFactory   *factory.WatchFactory
}

func NewMidoNetController(kubeClient kubernetes.Interface, wf *factory.WatchFactory) *Controller {
	return &Controller{
		Kube:                     &kube.Kube{KClient: kubeClient},
		watchFactory:             wf,
	}
}

func (c *Controller) Run() {
	c.WatchNodes()
	c.WatchPods()
}

func (c *Controller) WatchPods() {
	c.watchFactory.AddPodHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*kapi.Pod)
			log.WithFields(log.Fields{
				"pod": pod,
			}).Info("Add")
		},
		UpdateFunc: func(old, new interface{}) {
			oldPod := old.(*kapi.Pod)
			newPod := new.(*kapi.Pod)
			log.WithFields(log.Fields{
				"oldPod": oldPod,
				"newPod": newPod,
			}).Info("Update")
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*kapi.Pod)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					log.Errorf("couldn't get object from tombstone %+v", obj)
					return
				}
				pod, ok = tombstone.Obj.(*kapi.Pod)
				if !ok {
					log.Errorf("tombstone contained object that is not a pod %#v", obj)
					return
				}
			}
			log.WithFields(log.Fields{
				"pod": pod,
			}).Info("Delete")
		},
	}, func(pods []interface{}) {
		for i, v := range pods {
			pod, ok := v.(*kapi.Pod)
			if !ok {
				log.WithFields(log.Fields{
					"obj": v,
				}).Fatal("Spurious object in Sync")
			}
			log.WithFields(log.Fields{
				"index" : i,
				"pod": pod,
			}).Info("Sync")
		}
	})
}

func (c *Controller) WatchNodes() {
	c.watchFactory.AddNodeHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*kapi.Node)
			log.WithFields(log.Fields{
				"node": node,
			}).Info("Add")
		},
		UpdateFunc: func(old, new interface{}) {
			oldNode := old.(*kapi.Node)
			newNode := new.(*kapi.Node)
			log.WithFields(log.Fields{
				"oldNode": oldNode,
				"newNode": newNode,
			}).Info("Update")
		},
		DeleteFunc: func(obj interface{}) {
			node, ok := obj.(*kapi.Node)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					log.Errorf("couldn't get object from tombstone %+v", obj)
					return
				}
				node, ok = tombstone.Obj.(*kapi.Node)
				if !ok {
					log.Errorf("tombstone contained object that is not a node %#v", obj)
					return
				}
			}
			log.WithFields(log.Fields{
				"node": node,
			}).Info("Delete")
		},
	}, func(nodes []interface{}) {
		for i, v := range nodes {
			node, ok := v.(*kapi.Node)
			if !ok {
				log.WithFields(log.Fields{
					"obj": v,
				}).Fatal("Spurious object in Sync")
			}
			log.WithFields(log.Fields{
				"index" : i,
				"node": node,
			}).Info("Sync")
		}
	})
}
