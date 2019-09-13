package main

import (
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("hi")
	usage := `iOS client v 0.01

Usage:
  qvh devices 
 
Options:
  -h --help     Show this screen.
  --version     Show version.
  -u=<udid>, --udid     UDID of the device.
  -o=<filepath>, --output
  `
	arguments, _ := docopt.ParseDoc(usage)
	devices, _ := arguments.Bool("devices")
	if devices {
		log.Info("devices")
	}

}
