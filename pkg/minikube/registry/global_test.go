/*
Copyright 2018 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGlobalRegister(t *testing.T) {
	globalRegistry = newRegistry()
	foo := DriverDef{Name: "foo"}
	if err := Register(foo); err != nil {
		t.Errorf("Register = %v, expected nil", err)
	}
	if err := Register(foo); err == nil {
		t.Errorf("Register = nil, expected duplicate err")
	}
}

func TestGlobalDriver(t *testing.T) {
	foo := DriverDef{Name: "foo"}
	globalRegistry = newRegistry()

	if err := Register(foo); err != nil {
		t.Errorf("Register = %v, expected nil", err)
	}

	d := Driver("foo")
	if d.Empty() {
		t.Errorf("driver.Empty = true, expected false")
	}

	d = Driver("bar")
	if !d.Empty() {
		t.Errorf("driver.Empty = false, expected true")
	}
}

func TestGlobalList(t *testing.T) {
	foo := DriverDef{Name: "foo"}
	globalRegistry = newRegistry()
	if err := Register(foo); err != nil {
		t.Errorf("register returned error: %v", err)
	}

	if diff := cmp.Diff(List(), []DriverDef{foo}); diff != "" {
		t.Errorf("list mismatch (-want +got):\n%s", diff)
	}
}

func TestGlobalInstalled(t *testing.T) {
	globalRegistry = newRegistry()

	if err := Register(DriverDef{Name: "foo"}); err != nil {
		t.Errorf("register returned error: %v", err)
	}

	bar := DriverDef{
		Name:     "bar",
		Priority: Default,
		Status:   func() State { return State{Installed: true} },
	}
	if err := Register(bar); err != nil {
		t.Errorf("register returned error: %v", err)
	}

	expected := []DriverState{
		DriverState{
			Name:     "bar",
			Priority: Default,
			State: State{
				Installed: true,
			},
		},
	}

	if diff := cmp.Diff(Installed(), expected); diff != "" {
		t.Errorf("installed mismatch (-want +got):\n%s", diff)
	}
}

func TestGlobalStatus(t *testing.T) {
	globalRegistry = newRegistry()

	if err := Register(DriverDef{Name: "foo"}); err != nil {
		t.Errorf("register returned error: %v", err)
	}

	expected := State{Installed: true, Healthy: true}
	bar := DriverDef{
		Name:     "bar",
		Priority: Default,
		Status:   func() State { return expected },
	}
	if err := Register(bar); err != nil {
		t.Errorf("register returned error: %v", err)
	}

	if diff := cmp.Diff(Status("bar"), expected); diff != "" {
		t.Errorf("status mismatch (-want +got):\n%s", diff)
	}
}
