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

const (
	// HostIDAnnotation annotates MidoNet Host ID for the Node.
	HostIDAnnotation = "midonet.org/host-id"

	// TunnelZoneIDAnnotation annotates MidoNet Tunnel Zone ID for the Node.
	// An empty string means the auto-created default tunnel zone.
	TunnelZoneIDAnnotation = "midonet.org/tunnel-zone-id"

	// TunnelEndpointIPAnnotation annotates MidoNet Tunnel Endpoint IP for
	// the Node.
	TunnelEndpointIPAnnotation = "midonet.org/tunnel-endpoint-ip"

	// MACAnnotation annotates MAC address for the Pod/Node.
	MACAnnotation = "midonet.org/mac-address"
)
