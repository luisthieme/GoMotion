package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luisthieme/GoMotion/api"
	"github.com/luisthieme/GoMotion/internal"
)



func main() {
	internal.InitDB()
	api.RegisterRoutes()
	log.Printf("Starting %s on port %d...\n", internal.EngineName, internal.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", internal.Port), nil))
}
