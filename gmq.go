package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var configfile string

func main() {

	flag.StringVar(&configfile, "f", "gmq.conf", "Location of configuation file")
	flag.Parse()

	if err := parseConfigFile(configfile); err != nil {
		return errors.New(fmt.Printf("Unable to parse configuration file %s", configfile))
	}

}

func parseConfigFile() error {

}
