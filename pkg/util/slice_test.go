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

package util

import (
	"reflect"
	"testing"
)

func TestRemoveFirst(t *testing.T) {
	xs := []string{
		"a", "b", "c", "d", "a",
	}
	ys, removed := RemoveFirst(xs, "b")
	if !removed {
		t.FailNow()
	}
	expected := []string{"a", "c", "d", "a"}
	if !reflect.DeepEqual(ys, expected) {
		t.Errorf("expected %v actual %v", expected, ys)
		t.FailNow()
	}
	ys, removed = RemoveFirst(ys, "b")
	if removed {
		t.FailNow()
	}
	expected = []string{"a", "c", "d", "a"}
	if !reflect.DeepEqual(ys, expected) {
		t.Errorf("expected %v actual %v", expected, ys)
		t.FailNow()
	}
	ys, removed = RemoveFirst(ys, "a")
	if !removed {
		t.FailNow()
	}
	expected = []string{"c", "d", "a"}
	if !reflect.DeepEqual(ys, expected) {
		t.Errorf("expected %v actual %v", expected, ys)
		t.FailNow()
	}
	ys, removed = RemoveFirst(ys, "a")
	if !removed {
		t.FailNow()
	}
	expected = []string{"c", "d"}
	if !reflect.DeepEqual(ys, expected) {
		t.Errorf("expected %v actual %v", expected, ys)
		t.FailNow()
	}
}
