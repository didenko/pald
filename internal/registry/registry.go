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
	byname    map[string]*service
	byport    map[uint16]*service
	port_min  uint16
	port_max  uint16
	port_next uint16
}

func New(min, max uint16) (*Registry, error) {
	if min > max {
		return nil,
			fmt.Errorf("Minimum port %d must be less than maximum port %d",
				min, max)
	}

	return &Registry{
		byname:    make(map[string]*service, 100),
		byport:    make(map[uint16]*service, 100),
		port_min:  min,
		port_max:  max,
		port_next: min,
	}, nil
}

func (r *Registry) Query(name string) (uint16, []string, error) {

	r.RLock()
	svc, ok := r.byname[name]
	r.RUnlock()

	if !ok {
		return 0, nil, fmt.Errorf("Name \"%q\" not found in the port registry", name)
	}

	return svc.port, svc.addr, nil
}

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

func (r *Registry) Forget(port uint16) {

	r.Lock()
	defer r.Unlock()

	delete(r.byport, port)

	if svc, ok := r.byport[port]; ok {
		delete(r.byname, svc.name)
	}

	if port < r.port_next {
		r.port_next = port
	}
}

func (r *Registry) Fixed(port uint16, name string, addr ...string) error {

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

	for p, next := r.port_next, r.port_next <= r.port_max; next; p, next = p+1, p < r.port_max {

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
