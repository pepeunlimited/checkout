package server

import (
	"context"
	"github.com/pepeunlimited/accounts/accountsrpc"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/checkout/internal/app/app1/validator"
)

type CheckoutServer struct {
	validator validator.CheckoutServerValidator
	accounts accountsrpc.AccountService
}

func (server CheckoutServer) CreateAppleIAP(ctx context.Context, params *checkoutrpc.CreateAppleIAPParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateAppleIAP(params)
	if err != nil {
		return nil, err
	}
	// validate the IAP from service

	// if OK
	// else return error

	// => GetProductByID(productId)

	// => COIN (amount for the deposit..)
	_, err = server.accounts.CreateDeposit(ctx, &accountsrpc.CreateDepositParams{
		UserId:      params.UserId,
		Amount:      500, // TODO: check the product_id
		AccountType: "COIN",
	})
	if err != nil {
		return nil, err
	}

	// => FEATURE

	// => mark to the purchases
	return &checkoutrpc.Checkout{}, nil
}

func (server CheckoutServer) CreateOrder(ctx context.Context, params *checkoutrpc.CreateOrderParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateOrder(params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func NewCheckoutServer(accounts accountsrpc.AccountService) CheckoutServer {
	return CheckoutServer {
		validator: 	validator.NewCheckoutServerValidator(),
		accounts: 	accounts,
	}
}