package twirp

import (
	"context"
	"github.com/pepeunlimited/accounts/pkg/accountsrpc"
	"github.com/pepeunlimited/apple-iap/pkg/appleiap"
	"github.com/pepeunlimited/apple-iap/pkg/applerpc"
	"github.com/pepeunlimited/checkout/pkg/checkoutrpc"
	"github.com/twitchtv/twirp"
	"testing"
)

func TestCheckoutServer_CreateAppleIAP(t *testing.T) {
	ctx := context.TODO()
	account := accountsrpc.NewAccountsMock(nil, nil)
	appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{0}))

	server := NewCheckoutServer(account, appleiap)
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
	if account.Account.Balance != 140 {
		t.FailNow()
	}
}

func TestCheckoutServer_CreateAppleIAPError(t *testing.T) {
	ctx := context.TODO()
	account := accountsrpc.NewAccountsMock(nil, nil)
	appleiap := applerpc.NewAppleIAPMock(appleiap.NewAppStoreMock([]int{200}))

	server := NewCheckoutServer(account, appleiap)
	_, err := server.CreateAppleIAP(ctx, &checkoutrpc.CreateAppleIAPParams{
		IapReceipt: "1",
		UserId:     1,
		ProductId:  1,
	})
	if err == nil {
		t.FailNow()
	}
	if err.(twirp.Error).Msg() != "apple_iap_internal" {
		t.FailNow()
	}
}