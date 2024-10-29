package main

import (
	"log"
	"net/http"

	"github.com/luisthieme/GoMotion/api"
	"github.com/luisthieme/GoMotion/pkg/middleware"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var middlewares = []Middleware{
	middleware.TokenAuthMiddleware,
}

func main() {
	var handler http.HandlerFunc = api.HandleClientProfile

	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	http.HandleFunc("/user/profile", handler)

	log.Println("Server is running port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
