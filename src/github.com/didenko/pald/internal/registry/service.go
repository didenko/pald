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
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type service struct {
	port uint16
	name string
	addr []string
}

func (s *service) equal(sr *service) bool {
	if s.port != sr.port ||
		s.name != sr.name ||
		!reflect.DeepEqual(s.addr, sr.addr) {

		return false
	}
	return true
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

	addr := compact(strings.Split(fields[3], ","))
	if len(addr) == 0 {
		addr = nil
	}

	return &service{
		uint16(port64),
		fields[1],
		addr,
	}, nil
}
