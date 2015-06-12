package main

import (
	"flag"
	"fmt"
)

var configfile string

func main() {

	flag.StringVar(&configfile, "f", "gmq.conf", "Location of configuation file")
	flag.Parse()

	fmt.Printf("GMQ Server running with configurations from %s", configfile)
}
