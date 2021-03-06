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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/didenko/pald/internal/persist"
	"github.com/didenko/pald/internal/registry"
)

var (
	reg *registry.Registry
	err error

	flusher chan struct{}
)

func init() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
	http.HandleFunc("/del", del)
}

func Run(portSvr, portMin, portMax uint16, dumpName string) {

	reg, err = registry.New(portMin, portMax)
	if err != nil {
		log.Panic(err)
	}

	dump, err := os.OpenFile(dumpName, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		log.Panic(err)
	}

	err = persist.Load(dump, reg)
	if err != nil {
		log.Panic(err)
	}

	flusher = persist.Persist(reg, dump, time.Second)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portSvr), nil))
}
