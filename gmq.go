package main

import (
	"flag"
)

var configfile string

func main() {

	flag.StringVar(&configfile, "f", "/etc/gmq/gmq.json", "Configuration file")
	flag.Parse()

	//	if err := parseConfigFile(configfile); err != nil {
	//		panic(errors.New("Unable to parse configuration file %s"))
	//	}
}
