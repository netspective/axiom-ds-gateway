package main

import (
	"flag"
	"fmt"
	"netspective/axiomds/service"
	"netspective/axiomds/config"
)
const VERSION = "0.1.0"

func main() {
	var options service.Options
	if(options.InitCommandLine()) {
		if(options.Help) {
			flag.Usage()
		}

		config := *config.NewConfiguration()
		config.Configure(options.ConfigFile)
	} else {
		fmt.Printf("Invalid command line, aborting.")
		flag.Usage()
	}
}
