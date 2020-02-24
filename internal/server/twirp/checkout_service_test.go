package twirp

import (
	"context"
	"github.com/pepeunlimited/accounts/pkg/accountsrpc"
	"github.com/pepeunlimited/apple-iap/pkg/appleiap"
	"github.com/pepeunlimited/apple-iap/pkg/applerpc"
	"github.com/pepeunlimited/billing/pkg/orderrpc"
	"github.com/pepeunlimited/billing/pkg/paymentrpc"
	"github.com/pepeunlimited/checkout/pkg/rpc/checkout"
	"github.com/pepeunlimited/products/pkg/rpc/price"
	"github.com/pepeunlimited/products/pkg/rpc/product"
	"net/http"
	"testing"
)

func TestCheckoutServer_CreateCheckout(t *testing.T) {
	ctx := context.TODO()
	account := accountsrpc.NewAccountsMock(nil, nil)

	//appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{200}))
	appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{0}))
	payments := paymentrpc.NewPaymentServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)
	orders 	 := orderrpc.NewOrderServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)
	products := product.NewProductServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)
	prices   := price.NewPriceServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)

	paymentInstrumentId 	:= uint32(11)
	userId			   		:=  int64(2)
	productId          		:=  int64(77)

	server := NewCheckoutServer(account, appleiap, orders, payments, products, prices)
	_, err := server.CreateCheckout(ctx, &checkout.CreateCheckoutParams{
		PaymentInstrumentId: paymentInstrumentId,
		UserId:              userId,
		ProductId:           productId,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// initial account balance is 100 after deposit account should be 200
	if account.Account.Balance != 200 {
		t.FailNow()
	}
}
