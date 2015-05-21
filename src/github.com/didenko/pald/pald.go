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

package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/didenko/pald/internal/server"
	"github.com/spf13/viper"
	"github.com/takama/daemon"
)

const (
	daemonName = "pald"
	daemonDesc = "Port Allocator Daemon"
)

var (
	stdlog, errlog            *log.Logger
	portMin, portMax, portSvr uint16
)

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {
	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status", daemonName)

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	server.Run(portSvr, portMin, portMax)
	// never happen, but need to complete code
	return usage, nil
}

func downcast(i int, name string) uint16 {
	if i < 0 || i > 65535 {
		panic(fmt.Sprintf("Variable %s = %d is not convertible to uint16", name, i))
	}
	return uint16(i)
}

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configPath := usr.HomeDir + "/.pald"
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	viper.SetDefault("port_min", 49201)
	viper.SetDefault("port_max", 49999)
	viper.SetDefault("port_listen", 49200)
	viper.SetDefault("dump_file", configPath+"/dump")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	portMin = downcast(viper.GetInt("port_min"), "port_min")
	portMax = downcast(viper.GetInt("port_max"), "port_max")
	portSvr = downcast(viper.GetInt("port_listen"), "port_listen")

	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	srv, err := daemon.New(daemonName, daemonDesc)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
