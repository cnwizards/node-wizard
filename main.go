package main

import (
	"github.com/pehlicd/node-wizard/pkg"
	"github.com/pehlicd/node-wizard/pkg/logger"
)

func main() {
	logger.SetupLogger()
	pkg.Run()
}
