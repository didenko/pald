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
	portStr := r.Form.Get("port")

	if len(portStr) > 0 {

		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = reg.Fix(uint16(port), service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintln(w, "OK")

	} else {

		port, err := reg.Alloc(service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "%d\n", port)
	}
}

func del(w http.ResponseWriter, r *http.Request) {

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	portStr := r.Form.Get("port")

	if portStr == "" {
		http.Error(w, "Service name is missing", http.StatusBadRequest)
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
