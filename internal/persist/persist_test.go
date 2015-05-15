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
	"fmt"
	"testing"
	"time"

	"github.com/didenko/pald/internal/registry"
)

type StringWriter struct {
	s string
}

func (sw *StringWriter) Write(p []byte) (n int, err error) {
	sw.s = fmt.Sprintf("%s", p)
	return len(sw.s), nil
}

func (sw *StringWriter) Get() string {
	return sw.s
}

func TestPersist(t *testing.T) {
	r, _ := registry.New(0, 10)
	s := new(StringWriter)
	d := 100 * time.Millisecond
	var b struct{}

	flush := Persist(r, s, d)

	if s.Get() != "" {
		t.Error("Persistence test setup failed")
	}

	_, _ = r.Alloc("svc 0")
	flush <- b

	if s.Get() != "" {
		t.Error("Flushing in under timeout should have not produce data")
	}

	time.Sleep(200 * time.Millisecond) // wait for throttle to expire
	flush <- b
	time.Sleep(100 * time.Millisecond) // wait for Dump to happen in another goroutine

	if w := s.Get(); w != "svc 0\t0\t\n" {
		t.Errorf("A wrong string written. Received %q\n", w)
	}

	flush <- b
	_, _ = r.Alloc("svc 1", "127.0.0.1", "::1")
	flush <- b

	if w := s.Get(); w != "svc 0\t0\t\n" {
		t.Errorf("Flushing in under timeout should have not change data. Received %q\n", w)
	}

	time.Sleep(200 * time.Millisecond) // wait for throttle to expire
	flush <- b
	time.Sleep(100 * time.Millisecond) // wait for Dump to happen in another goroutine

	if w := s.Get(); w != "svc 0\t0\t\nsvc 1\t1\t127.0.0.1,::1\n" {
		t.Errorf("A wrong string written. Received %q\n", w)
	}
}
