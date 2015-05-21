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

package registry

import (
	"fmt"
	"testing"
)

// An overflow concern relevant to the code and tests
// is clarified at the StackOverflow question at
// http://stackoverflow.com/questions/29880038

// Test the creation of a registry with misc bounds
// and flipped boundaries
func TestNew(t *testing.T) {

	boundsSet := [][2]uint16{
		{0, 0}, {2, 2}, {0, 2}, {0, 65535},
		{65534, 65535}, {65535, 65535},
	}

	for _, bounds := range boundsSet {
		_, err := New(bounds[0], bounds[1])
		if err != nil {
			t.Error(err)
		}
	}

	_, err := New(1, 0)
	if err == nil {
		t.Error("Registry creation should fail with min=1 and max=0")
	}
}

// Test port allocation at the high end - catch potential overflow bugs
func TestAddAllocHigh(t *testing.T) {

	const (
		max uint16 = ^uint16(0)
		min uint16 = max - 2
	)

	reg, err := New(min, max)
	if err != nil {
		t.Error(err)
	}

	for i := min; i <= max; i++ {
		p, err := reg.Alloc(fmt.Sprintf("svc%d", i), fmt.Sprintf("%d.%d.%d.%d", i, i, i, i))
		if err != nil {
			t.Errorf("Alloc registration # %d should not fail in [%d, %d] registry",
				i, min, max)
		}

		if p != i {
			t.Errorf("Allocated port %d should have been same as the iterator %d", p, i)
		}

		if i == max {
			break
		}
	}
}

type action int

const (
	add action = iota
	del
	chk
)

// Test generic use cases of allocating, querying,
// and forgetting port assignments
func TestRegistry(t *testing.T) {
	const (
		min uint16 = 0
		max uint16 = 3
	)

	reg, err := New(min, max)
	if err != nil {
		t.Errorf("Failed to create a registry with min=%d and max=%d", min, max)
	}

	mocks := [][]struct {
		act  action
		name string
		port uint16
		ok   bool
	}{
		{
			{act: add, name: "svc 3", port: 0, ok: true},
			{act: add, name: "svc 6", port: 1, ok: true},
			{act: add, name: "svc 3", port: 0, ok: false},
			{act: add, name: "svc 1", port: 2, ok: true},
			{act: add, name: "svc 9", port: 3, ok: true},
			{act: add, name: "extra", port: 0, ok: false},
		},
		{
			{act: del, name: "svc 6", port: 1, ok: true},
			{act: del, name: "svc 9", port: 3, ok: true},
			{act: chk, name: "svc 6", port: 1, ok: false},
			{act: chk, name: "svc 9", port: 3, ok: false},
		},
		{
			{act: add, name: "svc 9", port: 1, ok: true},
			{act: add, name: "svc 1", port: 0, ok: false},
			{act: add, name: "svc 2", port: 3, ok: true},
		},
		{
			{act: chk, name: "svc 3", port: 0, ok: true},
			{act: chk, name: "svc 1", port: 1, ok: true},
			{act: chk, name: "svc 1", port: 2, ok: true},
			{act: chk, name: "svc 2", port: 3, ok: true},
		},
	}

	for runid, run := range mocks {

		for line, mock := range run {

			switch mock.act {

			case add:

				p, err := reg.Alloc(mock.name)
				if (err == nil) != mock.ok {
					if err == nil {
						t.Errorf("Run %d mock %d: got port %d unexpectedly", runid, line, p)
					} else {
						t.Errorf("Run %d mock %d: %q", runid, line, err.Error())
					}
					continue
				}

			case del:
				reg.Forget(mock.port)

			case chk:

				p, _, err := reg.Lookup(mock.name)
				if (err == nil) != mock.ok {
					if err == nil {
						t.Errorf("Run %d mock %d: got port %d unexpectedly", runid, line, p)
					} else {
						t.Errorf("Run %d mock %d: %q", runid, line, err.Error())
					}
					continue
				}
			}
		}
	}

	pmax := ^uint16(0)
	for p, next := uint16(max+1), 0 < pmax; next; p, next = p+1, p < pmax {
		if svc, sOk := reg.byport[p]; sOk {
			t.Errorf("Unexpected service registered at port %d as %q", p, svc.name)
		}
	}
}

func matches(r *Registry, name string, port uint16) error {
	p, _, err := r.Lookup(name)
	if err != nil {
		return err
	}
	if p != port {
		return fmt.Errorf("Looking up %q returned %d instead of expected %d",
			name, p, port)
	}
	return nil
}

// Benchmark of allocation. Hopefully no software will allocate all
// 64K ports as this will get ugly (hash map lookups and mallocs)
func BenchmarkAddAlloc(b *testing.B) {

	const (
		min uint16 = 0
		max uint16 = ^uint16(0)
	)

	reg, err := New(min, max)
	if err != nil {
		b.Error(err)
	}

	for i := uint16(0); i < uint16(b.N); i++ {

		_, err := reg.Alloc(fmt.Sprintf("service %d", i))

		if err != nil {
			b.Errorf("Alloc registration %d should not fail in [%d, %d] registry",
				i, min, max)
		}

	}
}
