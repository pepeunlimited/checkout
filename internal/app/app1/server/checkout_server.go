package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/pepeunlimited/accounts/accountsrpc"
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/checkout/internal/app/app1/ent"
	"github.com/pepeunlimited/checkout/internal/app/app1/validator"
	"github.com/pepeunlimited/checkout/internal/app/app1/vaultsrepo"
	"github.com/pepeunlimited/microservice-kit/rpcz"
	"github.com/twitchtv/twirp"
	"log"
)

type CheckoutServer struct {
	validator validator.CheckoutServerValidator
	accounts accountsrpc.AccountService
	vaults   vaultsrepo.VaultsRepository
}

func (server CheckoutServer) CreateAppleIAP(ctx context.Context, params *checkoutrpc.CreateAppleIAPParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateAppleIAP(params)
	if err != nil {
		return nil, err
	}
	// validate the IAP from AppleStore

	// if OK
	// else return error

	// => GetProductByID(productId)
	//params.ProductId

	// == COIN
	// => COIN (amount for the deposit..)

	_, err = server.accounts.CreateDeposit(ctx, &accountsrpc.CreateDepositParams{
		UserId:      params.UserId,
		Amount:      500, // TODO: check the product_id
		AccountType: "COIN",
	})
	if err != nil {
		return nil, err
	}

	// ELSE ==
	// => FEATURE

	// => mark to the purchases
	// => save tx to the `vault`

	return &checkoutrpc.Checkout{}, nil
}

func (server CheckoutServer) CreateOrder(ctx context.Context, params *checkoutrpc.CreateOrderParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateOrder(params)
	if err != nil {
		return nil, err
	}

	// => GetProductByID(productId)
	// params.ProductId

	// GetProductPriceByProductID(productID)

	// == LIST_OF_PLACES

	// fromAmount :=      => ProductPrice
	// toUserID   :=      => how 'owns' the product
	// toAmount   :=

	cutAmount 	:= int64(60) //=> cut 60% fromAmount
	toAmount 	:= int64(40) //=> cut 60% fromAmount

	// => save amount&reference number to the `vault`
	referenceNumber := uuid.New().String()
	_, tx, err := server.vaults.Add(ctx, cutAmount, referenceNumber)
	if err != nil {
		return nil, server.isVaultError(err)
	}

	// => LIST_OF_PLACES (amount for the deposit..)
	transfer, err := server.accounts.CreateTransfer(ctx, &accountsrpc.CreateTransferParams{
		FromUserId: params.UserId,
		FromAmount: -100,
		ToUserId:   1,
		ToAmount:   toAmount,
		ReferenceNumber: &wrappers.StringValue{
			Value: referenceNumber,
		},
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Print("checkout: cant rollback: "+err.Error())
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("checkout: cant commit: %v.. transfer=[%v]", err.Error(), transfer)
	}
	// ELSE ==
	// => APP_FEATURE

	// => mark to the purchases
	return &checkoutrpc.Checkout{}, nil
}

func (server CheckoutServer) isVaultError(err error) error {
	switch err {
	case vaultsrepo.ErrReferenceNumberExist:
		return twirp.NewError(twirp.Aborted, "reference number exist").WithMeta(rpcz.Reason, checkoutrpc.ReferenceNumberExist)
	}
	log.Print("vaults-service: unknown error: "+err.Error())
	return twirp.InternalErrorWith(err)
}


func NewCheckoutServer(accounts accountsrpc.AccountService, client *ent.Client) CheckoutServer {
	return CheckoutServer {
		validator: 	validator.NewCheckoutServerValidator(),
		accounts: 	accounts,
		vaults: vaultsrepo.NewVaultsRepository(client),
	}
}