package main

import (
	"github.com/pehlicd/node-wizard/pkg"
)

func main() {
	err := pkg.DrainNode()
	if err != nil {
		panic(err.Error())
	}
}
