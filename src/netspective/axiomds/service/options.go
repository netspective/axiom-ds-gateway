package service

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	ConfigFile   string
	Help         bool
	Verbose      bool
	OptionsValid bool
}

func (o *Options) logError(message string) {
	o.OptionsValid = false
	fmt.Printf(message)
}

func (o *Options) InitCommandLine() bool {
	o.OptionsValid = true

	flag.StringVar(&o.ConfigFile, "conf", "axiomds.conf.json", "Configuration file path, name, and extension")
	flag.BoolVar(&o.Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&o.Help, "help", false, "Display usage information")

	flag.Parse()

	_, err := os.Stat(o.ConfigFile);
	if(err != nil) {
		o.logError(fmt.Sprintf("Config file '%s' error '%s'", o.ConfigFile, err))
		return o.OptionsValid
	}

	return o.OptionsValid
}


