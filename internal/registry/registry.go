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

package registry

import (
	"fmt"
	"sync"
)

type Registry struct {
	sync.RWMutex
	byname   map[string]*service
	byport   map[uint16]*service
	portMin  uint16
	portMax  uint16
	portNext uint16
}

// Create New port registry with given boundaries
func New(min, max uint16) (*Registry, error) {
	if min > max {
		return nil,
			fmt.Errorf("Minimum port %d must be less than maximum port %d",
				min, max)
	}

	return &Registry{
		byname:   make(map[string]*service, 100),
		byport:   make(map[uint16]*service, 100),
		portMin:  min,
		portMax:  max,
		portNext: min,
	}, nil
}

// Lookup service details by it's symbolic name
func (r *Registry) Lookup(name string) (uint16, []string, error) {

	r.RLock()
	svc, ok := r.byname[name]
	r.RUnlock()

	if !ok {
		return 0, nil, fmt.Errorf("Name %q not found in the port registry", name)
	}

	return svc.port, svc.addr, nil
}

// Alloc attempts registering a service by it's symbolic name
// into a dynamically found port number. The function fails if
// the symbolic name is already registered or no more dynamic
// ports available in the pool
func (r *Registry) Alloc(name string, addr ...string) (uint16, error) {

	r.Lock()
	defer r.Unlock()

	_, name_taken := r.byname[name]

	if name_taken {
		return 0, fmt.Errorf("Name %q is already taken", name)
	}

	port, err := r.portFind()

	if err != nil {
		return 0, err
	}

	r.setSvc(port, name, addr...)

	return port, nil
}

// Forget removes service associated with the specified port.
// If the port is not in the registry, no error generated.
func (r *Registry) Forget(port uint16) {

	r.Lock()
	defer r.Unlock()

	if svc, ok := r.byport[port]; ok {
		delete(r.byname, svc.name)
	}

	delete(r.byport, port)

	if port < r.portNext {
		r.portNext = port
	}
}

// Fix records a static port mapipng to a symbolic name. It errors
// if either port or service are already in the registry. Fix allows
// to register a port at any unassigned number, even outside the
// registry's min-max range for dynamic allocations.
func (r *Registry) Fix(port uint16, name string, addr ...string) error {

	r.Lock()
	defer r.Unlock()

	_, name_taken := r.byname[name]
	_, port_taken := r.byport[port]

	if name_taken && port_taken {
		return fmt.Errorf("Name %q and port %d are already registered", name, port)
	}

	if name_taken {
		return fmt.Errorf("Name %q is already taken", name)
	}

	if port_taken {
		return fmt.Errorf("Port %d is already taken", port)
	}

	r.setSvc(port, name, addr...)

	return nil
}

type service struct {
	port uint16
	name string
	addr []string
}

func (r *Registry) portFind() (uint16, error) {

	for p, next := r.portNext, r.portNext <= r.portMax; next; p, next = p+1, p < r.portMax {

		_, taken := r.byport[p]

		if !taken {
			return p, nil
		}
	}

	return 0, fmt.Errorf("No ports available")
}

func (r *Registry) setSvc(port uint16, name string, addr ...string) {
	svc := &service{port, name, addr}
	r.byname[name] = svc
	r.byport[port] = svc
}
