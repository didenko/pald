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

// Test generic use cases of fixing, allocating, querying,
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

	mocks := []struct {
		port  uint16
		name  string
		addok bool
		del   bool
	}{
		{port: 2, name: "fixed", addok: true, del: false},
		{port: 0, name: "fixed", addok: true, del: false},
		{port: 1, name: "alloc", addok: true, del: true},
		{port: 3, name: "alloc", addok: true, del: false},
		{port: 1, name: "fixed", addok: false, del: true},
		{port: 4, name: "alloc", addok: false, del: false},
		{port: 9, name: "fixed", addok: true, del: true},
	}

	mockMap := make(map[uint16]int)

	for _, mock := range mocks {

		name := fmt.Sprintf("%s %d", mock.name, mock.port)

		if mock.name == "fixed" {

			err = reg.Fix(mock.port, name)
			if (err == nil) != mock.addok {
				t.Error(err)
				continue
			}

		} else {

			_, err = reg.Alloc(name)
			if (err == nil) != mock.addok {
				t.Error(err)
				continue
			}
		}
	}

	for line, mock := range mocks {
		if mock.del {
			reg.Forget(mock.port)
		} else {
			if mock.addok {
				mockMap[mock.port] = line
			}
		}
	}

	pmax := ^uint16(0)
	for p, next := uint16(0), 0 < pmax; next; p, next = p+1, p < pmax {

		midx, mOk := mockMap[p]

		svc, sOk := reg.byport[p]

		if sOk == mOk {
			if mOk && svc.name != fmt.Sprintf("%s %d", mocks[midx].name, mocks[midx].port) {
				t.Errorf("Service names do not match for port %d", p)
				if err := matches(reg, mocks[midx].name, mocks[midx].port); err != nil {
					t.Error(err)
				}
			}

		} else if sOk {
			t.Errorf("Unexpected service registered at port %d as %q", p, svc.name)

		} else if mocks[midx].addok && !mocks[midx].del {
			t.Errorf("Missing service expected at port %d as %q", p, mocks[midx].name)
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
