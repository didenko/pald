/*
	(c) Copyright 2015 Vlad Didenko

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

package persist

import (
	"testing"
	"time"

	"github.com/didenko/pald/internal/registry"
)

func TestPersist(t *testing.T) {
	r, _ := registry.New(0, 10)
	s := new(StringRWST)
	d := 100 * time.Millisecond
	var b struct{}

	flush := Persist(r, s, d)

	if s.Get() != "" {
		t.Error("Persistence test setup failed")
	}

	_, _ = r.Alloc("svc_0")
	flush <- b

	if s.Get() != "" {
		t.Error("Flushing in under timeout should have not produce data")
	}

	time.Sleep(200 * time.Millisecond) // wait for throttle to expire
	flush <- b
	time.Sleep(100 * time.Millisecond) // wait for Dump to happen in another goroutine

	if w := s.Get(); w != "svc_0\t0\t\n" {
		t.Errorf("A wrong string written. Received %q\n", w)
	}

	flush <- b
	_, _ = r.Alloc("svc_1", "127.0.0.1", "::1")
	flush <- b

	if w := s.Get(); w != "svc_0\t0\t\n" {
		t.Errorf("Flushing in under timeout should have not change data. Received %q\n", w)
	}

	time.Sleep(200 * time.Millisecond) // wait for throttle to expire
	flush <- b
	time.Sleep(100 * time.Millisecond) // wait for Dump to happen in another goroutine

	if w := s.Get(); w != "svc_0\t0\t\nsvc_1\t1\t127.0.0.1,::1\n" {
		t.Errorf("A wrong string written. Received %q\n", w)
	}
}

func TestLoad(t *testing.T) {
	s := new(StringRWST)
	s.Set("svc_0\t0\t\nsvc_1\t1\t127.0.0.1,::1\n")

	rBuild, _ := registry.New(0, 10)
	_, _ = rBuild.Alloc("svc_0")
	_, _ = rBuild.Alloc("svc_1", "127.0.0.1", "::1")

	rLoaded, _ := registry.New(0, 10)
	if err := Load(s, rLoaded); err != nil {
		t.Fatal(err)
	}

	if !rBuild.Equal(rLoaded) {
		t.Error("Build and loaded registries differ.\n", rBuild, "\n", rLoaded)
	}
}
