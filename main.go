package main

import (
	"time"

	"github.com/pehlicd/node-wizard/pkg"
)

func main() {
	err := pkg.DrainNode()
	if err != nil {
		panic(err.Error())
	}
	for {
		time.Sleep(100 * time.Second)
	}
}
