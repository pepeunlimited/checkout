package twirp

import (
	"context"
	"github.com/pepeunlimited/accounts/pkg/accountsrpc"
	"github.com/pepeunlimited/apple-iap/pkg/appleiap"
	"github.com/pepeunlimited/apple-iap/pkg/applerpc"
	"github.com/pepeunlimited/billing/pkg/orderrpc"
	"github.com/pepeunlimited/billing/pkg/paymentrpc"
	"github.com/pepeunlimited/checkout/pkg/checkoutrpc"
	"net/http"
	"testing"
)

func TestCheckoutServer_CreateCheckout(t *testing.T) {
	ctx := context.TODO()
	account := accountsrpc.NewAccountsMock(nil, nil)

	//appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{200}))
	appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{0}))


	payments := paymentrpc.NewPaymentServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)
	orders := orderrpc.NewOrderServiceProtobufClient("api.dev.pepeunlimited.com", http.DefaultClient)


	paymentInstrumentId 	:= uint32(1)
	userId			   		:=  int64(2)
	productId          		:=  int64(3)

	server := NewCheckoutServer(account, appleiap, orders, payments)
	iap, err := server.CreateCheckout(ctx, &checkoutrpc.CreateCheckoutParams{
		PaymentInstrumentId: paymentInstrumentId,
		UserId:              userId,
		ProductId:           productId,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if iap == nil {
		t.FailNow()
	}
	if account.Account.Balance != 140 {
		t.FailNow()
	}
}
