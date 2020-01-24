package server

import (
	"context"
	"github.com/pepeunlimited/accounts/accountsrpc"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"testing"
)

func TestCheckoutServer_CreateAppleIAP(t *testing.T) {
	ctx := context.TODO()
	mock := accountsrpc.NewAccountsMock(nil, nil)
	server := NewCheckoutServer(mock)
	iap, err := server.CreateAppleIAP(ctx, &checkoutrpc.CreateAppleIAPParams{
		IapReceipt: "1",
		UserId:     1,
		ProductId:  1,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if iap == nil {
		t.FailNow()
	}
	if mock.Account.Balance != 140 {
		t.FailNow()
	}
}