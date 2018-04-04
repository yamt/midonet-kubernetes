package midonet

import (
	log "github.com/sirupsen/logrus"
)

type SubResource interface {
	Convert(key string) ([]APIResource, error)
}

type SubResourceMap map[string]SubResource

type Handler struct {
	converter Converter
	config    *Config
	knownSubResources map[string]map[string]SubResource
}

func NewHandler(converter Converter, config *Config) *Handler {
	return &Handler{converter, config, make(map[string]map[string]SubResource)}
}

func (h *Handler) deletedSubResources(key string, rs map[string]SubResource) map[string]SubResource {
	known := h.knownSubResources[key]
	deleted := make(map[string]SubResource)
	for k, r := range known {
		if _, ok := rs[k]; !ok {
			deleted[k] = r
		}
	}
	return deleted
}

func merge(a map[string]SubResource, b map[string]SubResource) map[string]SubResource {
	c := make(map[string]SubResource)
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

func (h *Handler) handleSubResources(key string, added map[string]SubResource, deleted map[string]SubResource, cli *Client, clog *log.Entry) error {
	if h.knownSubResources[key] == nil {
		h.knownSubResources[key] = make(map[string]SubResource)
	}
	for k, r := range merge(h.deletedSubResources(key, added), deleted) {
		converted, err := r.Convert(k)
		err = cli.Delete(converted)
		if err != nil {
			clog.WithError(err).WithFields(log.Fields{
				"key": key,
				"sub-key": k,
			}).Error("failed to delete a sub resource")
			return err
		}
		delete(h.knownSubResources[key], k)
	}
	for k, r := range added {
		convertedSub, err := r.Convert(k)
		err = cli.Push(convertedSub)
		// Remember the resource regardless of err as we might have
		// partially pushed.
		h.knownSubResources[key][k] = r
		if err != nil {
			clog.WithError(err).WithFields(log.Fields{
				"key": key,
				"sub-key": k,
			}).Error("failed to push a sub resource")
			return err
		}
	}
	return nil
}

func (h *Handler) Update(key string, obj interface{}) error {
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	converted, subResources, err := h.converter.Convert(key, obj, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	cli := NewClient(h.config)
	err = cli.Push(converted)
	if err != nil {
		clog.WithError(err).Error("Failed to push")
		return err
	}
	err = h.handleSubResources(key, subResources, nil, cli, clog)
	if err != nil {
		clog.WithError(err).Error("handleSubResources")
		return err
	}
	// TODO: annotate kubernetes obj
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	converted, subResources, err := h.converter.Convert(key, nil, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	cli := NewClient(h.config)
	err = h.handleSubResources(key, nil, subResources, cli, clog)
	if err != nil {
		clog.WithError(err).Error("handleSubResources")
		return err
	}
	err = cli.Delete(converted)
	if err != nil {
		clog.WithError(err).Error("Failed to delete")
		return err
	}
	return nil
}
