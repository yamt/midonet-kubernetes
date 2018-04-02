package midonet

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Push(resources []*APIResource, config *Config) error {
	for _, res := range resources {
		err := Post(res, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func Post(resource *APIResource, config *Config) error {
	data, err := json.Marshal(resource.Body)
	if err != nil {
		return err
	}
	url := config.API + resource.PathForPost
	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", resource.MediaType)
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"request": request,
		"url": url,
		"request-json": string(data),
		"response": response,
	}).Info("Do")
	return nil
}