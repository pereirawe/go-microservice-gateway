package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pereirawe/go-microservice-gateway/config"
	"github.com/pereirawe/go-microservice-gateway/handlers"
	"github.com/pereirawe/go-microservice-gateway/microservices"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Gateway is alive!")
	})

	router.HandleFunc("/login", handlers.LoginHandler)

	apiRouter := router.NewRoute().Subrouter()
	apiRouter.Use(handlers.JWTMiddleware)

	microservices := microservices.GetMicroservices()

	for microservice, url := range microservices {
		path := "/" + microservice + "/{rest:.*}"
		apiRouter.HandleFunc(path, handlers.CreateProxyHandler(url))
	}

	port := cfg.APPPort
	fmt.Printf("\nServer listening on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
