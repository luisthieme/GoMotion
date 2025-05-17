package main

import (
	"github.com/luisthieme/GoMotion/core"
)



func main() {
	engine := core.NewEngine("luis_engine", "6969")
	engine.Start()

	// engine2 := core.NewEngine("backup_engine", "8787")

	// engine2.Start()
}
