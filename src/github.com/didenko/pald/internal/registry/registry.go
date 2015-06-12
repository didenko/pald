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
	"bufio"
	"fmt"
	"io"
	"strings"
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

	r.createSvc(port, name, addr...)

	return port, nil
}

// Forget removes the service associated with the specified port.
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

// Dump writes out all registry's services
func (r *Registry) Dump(w io.Writer) (int, error) {

	r.RLock()
	defer r.RUnlock()

	var (
		wrote, n int = 0, 0
		err      error
		buf      *bufio.Writer = bufio.NewWriter(w)
		min      uint16        = 0
		max      uint16        = ^uint16(0)
	)

	for p, next := min, min < max; next; p, next = p+1, p < max {

		if s, ok := r.byport[p]; ok {

			n, err = fmt.Fprintf(buf, "%s\t%d\t%s\n", s.name, s.port, strings.Join(s.addr, ","))
			wrote += n
			if err != nil {
				return wrote, err
			}
		}
	}
	return wrote, buf.Flush()
}

// Load reads all the services from r.
func (reg *Registry) Load(r io.Reader) (err error) {

	reg.Lock()
	defer reg.Unlock()

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {

		service, err := parseSvc(scanner.Text())
		if err != nil {
			return err
		}

		reg.setSvc(service)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
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

func (r *Registry) setSvc(svc *service) error {
	// check valid service values here
	r.byname[svc.name] = svc
	r.byport[svc.port] = svc
	return nil
}

func (r *Registry) createSvc(port uint16, name string, addr ...string) error {
	svc := &service{port, name, addr}
	return r.setSvc(svc)
}

// Equal is used in testing only. It is made to be able to debug
// where the comparison actually fails. The reflect.DeepEqual does
// not provide enough details to find mismatches deep in the structure.
func (r *Registry) Equal(rr *Registry) bool {

	if r.portMin != rr.portMin ||
		r.portMax != rr.portMax ||
		r.portNext != rr.portNext ||
		len(r.byport) != len(rr.byport) ||
		len(r.byname) != len(rr.byname) {

		return false
	}

	for p, s := range r.byport {
		sr, ok := rr.byport[p]
		if !ok {
			return false
		}
		if !s.equal(sr) {
			return false
		}
	}

	for n, s := range r.byname {
		sr, ok := rr.byname[n]
		if !ok {
			return false
		}
		if !s.equal(sr) {
			return false
		}
	}

	return true
}
