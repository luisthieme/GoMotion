package main

import (
	"github.com/luisthieme/GoMotion/core"
)



func main() {
	engine := core.NewEngine("go_motion", "localhost:777")
	engine.Start()
}
