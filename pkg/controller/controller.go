package controller

import (
	"time"

	log "github.com/sirupsen/logrus"

	// kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	informerFactory informers.SharedInformerFactory
}

func (c *Controller) Start(stopCh <-chan struct{}) {
	c.informerFactory.Start(stopCh)
}

func NewController(kc *kubernetes.Clientset, resyncPeriod time.Duration) (*Controller, error) {
	si := informers.NewSharedInformerFactory(kc, resyncPeriod)
	addHandler(si.Core().V1().Nodes().Informer(), "Node")
	addHandler(si.Core().V1().Pods().Informer(), "Pod")
	return &Controller{informerFactory: si}, nil
}

func addHandler(informer cache.SharedIndexInformer, kind string) {
	rateLimiter := workqueue.DefaultControllerRateLimiter()
	queue := workqueue.NewNamedRateLimitingQueue(rateLimiter, kind)
	handler := newHandler(kind, queue)
        informer.AddEventHandler(handler)
}

func newHandler(kind string, queue workqueue.Interface) cache.ResourceEventHandler {
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				log.WithError(err).Fatal("MetaNamespaceKeyFunc")
			}
			meta, err := meta.Accessor(obj)
			if err != nil {
				log.WithError(err).Fatal("meta.Accessor")
			}
			log.WithFields(log.Fields{
				"kind": kind,
				"key": key,
				"uid": meta.GetUID(),
				"version": meta.GetResourceVersion(),
			}).Info("Add")
			queue.Add(key)
		},
		UpdateFunc: func(old, new interface{}) {
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
				"kind": kind,
				"key": key,
				"uid": metaOld.GetUID(),
				"oldVersion": metaOld.GetResourceVersion(),
				"newVersion": metaNew.GetResourceVersion(),
			}).Info("Update")
			queue.Add(key)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				log.WithError(err).Fatal("DeletionHandlingMetaNamespaceKeyFunc")
			}
			meta, err := meta.Accessor(obj)
			if err != nil {
				log.WithError(err).Fatal("meta.Accessor")
			}
			log.WithFields(log.Fields{
				"kind": kind,
				"key": key,
				"uid": meta.GetUID(),
			}).Info("Delete")
			queue.Add(key)
		},
	}
	return handler
}
