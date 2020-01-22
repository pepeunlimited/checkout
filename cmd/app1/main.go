package main

import (
	"github.com/pepeunlimited/checkout/internal/app/app1/server"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/microservice-kit/headers"
	"github.com/pepeunlimited/microservice-kit/middleware"
	"log"
	"net/http"
)

const (
	Version = "0.1"
)

func main() {
	log.Printf("Starting the CheckoutServer... version=[%v]", Version)

	ts := checkoutrpc.NewTodoServiceServer(server.NewCheckoutServer(), nil)

	mux := http.NewServeMux()
	mux.Handle(ts.PathPrefix(), middleware.Adapt(ts, headers.Username()))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}
}