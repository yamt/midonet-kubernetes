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
	"net"
	"reflect"
	"testing"
)

func TestMAC(t *testing.T) {
	actual := MACForKey("hey")
	expected, err := net.ParseMAC("ac:ca:ba:fa:69:0b")
	if err != nil {
		t.Errorf("ParseMAC")
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestIDForKeyDifferBetweenObjects(t *testing.T) {
	ID1 := IDForKey("Foo", "obj1", &Config{Tenant: "MyTenant"})
	ID2 := IDForKey("Foo", "obj2", &Config{Tenant: "MyTenant"})
	if ID1 == ID2 {
		t.Errorf("Got the same ID for different objects")
	}
}

func TestIDForKeyDifferBetweenTenants(t *testing.T) {
	myID := IDForKey("Pod", "hoge/fuga", &Config{Tenant: "MyTenant"})
	yourID := IDForKey("Pod", "hoge/fuga", &Config{Tenant: "YourTenant"})
	if myID == yourID {
		t.Errorf("Got the same ID for different tenants")
	}
}

func TestIDForTenantDifferBetweenTenants(t *testing.T) {
	myID := idForTenant("MyTenant")
	yourID := idForTenant("YourTenant")
	if myID == yourID {
		t.Errorf("Got the same ID for different tenants")
	}
}
