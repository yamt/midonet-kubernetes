package midonet

import (
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	converter Converter
	config    *Config
}

func NewHandler(converter Converter, config *Config) *Handler {
	return &Handler{converter, config}
}

func (h *Handler) Update(key string, obj interface{}) error {
	clog := log.WithFields(log.Fields{
		"key": key,
		"obj": obj,
	})
	converted, err := h.converter.Convert(key, obj, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	cli := NewClient(h.config)
	err = cli.Push(converted)
	if err != nil {
		clog.WithError(err).Fatal("Failed to push")
	}
	return nil
}

func (h *Handler) Delete(key string) error {
	clog := log.WithField("key", key)
	converted, err := h.converter.Convert(key, nil, h.config)
	if err != nil {
		clog.WithError(err).Fatal("Failed to convert")
	}
	clog.WithField("converted", converted).Info("Converted")
	return nil
}
