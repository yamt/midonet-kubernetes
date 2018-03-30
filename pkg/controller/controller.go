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
			logAndQueue("Delete", kind, queue, obj, nil)
		},
		UpdateFunc: func(old, new interface{}) {
			logAndQueue("Delete", kind, queue, new, old)
		},
		DeleteFunc: func(obj interface{}) {
			logAndQueue("Delete", kind, queue, obj, nil)
		},
	}
	return handler
}

func logAndQueue(logMsg string, kind string, queue workqueue.Interface, obj interface{}, oldObj interface{}) error {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.WithError(err).Fatal("DeletionHandlingMetaNamespaceKeyFunc")
	}
	ctxLog := log.NewEntry(log.StandardLogger())
	if _, ok := obj.(cache.DeletedFinalStateUnknown); !ok {
		meta, err := meta.Accessor(obj)
		if err != nil {
			ctxLog.WithError(err).Fatal("meta.Accessor")
		} else {
			ctxLog = ctxLog.WithFields(log.Fields{
				"uid": meta.GetUID(),
				"version": meta.GetResourceVersion(),
			})
		}
	}
	if oldObj != nil {
		metaOld, err := meta.Accessor(oldObj)
		if err != nil {
			ctxLog.WithError(err).Fatal("meta.Accessor")
		}
		ctxLog = ctxLog.WithFields(log.Fields{
			"oldVersion": metaOld.GetResourceVersion(),
		})
	}
	ctxLog = ctxLog.WithFields(log.Fields{
		"kind": kind,
		"key": key,
	})
	ctxLog.Info(logMsg)
	queue.Add(key)
	return nil
}
