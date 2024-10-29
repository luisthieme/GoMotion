package main

import (
	"log"
	"net/http"

	"github.com/luisthieme/GoMotion/api"
)

func main() {
	api.RegisterRoutes()
	log.Println("Server is running port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
