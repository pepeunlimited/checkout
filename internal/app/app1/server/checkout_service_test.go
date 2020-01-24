package server

import (
	"context"
	"fmt"
	"github.com/pepeunlimited/accounts/accountsrpc"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/checkout/internal/app/app1/mysql"
	"testing"
)

func TestCheckoutServer_CreateOrder(t *testing.T) {
	ctx := context.TODO()
	mock := accountsrpc.NewAccountsMock(nil, nil, nil)
	server := NewCheckoutServer(mock, mysql.NewEntClient())
	server.vaults.Delete(ctx)
	_, err := server.CreateOrder(ctx, &checkoutrpc.CreateOrderParams{
		UserId:     1,
		ProductId:  1,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	vault, err := server.vaults.GetByReferenceNumber(ctx, mock.ReferenceNumber.Value)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if vault.Amount == 0 {
		t.FailNow()
	}
}

func TestCheckoutServer_CreateOrderRollback(t *testing.T) {
	ctx := context.TODO()
	mock := accountsrpc.NewAccountsMock([]error{fmt.Errorf("custom")}, nil, nil)
	server := NewCheckoutServer(mock, mysql.NewEntClient())
	server.vaults.Delete(ctx)
	_, err := server.CreateOrder(ctx, &checkoutrpc.CreateOrderParams{
		UserId:     1,
		ProductId:  1,
	})
	if err == nil {
		t.FailNow()
	}
	_, err = server.vaults.GetByReferenceNumber(ctx, mock.ReferenceNumber.Value)
	if err == nil {
		t.FailNow()
	}
}

func TestCheckoutServer_CreateAppleIAP(t *testing.T) {
	ctx := context.TODO()
	mock := accountsrpc.NewAccountsMock(nil, nil, nil)
	server := NewCheckoutServer(mock, mysql.NewEntClient())
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

	if mock.Coin.Balance != 600 {
		t.FailNow()
	}

}