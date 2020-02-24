package main

import (
	"github.com/pepeunlimited/accounts/pkg/rpc/account"
	"github.com/pepeunlimited/apple-iap/pkg/rpc/appleiap"
	"github.com/pepeunlimited/billing/pkg/rpc/order"
	"github.com/pepeunlimited/billing/pkg/rpc/payment"
	"github.com/pepeunlimited/checkout/internal/server/twirp"
	"github.com/pepeunlimited/checkout/pkg/rpc/checkout"
	"github.com/pepeunlimited/microservice-kit/middleware"
	"github.com/pepeunlimited/microservice-kit/misc"
	"github.com/pepeunlimited/products/pkg/rpc/price"
	"github.com/pepeunlimited/products/pkg/rpc/product"
	"log"
	"net/http"
)

const (
	Version = "0.1.4.1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Printf("Starting the CheckoutServer... version=[%v]", Version)

	accountsAddress := misc.GetEnv(account.RpcAccountsHost, "api.dev.pepeunlimited.com")
	appleIAPAddress := misc.GetEnv(appleiap.RpcAppleIapHost, "api.dev.pepeunlimited.com")
	paymentAddress  := misc.GetEnv(payment.RpcPaymentHost, "api.dev.pepeunlimited.com")
	orderAddress    := misc.GetEnv(order.RpcOrderHost, "api.dev.pepeunlimited.com")
	productsAddress := misc.GetEnv(product.RpcProductHost, "api.dev.pepeunlimited.com")
	pricesAddress   := misc.GetEnv(price.RpcPriceHost, "api.dev.pepeunlimited.com")

	accounts := account.NewAccountServiceProtobufClient(accountsAddress, http.DefaultClient)
	appleiap := appleiap.NewAppleIAPServiceProtobufClient(appleIAPAddress, http.DefaultClient)
	payments := payment.NewPaymentServiceProtobufClient(paymentAddress, http.DefaultClient)
	orders   := order.NewOrderServiceProtobufClient(orderAddress, http.DefaultClient)
	products := product.NewProductServiceProtobufClient(productsAddress, http.DefaultClient)
	prices   := price.NewPriceServiceProtobufClient(pricesAddress, http.DefaultClient)

	cs := checkout.NewCheckoutServiceServer(twirp.NewCheckoutServer(accounts, appleiap, orders, payments, products, prices), nil)

	mux := http.NewServeMux()
	mux.Handle(cs.PathPrefix(), middleware.Adapt(cs))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}
}