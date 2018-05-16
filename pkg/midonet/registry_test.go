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
	"reflect"
	"testing"
)

func TestTypeNameForObject(t *testing.T) {
	obj := Bridge{}
	actual := TypeNameForObject(obj)
	expected := "Bridge"
	if actual != expected {
		t.Errorf("got %v\nwant %v", string(actual), expected)
	}
}

func TestObjectByTypeName(t *testing.T) {
	obj := ObjectByTypeName("Bridge")
	actual, ok := obj.(*Bridge)
	if !ok {
		t.Errorf("wrong type %v", reflect.TypeOf(obj))
	}
	expected := &Bridge{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
	if *actual != *expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
	i := obj.(APIResource)
	actualType := i.MediaType()
	expectedType := "application/vnd.org.midonet.Bridge-v4+json"
	if actualType != expectedType {
		t.Errorf("got %v\nwant %v", actualType, expectedType)
	}
}
