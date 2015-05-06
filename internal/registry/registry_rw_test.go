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
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestReService(t *testing.T) {
	mocks := []struct {
		line     string
		expected []string
	}{
		{
			`pald  1001  127.0.0.1,::1,not-local.example.com# Port Allocator Daemon`,
			[]string{
				"pald  1001  127.0.0.1,::1,not-local.example.com# Port Allocator Daemon",
				"pald",
				"1001",
				"127.0.0.1,::1,not-local.example.com",
				" Port Allocator Daemon",
			},
		},
		{
			` pald  1001`,
			[]string{
				" pald  1001",
				"pald",
				"1001",
				"",
				"",
			},
		},
		{
			` p_a.l-d  01001`,
			[]string{
				" p_a.l-d  01001",
				"p_a.l-d",
				"01001",
				"",
				"",
			},
		},
		{`pald`, nil},
		{`  1001  pald  `, nil},
		{
			`pald  1001 # Port Allocator Daemon`,
			[]string{
				"pald  1001 # Port Allocator Daemon",
				"pald",
				"1001",
				"",
				" Port Allocator Daemon",
			},
		},
		{
			`  # Port Allocator Daemon`,
			[]string{"  # Port Allocator Daemon", "", "", "", " Port Allocator Daemon"},
		},
		{
			`# Port Allocator Daemon`,
			[]string{"# Port Allocator Daemon", "", "", "", " Port Allocator Daemon"},
		},
		{``, []string{"", "", "", "", ""}},
	}

	for _, mock := range mocks {
		if stringSlicesDiffer(reService.FindStringSubmatch(mock.line), mock.expected) {
			t.Errorf("The line failed service parsing: %q", mock.line)
		}
	}
}

func stringSlicesDiffer(a, b []string) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}
	if len(a) != len(b) {
		return true
	}
	for i, e := range a {
		if e != b[i] {
			return true
		}
	}
	return false
}

func TestWrite(t *testing.T) {
	write_dst := "./reg_write.test"

	defer os.Remove(write_dst)
	reg := mockRegistry(t)

	dst, err := os.Create(write_dst)
	if err != nil {
		t.Error(err)
	}

	reg.Write(dst)

	result, err := ioutil.ReadFile(write_dst)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result, mockText) {
		t.Errorf("Expected and written differ. Written: %s", result)
	}
}

func TestRead(t *testing.T) {
	t.Error("Test not implemented")
}

var mockText = []byte(`fixed_0	0	0.0.0.0
alloc_1	1	0.0.0.0,1.1.1.1
fixed_2	2	
alloc_3	3	0.0.0.0,1.1.1.1,2.2.2.2
fixed_9	9	0.0.0.0,1.1.1.1,2.2.2.2,3.3.3.3
`)

func mockRegistry(t *testing.T) *Registry {

	mocks := []struct {
		port uint16
		name string
	}{
		{port: 2, name: "fixed"},
		{port: 0, name: "fixed"},
		{port: 1, name: "alloc"},
		{port: 3, name: "alloc"},
		{port: 9, name: "fixed"},
	}

	reg, err := New(1, 10)
	if err != nil {
		t.Fatal("Failed to create a mock registry")
	}

	for m, mock := range mocks {

		name := fmt.Sprintf("%s_%d", mock.name, mock.port)

		addr := make([]string, m)
		for i := 0; i < m; i++ {
			addr[i] = fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)
		}

		if mock.name == "fixed" {
			err = reg.Fix(mock.port, name, addr...)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			_, err = reg.Alloc(name, addr...)
			if err != nil {
				t.Error(err)
			}
		}
	}

	return reg
}
