package main

import (
	"github.com/pepeunlimited/accounts/accountsrpc"
	"github.com/pepeunlimited/checkout/internal/app/app1/mysql"
	"github.com/pepeunlimited/checkout/internal/app/app1/server"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/microservice-kit/middleware"
	"github.com/pepeunlimited/microservice-kit/misc"
	"log"
	"net/http"
)

const (
	Version = "0.1.2"
)

func main() {
	log.Printf("Starting the CheckoutServer... version=[%v]", Version)

	accounts := accountsrpc.NewAccountServiceProtobufClient(misc.GetEnv(accountsrpc.RpcAccountsHost, "api.dev.pepeunlimited.com"), http.DefaultClient)

	client := mysql.NewEntClient()

	cs := checkoutrpc.NewCheckoutServiceServer(server.NewCheckoutServer(accounts, client), nil)

	mux := http.NewServeMux()
	mux.Handle(cs.PathPrefix(), middleware.Adapt(cs))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}
}