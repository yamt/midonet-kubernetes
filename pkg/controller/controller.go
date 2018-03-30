package controller

import (
	log "github.com/sirupsen/logrus"

	// kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Handler interface {
	Handle(interface{}) error
}

type Controller struct {
	informerFactory informers.SharedInformerFactory
	queue workqueue.RateLimitingInterface
	handler Handler
}

func NewController(si informers.SharedInformerFactory, queue workqueue.RateLimitingInterface, handler Handler) *Controller {
	return &Controller{
		informerFactory: si,
		queue: queue,
		handler: handler,
	}
}

func (c *Controller) startInformer(stopCh <-chan struct{}) {
	c.informerFactory.Start(stopCh)
	c.informerFactory.WaitForCacheSync(stopCh)
	log.Info("Cache synced")
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	c.startInformer(stopCh)

	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	queue := c.queue
	key, quit := queue.Get()
	if quit {
		return false
	}
	defer queue.Done(key)
	clog := log.WithFields(log.Fields{
		"key": key,
	})

	clog.Info("Start processing")
	err := c.handler.Handle(key)
	if err == nil {
		clog.Info("Done")
		queue.Forget(key)
		return true
	}
	clog.WithError(err).Error("Failed")
	queue.AddRateLimited(key)
	return true
}

func AddHandler(informer cache.SharedIndexInformer, kind string) workqueue.RateLimitingInterface {
	rateLimiter := workqueue.DefaultControllerRateLimiter()
	queue := workqueue.NewNamedRateLimitingQueue(rateLimiter, kind)
	handler := newHandler(kind, queue)
	informer.AddEventHandler(handler)
	return queue
}

func newHandler(kind string, queue workqueue.Interface) cache.ResourceEventHandler {
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			logAndQueue("Add", kind, queue, obj, nil)
		},
		UpdateFunc: func(old, new interface{}) {
			logAndQueue("Update", kind, queue, new, old)
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
