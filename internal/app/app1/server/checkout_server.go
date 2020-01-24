package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/pepeunlimited/accounts/accountsrpc"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/checkout/internal/app/app1/validator"
)

type CheckoutServer struct {
	validator validator.CheckoutServerValidator
	accounts accountsrpc.AccountService
}

func (server CheckoutServer) CreateGiftVoucherOrder(ctx context.Context, params *checkoutrpc.CreateGiftVoucherOrderParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateGiftVoucherOrder(params)
	if err != nil {
		return nil, err
	}



	return &checkoutrpc.Checkout{}, nil
}

func (server CheckoutServer) CreateAppleIAP(ctx context.Context, params *checkoutrpc.CreateAppleIAPParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateAppleIAP(params)
	if err != nil {
		return nil, err
	}
	// execute validation for the IAP from AppleStore

	// if OK
	// else return error

	// => GetProductByID(productId)
	// params.ProductId

	// GetProductPriceByProductID(productID)

	// == LIST_OF_PLACES

	// fromAmount :=      => ProductPrice
	// toUserID   :=      => how 'owns' the product
	// toAmount   :=

	//cutAmount 	:= int64(60) //=> cut 60% fromAmount
	toAmount 	:= int64(40) //=> cut 60% fromAmount
	// => save amount&reference number to the `vault`
	referenceNumber := uuid.New().String()
	_, err = server.accounts.CreateDeposit(ctx, &accountsrpc.CreateDepositParams{
		UserId:          params.UserId,
		Amount:          toAmount,
		ReferenceNumber: &wrappers.StringValue{Value: referenceNumber},
	})
	if err != nil {
		return nil, err
	}

	// ELSE ==
	// => APP_FEATURE

	// => mark to the purchases

	return &checkoutrpc.Checkout{}, nil
}

func NewCheckoutServer(accounts accountsrpc.AccountService) CheckoutServer {
	return CheckoutServer {
		validator: 	validator.NewCheckoutServerValidator(),
		accounts: 	accounts,
	}
}