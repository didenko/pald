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
	"regexp"
	"strconv"
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

// Write writes out all registry's services
func (r *Registry) Dump(w io.Writer) (err error) {

	r.RLock()
	defer r.RUnlock()

	buf := bufio.NewWriter(w)
	min, max := uint16(0), ^uint16(0)

	for p, next := min, min < max; next; p, next = p+1, p < max {

		if s, ok := r.byport[p]; ok {

			_, err = fmt.Fprintf(buf, "%s\t%d\t%s\n", s.name, s.port, strings.Join(s.addr, ","))
			if err != nil {
				return err
			}
		}
	}
	return buf.Flush()
}

var reService = regexp.MustCompile(`^(?:\s*(?P<name>[\w\-\.]+)\s+(?P<port>\d+)(?:\s+(?P<addr>[,\.\-_:\w]+))?)?\s*(?:\#(?P<comment>.*))?$`)

func parseSvc(line string) (*service, error) {

	compact := func(a []string) []string {
		b := make([]string, 0, len(a)/2)
		for _, e := range a {
			if e != "" {
				b = append(b, e)
			}
		}
		return b
	}

	// o:line, 1:name, 2:port, 3:addr, 4:comment
	fields := reService.FindStringSubmatch(line)
	if fields == nil {
		return nil, fmt.Errorf("The line fails to match a service definition: %q", line)
	}

	port64, err := strconv.ParseUint(fields[2], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Port spec %q failes to parse into 16-bit unsigned integer: %s", fields[2], err.Error())
	}

	return &service{
		uint16(port64),
		fields[1],
		compact(strings.Split(fields[3], ",")),
	}, nil
}

// Read reads all the services from r.
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
