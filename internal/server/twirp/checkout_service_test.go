package twirp

import (
	"context"
	"github.com/pepeunlimited/apple-iap/pkg/appstore"
	"github.com/pepeunlimited/apple-iap/pkg/rpc/appleiap"
	"github.com/pepeunlimited/billing/pkg/rpc/order"
	"github.com/pepeunlimited/billing/pkg/rpc/payment"
	"github.com/pepeunlimited/accounts/pkg/rpc/account"
	"github.com/pepeunlimited/checkout/pkg/rpc/checkout"
	"github.com/pepeunlimited/products/pkg/rpc/price"
	"github.com/pepeunlimited/products/pkg/rpc/product"
	"net/http"
	"testing"
)

func TestCheckoutServer_CreateCheckout(t *testing.T) {
	ctx := context.TODO()
	account := account.NewAccountsMock(nil, nil)

	//appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{200}))

	appleiap := appleiap.NewAppleIAPMock(appstore.NewAppStoreMock([]int{0}))
	payments := payment.NewPaymentServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)
	orders 	 := order.NewOrderServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)
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
