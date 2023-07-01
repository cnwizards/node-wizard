package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pehlicd/node-wizard/pkg"
	"github.com/pehlicd/node-wizard/pkg/logger"
)

func main() {
	logger.SetupLogger()
	err := pkg.DrainNode()
	if err != nil {
		log.Panicf("Error draining node: %v", err)
	}
	for {
		time.Sleep(100 * time.Second)
	}
}
