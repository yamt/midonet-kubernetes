package controller

import (
	"time"

	log "github.com/sirupsen/logrus"

	kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	informerFactory informers.SharedInformerFactory
}

func (c *Controller) Start(stopCh <-chan struct{}) {
	c.informerFactory.Start(stopCh)
}

func NewController(kc *kubernetes.Clientset, resyncPeriod time.Duration) (*Controller, error) {
	si := informers.NewSharedInformerFactory(kc, resyncPeriod)
	si.Core().V1().Nodes().Informer().AddEventHandler(
			cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*kapi.Node)
			log.WithFields(log.Fields{
				"node": node,
			}).Debug("Add")
		},
		UpdateFunc: func(old, new interface{}) {
			oldNode := old.(*kapi.Node)
			newNode := new.(*kapi.Node)
			key, err := cache.MetaNamespaceKeyFunc(old)
			if err != nil {
				log.WithError(err).Fatal("MetaNamespaceKeyFunc")
			}
			metaOld, err := meta.Accessor(old)
			if err != nil {
				log.WithError(err).Fatal("meta.Accessor for old")
			}
			metaNew, err := meta.Accessor(new)
			if err != nil {
				log.WithError(err).Fatal("meta.Accessor for new")
			}
			log.WithFields(log.Fields{
				"key": key,
				"uid": metaOld.GetUID(),
				"oldVersion": metaOld.GetResourceVersion(),
				"newVersion": metaNew.GetResourceVersion(),
			}).Info("Add")
			log.WithFields(log.Fields{
				"oldNode": oldNode,
				"newNode": newNode,
			}).Debug("Update")
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
			}).Debug("Delete")
		},
	})
	return &Controller{informerFactory: si}, nil
}
