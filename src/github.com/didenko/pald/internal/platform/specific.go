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

// Config provides the information on where to locate information files
// for this applciation instance
type Config interface {

	// DirSystem returns a system-wide location of config files for this application
	DirSystem() string

	// DirUser returns a user-specific location of config files for this application
	DirUser() string

	// DirUser returns q location of config files for this running instance
	DirInstance() string
}

func GetConfig() Config {
	return platformConfig()
}
