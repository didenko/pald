/*
	(c) 2015 Vlad Didenko. All Rights Reserved.

	This file is part of Port ALLocator Daemon a.k.a. pald.

	pald is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	pald is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with pald.  If not, see <http://www.gnu.org/licenses/>.
*/

package store

import (
	"fmt"
	"testing"

	"github.com/didenko/pald/internal/registry"
)

func TestSave(t *testing.T) {
	_ = mockRegistry(t)
	// reg := mockRegistry(t)

	t.Error("Test not implemented")
}

func TestLoad(t *testing.T) {
	t.Error("Test not implemented")
}

func mockRegistry(t *testing.T) *registry.Registry {

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

	reg, err := registry.New(1, 10)
	if err != nil {
		t.Fatal("Failed to create a mock registry")
	}

	for _, mock := range mocks {

		name := fmt.Sprintf("%s %d", mock.name, mock.port)

		if mock.name == "fixed" {
			err = reg.Fix(mock.port, name)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			_, err = reg.Alloc(name)
			if err != nil {
				t.Error(err)
			}
		}
	}

	return reg
}
