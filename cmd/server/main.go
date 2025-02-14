package main

import (
	"github.com/luisthieme/GoMotion/core"
	"github.com/luisthieme/GoMotion/internal"
)



func main() {
	internal.InitDB()
	engine := core.Engine{}
	engine.Start()
}
