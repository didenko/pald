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

package platform

import (
	"log"
	"os"
	u "os/user"
	"path"
)

type darwinConfig struct {
	dir struct {
		system   string
		user     string
		instance string
	}
	user *u.User
}

var (
	procName string
	config   darwinConfig
)

func init() {

	procName = path.Base(os.Args[0])

	var err error

	config.user, err = u.Current()
	if err != nil {
		log.Fatal(err)
	}

	config.dir.system = "/Library/Application Support/" + procName
	config.dir.user = config.user.HomeDir + "/." + procName
	config.dir.instance = "."
}

func platformConfig() Config {
	return &config
}

func (dc *darwinConfig) DirSystem() string {
	return dc.dir.system
}

func (dc *darwinConfig) DirUser() string {
	return dc.dir.user
}

func (dc *darwinConfig) DirInstance() string {
	return dc.dir.instance
}
