package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/didenko/pald/palc"
)

const uint64max = uint(^uint16(0))

var (
	paldPort, port uint16
	action         string
	err            error
)

func failOnErr(err error) {
	if err != nil {
		println(err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func init() {
	var portParam uint

	flag.UintVar(&portParam, "p", 49200, "PALD port to send the query to.")
	flag.StringVar(&action, "a", "get", "A query action to make - one out of \"get\", \"set\", or \"del\",.")

	flag.Parse()

	if portParam > uint64max {
		fmt.Fprintf(os.Stderr, "%d is too high for a port number (max %d)\n", portParam, uint64max)
		flag.PrintDefaults()
		os.Exit(1)
	}
	paldPort = uint16(portParam)
}

func main() {
	p := palc.New("localhost", paldPort)

	switch action {

	case "get":
		{
			port, err = p.Get(flag.Arg(0))
			failOnErr(err)
			fmt.Println(port)
		}

	case "set":
		{
			port, err = p.Set(flag.Arg(0))
			failOnErr(err)
			fmt.Println(port)
		}

	case "del":
		{
			err = p.Del(flag.Arg(0))
			failOnErr(err)
			fmt.Println("OK")
		}

	default:
		failOnErr(fmt.Errorf("Unrecognised action %q", action))
	}
}
