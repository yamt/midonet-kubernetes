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

package converter

import (
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Space constants for Version 5 UUID.
const (
	// midonetTenantSpaceUUIDString is used to generate MidoNet UUID
	// (eg. PortID) for the MidoNet tenant.  It's used for global resources.
	midonetTenantSpaceUUIDString = "3978567E-91C4-465C-A0D1-67575F6B4C7F"
)

func idForString(space uuid.UUID, key string) uuid.UUID {
	return uuid.NewSHA1(space, []byte(fmt.Sprintf("%s/%s", TranslationVersion, key)))
}

func idForTenant(tenant string) uuid.UUID {
	space, err := uuid.Parse(midonetTenantSpaceUUIDString)
	if err != nil {
		log.WithError(err).Fatal("space")
	}
	return idForString(space, tenant)
}

// IDForKey deterministically generates a MidoNet UUID for a given Kubernetes
// resource key, that is "Namespace/Name" string.
// Note: This function generates different UUIDs for different Tenants.
// Note: This function is also (ab)used for our pseudo resources;
// ServicePort, ServicePortSub, and Endpoint.
func IDForKey(kind string, key string, config *Config) uuid.UUID {
	// kubernetesSpaceUUID is used to generate MidoNet UUID (eg. PortID)
	// for Kubernetes resouces (Namespace/Name)
	kubernetesSpaceUUID := idForTenant(config.Tenant)
	return idForString(kubernetesSpaceUUID, fmt.Sprintf("%s/%s", kind, key))
}

// SubID deterministically generates another MidoNet UUID for the resource
// identified by the given UUID.  It's used e.g. when more than two MidoNet
// resources (thus UUIDs) are necessary for a Kubernetes resource.
func SubID(id uuid.UUID, s string) uuid.UUID {
	return uuid.NewSHA1(id, []byte(s))
}

// MACForKey deterministically generates an MAC address from a given string,
// which is typically a Kubernetes resource key.
// Note: Don't assume this conflict-free as the output space is rather
// small. (only 24 bits)
func MACForKey(key string) net.HardwareAddr {
	hash := sha256.Sum256([]byte(key))
	// AC-CA-BA  Midokura Co., Ltd.
	addr := [6]byte{0xac, 0xca, 0xba, hash[0], hash[1], hash[2]}
	return net.HardwareAddr(addr[:])
}
