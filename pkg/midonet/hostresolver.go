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
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// REVISIT(yamamoto): Consider to cache the mapping
// REVISIT(yamamoto): Consider making this a separate controller to
// annotate nodes with MidoNet Host IDs

type HostResolver struct {
	client *Client
}

func NewHostResolver(client *Client) *HostResolver {
	return &HostResolver{client}
}

func (h *HostResolver) ResolveHost(hostname string) (*uuid.UUID, error) {
	clog := log.WithField("hostname", hostname)
	clog.Debug("Start resolving")
	hosts, err := listHosts(h.client)
	if err != nil {
		clog.WithError(err).Error("listHosts")
		return nil, err
	}
	clog.WithField("hosts", hosts).Debug("Got hosts")
	for _, h := range hosts {
		if h.Name == hostname {
			clog.WithField("ID", h.ID).Info("Resolved")
			return h.ID, nil
		}
	}
	clog.Info("No host found")
	return nil, nil
}

func listHosts(c *Client) ([]Host, error) {
	var hosts []Host
	_, err := c.List(&hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}
