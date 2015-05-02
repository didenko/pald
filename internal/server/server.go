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
	"log"
	"net/http"

	"github.com/didenko/pald/internal/registry"
)

var reg, err = registry.New(49200, 49299)

func init() {
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
	http.HandleFunc("/del", del)

}

func Run(port uint16) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
