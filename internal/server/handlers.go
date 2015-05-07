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

package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func get(w http.ResponseWriter, r *http.Request) {

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service := r.Form.Get("service")

	if service == "" {
		http.Error(w, "Service name is missing", http.StatusBadRequest)
		return
	}

	port, _, err := reg.Lookup(service)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "%d\n", port)
}

func set(w http.ResponseWriter, r *http.Request) {

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service := r.Form.Get("service")

	if service == "" {
		http.Error(w, "Service name is missing", http.StatusBadRequest)
		return
	}

	port, err := reg.Alloc(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "%d\n", port)
}

func del(w http.ResponseWriter, r *http.Request) {

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	portStr := r.Form.Get("port")

	if portStr == "" {
		http.Error(w, "Port number is missing", http.StatusBadRequest)
		return
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reg.Forget(uint16(port))

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintln(w, "OK")
}
