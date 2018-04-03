package midonet

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	return &Client{config}
}

func (c *Client) Push(resources []*APIResource) error {
	for _, res := range resources {
		// REVISIT: maybe we should save updates (and thus zk and
		// midolman loads) by performing GET and compare first.
		// Or we can make the MidoNet API detect and ignore no-op updates.
		resp, err := c.post(res)
		if err != nil {
			return err
		}
		if resp.StatusCode == 409 && res.PathForPut != "" {
			_, err := c.put(res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: func (c *) Delete(resources []*APIResource) error

func (c *Client) post(res *APIResource) (*http.Response, error) {
	return c.postOrPut("POST", res.PathForPost, res)
}

func (c *Client) put(res *APIResource) (*http.Response, error) {
	return c.postOrPut("PUT", res.PathForPut, res)
}

func (c *Client) postOrPut(method string, path string, res *APIResource) (*http.Response, error) {
	data, err := json.Marshal(res.Body)
	if err != nil {
		return nil, err
	}
	url := c.config.API + path
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	clog := log.WithFields(log.Fields{
		"request":      req,
		"url":          url,
		"request-json": string(data),
	})
	req.Header.Add("Content-Type", res.MediaType)
	return c.doRequest(req, clog)
}

func (c *Client) doRequest(req *http.Request, clog *log.Entry) (*http.Response, error) {
	// TODO: login
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	clog.WithFields(log.Fields{
		"response": resp,
	}).Info("Do")
	return resp, nil
}
