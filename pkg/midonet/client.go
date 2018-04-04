package midonet

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	return &Client{config}
}

func (c *Client) Push(resources []APIResource) error {
	for _, res := range resources {
		// REVISIT: maybe we should save updates (and thus zk and
		// midolman loads) by performing GET and compare first.
		// Or we can make the MidoNet API detect and ignore no-op updates.
		resp, err := c.post(res)
		if err != nil {
			return err
		}
		if resp.StatusCode == 409 && res.Path("PUT") != "" {
			_, err := c.put(res)
			if err != nil {
				return err
			}
		}
		// TODO: check resp.StatusCode 
	}
	return nil
}

func (c *Client) Delete(resources []APIResource) error {
	for _, res := range resources {
		_, err := c.doRequest("DELETE", res.Path("DELETE"), nil)
		if err != nil {
			return err
		}
		// TODO: check resp.StatusCode 
	}
	return nil
}

func (c *Client) post(res APIResource) (*http.Response, error) {
	return c.doRequest("POST", res.Path("POST"), res)
}

func (c *Client) put(res APIResource) (*http.Response, error) {
	return c.doRequest("PUT", res.Path("PUT"), res)
}

func (c *Client) doRequest(method string, path string, res APIResource) (*http.Response, error) {
	url := c.config.API + path
	clog := log.WithField("url", url)
	var body io.Reader
	if res != nil {
		data, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
		clog = clog.WithField("request-json", string(data))
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if res != nil {
		req.Header.Add("Content-Type", res.MediaType())
	}
	clog = clog.WithField("request", req)

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
