package twirp

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/pepeunlimited/accounts/pkg/accountsrpc"
	"github.com/pepeunlimited/apple-iap/pkg/applerpc"
	"github.com/pepeunlimited/checkout/internal/server/validator"
	"github.com/pepeunlimited/checkout/pkg/checkoutrpc"
	"github.com/twitchtv/twirp"
	"log"
)

type CheckoutServer struct {
	validator validator.CheckoutServerValidator
	accounts  accountsrpc.AccountService
	//billing
	//products
	iap       applerpc.AppleIAPService
}

func (server CheckoutServer) CreateSubscription(ctx context.Context, params *checkoutrpc.CreateSubscriptionParams) (*checkoutrpc.Checkout, error) {
	product, err := server.products.GetProductById(params.ProductId)
	if err != nil {
		return nil, err
	}
	// is the product trialable?
	if params.UseTrial {
		// does the user has used it?
		server.billing.startTrial()
		return &checkoutrpc.Checkout{}, nil
	}
	// is the product subscribable?
	// => mark to subscription
	server.billing.startSubscription()
	return &checkoutrpc.Checkout{}, nil
}

func (server CheckoutServer) CreateCheckout(ctx context.Context, params *checkoutrpc.CreateCheckoutParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateAppleIAP(params)
	if err != nil {
		return nil, err
	}

	isFree, err := server.isProductFree(params.ProductId)
	if err != nil {
		log.Print("can't verify is the product set as free: "+err.Error())
		return nil, err
	}
	if *isFree {
		server.billing.createOrder(productId, userId)
		return &checkoutrpc.Checkout{}, nil
	}
	//paymentInstrument, err := server.billing.GetPaymentInstrumentByID(params.PaymentInstrumentId)
	switch paymentInstrument.Name() {
	case "APPLE_IAP":
		err := server.appleiap(ctx, params.PaymentInstrumentType, params.UserId, params.ProductId)
		if err != nil {
			return nil, err
		}
	case "GIFT_VOUCHER":
		//server.billing.validateGiftVoucher
		//server.billing.useGiftVoucher
	default:
		log.Print("not supported payment instrument: "+paymentInstrument.Name())
		return nil, twirp.NewError(twirp.Aborted, "not_supported_payment_instrument")
	}
	// => mark to the purchases
	server.billing.createOrder(productId, userId)
	return &checkoutrpc.Checkout{}, nil
}

func (server CheckoutServer) isProductFree(productId int64) (*bool, error) {
	return nil, nil
}

func (server CheckoutServer) appleiap(ctx context.Context, receipt string, userId int64, productId int64) error {
	// execute validation for the IAP from AppleStore
	verified, err := server.iap.VerifyReceipt(ctx, &applerpc.VerifyReceiptParams{
		Receipt: receipt,
	})
	if err != nil {
		log.Printf("iap validation failed: "+err.Error())
		return err
	}
	product, err := server.products.GetProductByID(productId)
	// == LIST
	switch product.Types {
	case "LIST":
		if verified.Type != applerpc.VerifyReceiptResponse_CONSUMABLE {
			log.Printf("FATAL: wrong product type set to the apple iap! Should be CONSUMABLE receipt=%v, productId=%v",receipt, productId)
		}
		//price, err := server.prices.GetProductPriceByID(productId)
		if err != nil {
			return err
		}
		// fromAmount :=      => ProductPrice
		// toUserID   :=      => how 'owns' the product
		// toAmount   :=
		// cutAmount  := 	  => int64(60) //=> cut 60% fromAmount
		toAmount 	:= int64(40) //=> cut 60% fromAmount
		// => save amount&reference number to the `vault`
		referenceNumber := uuid.New().String()
		// TODO: backoff https://github.com/jpillora/backoff
		_, err = server.accounts.CreateDeposit(ctx, &accountsrpc.CreateDepositParams{
			UserId:          userId,
			Amount:          toAmount,
			ReferenceNumber: &wrappers.StringValue{Value: referenceNumber},
		})
		if err != nil {
			log.Print("deposit failed: "+err.Error())
			return err
		}
	}
	return nil
}

func NewCheckoutServer(accounts accountsrpc.AccountService, iap applerpc.AppleIAPService) CheckoutServer {
	return CheckoutServer {
		validator: validator.NewCheckoutServerValidator(),
		iap:       iap,
		accounts:  accounts,
	}
}