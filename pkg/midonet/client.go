package midonet

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Post(resource *MidoNetAPIResource, config *Config) error {
	data, err := json.Marshal(resource.Body)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", config.API + resource.PathForPost, bytes.NewReader(data))
	if err != nil {
		return err
	}
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"response": response,
	}).Info("Response")
	return nil
}
