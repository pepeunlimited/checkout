package main

import (
	"github.com/pepeunlimited/accounts/pkg/accountsrpc"
	"github.com/pepeunlimited/apple-iap/pkg/applerpc"
	"github.com/pepeunlimited/checkout/internal/server/twirp"
	"github.com/pepeunlimited/checkout/pkg/checkoutrpc"
	"github.com/pepeunlimited/microservice-kit/middleware"
	"github.com/pepeunlimited/microservice-kit/misc"
	"log"
	"net/http"
)

const (
	Version = "0.1.4"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Printf("Starting the AppleIAPServer... version=[%v]", Version)
	accounts := accountsrpc.NewAccountServiceProtobufClient(misc.GetEnv(accountsrpc.RpcAccountsHost, "api.dev.pepeunlimited.com"), http.DefaultClient)
	appleiap := applerpc.NewAppleIAPServiceProtobufClient(misc.GetEnv(applerpc.RpcAppleIapHost, "api.dev.pepeunlimited.com"), http.DefaultClient)

	cs := checkoutrpc.NewCheckoutServiceServer(twirp.NewCheckoutServer(accounts, appleiap), nil)

	mux := http.NewServeMux()
	mux.Handle(cs.PathPrefix(), middleware.Adapt(cs))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}
}