// Copyright (C) 2018 Midokura SARL.
// All rights reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package midonet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

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
		if resp.StatusCode == 404 {
			if _, ok := res.(HasParent); ok {
				log.Info("Parent doesn't exist yet?")
				return fmt.Errorf("Parent doesn't exist yet?")
			}
		}
		if resp.StatusCode == 409 {
			if res.Path("PUT") != "" {
				resp, err = c.put(res)
				if err != nil {
					return err
				}
			} else {
				// assume 409 meant ok
				continue
			}
		}
		if resp.StatusCode/100 != 2 {
			log.WithField("statusCode", resp.StatusCode).Fatal("Unexpected status code")
		}
	}
	return nil
}

func (c *Client) Delete(resources []APIResource) error {
	for _, res := range resources {
		resp, _, err := c.doRequest("DELETE", res.Path("DELETE"), nil, "")
		if err != nil {
			return err
		}
		// Ignore 404 assuming it's ok.  Even if we're the only one making
		// MidoNet topology modifications, it happens e.g. when a removal
		// of a Chain cascade-deleted Rules.
		if resp.StatusCode/100 != 2 && resp.StatusCode != 404 {
			log.WithField("statusCode", resp.StatusCode).Fatal("Unexpected status code")
		}
	}
	return nil
}

func (c *Client) post(res APIResource) (*http.Response, error) {
	resp, _, err := c.doRequest("POST", res.Path("POST"), res, "")
	return resp, err
}

func (c *Client) put(res APIResource) (*http.Response, error) {
	resp, _, err := c.doRequest("PUT", res.Path("PUT"), res, "")
	return resp, err
}

func (c *Client) List(rs interface{}) (*http.Response, error) {
	// assumption: rs is a pointer to an array of ListableResource
	// E.g. *[]Host
	t := reflect.TypeOf(rs)
	et := t.Elem().Elem()
	p := reflect.New(et)
	r := p.Interface().(ListableResource)
	resp, body, err := c.doRequest("GET", r.Path("LIST"), nil, r.CollectionMediaType())
	if err != nil {
		return resp, err
	}
	dec := json.NewDecoder(strings.NewReader(body))
	err = dec.Decode(rs)
	return resp, err
}

func (c *Client) doRequest(method string, path string, res APIResource, respType string) (*http.Response, string, error) {
	url := c.config.API + path
	clog := log.WithFields(log.Fields{
		"method": method,
		"url":    url,
	})
	var body io.Reader
	if res != nil {
		data, err := json.Marshal(res)
		if err != nil {
			return nil, "", err
		}
		body = bytes.NewReader(data)
		clog = clog.WithField("requestBody", string(data))
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, "", err
	}
	if res != nil {
		req.Header.Add("Content-Type", res.MediaType())
	}
	if respType != "" {
		req.Header.Add("Accept", respType)
	}

	// TODO: login
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	clog = clog.WithFields(log.Fields{
		"statusCode":   resp.StatusCode,
		"responseBody": string(respBody),
	})
	clog.Info("Do")
	return resp, string(respBody), nil
}
