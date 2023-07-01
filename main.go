package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pehlicd/node-wizard/pkg/logger"
	"github.com/pehlicd/node-wizard/pkg/utils"
)

func main() {
	logger.SetupLogger()
	err := utils.DrainNode()
	if err != nil {
		log.Panicf("Error draining node: %v", err)
	}
	for {
		time.Sleep(100 * time.Second)
	}
}
