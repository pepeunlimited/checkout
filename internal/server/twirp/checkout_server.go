package twirp

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/pepeunlimited/accounts/pkg/accountsrpc"
	"github.com/pepeunlimited/apple-iap/pkg/applerpc"
	"github.com/pepeunlimited/billing/pkg/orderrpc"
	"github.com/pepeunlimited/billing/pkg/paymentrpc"
	"github.com/pepeunlimited/checkout/internal/server/validator"
	"github.com/pepeunlimited/checkout/pkg/checkoutrpc"
	"github.com/pepeunlimited/products/pkg/rpc/price"
	"github.com/pepeunlimited/products/pkg/rpc/product"
	"github.com/twitchtv/twirp"
	"log"
)

type CheckoutServer struct {
	validator validator.CheckoutServerValidator
	accounts  accountsrpc.AccountService
	products  product.ProductService
	prices    price.PriceService
	orders 	  orderrpc.OrderService
	payments  paymentrpc.PaymentService
	iap       applerpc.AppleIAPService
}

func (server CheckoutServer) CreateCheckout(ctx context.Context, params *checkoutrpc.CreateCheckoutParams) (*checkoutrpc.Checkout, error) {
	err := server.validator.CreateCheckout(params)
	if err != nil {
		return nil, err
	}
	order, err := server.orders.CreateOrder(ctx, &orderrpc.CreateOrderParams{
		OrderItems: []*orderrpc.OrderItem{&orderrpc.OrderItem{PriceId:  params.UserId, Quantity: 1}},
		UserId: params.UserId,
	})
	if err != nil {
		return nil, err
	}
	product, err := server.products.GetProduct(ctx, &product.GetProductParams{ProductId: params.ProductId})
	if err != nil {
		return nil, err
	}
	price, err := server.prices.GetPrice(ctx, &price.GetPriceParams{
		ProductId: product.Id,
	})
	if err != nil {
		return nil, err
	}
	if price.Price == 0 || price.Discount == 0 {
		// mark as paid
		payment, err := server.payments.CreatePayment(ctx, &paymentrpc.CreatePaymentParams{
			OrderId:             order.Order.Id,
			PaymentInstrumentId: params.PaymentInstrumentId,
			UserId:              params.UserId,
		})
		if err != nil {
			return nil, err
		}
		return &checkoutrpc.Checkout{
			OrderId:             order.Order.Id,
			PaymentId:           payment.Id,
			PaymentInstrumentId: payment.PaymentInstrumentId,
		}, nil
	}
	paymentInstrument, err := server.payments.GetPaymentInstrument(ctx, &paymentrpc.GetPaymentInstrumentParams{Id: params.PaymentInstrumentId})
	if err != nil {
		return nil, err
	}
	switch paymentInstrument.Type {
	case "APPLE":
		err := server.apple(ctx, "receipt", params.UserId, product, price)
		if err != nil {
			return nil, err
		}
	case "GIFT_VOUCHER":
		//server.billing.validateGiftVoucher
		//server.billing.useGiftVoucher
	default:
		log.Print("not supported payment instrument: "+paymentInstrument.Type)
		return nil, twirp.NewError(twirp.Aborted, "not_supported_payment_instrument")
	}
	// => mark to the order & payments
	payment, err := server.payments.CreatePayment(ctx, &paymentrpc.CreatePaymentParams{
		OrderId:             order.Order.Id,
		PaymentInstrumentId: params.PaymentInstrumentId,
		UserId:              params.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &checkoutrpc.Checkout{
		OrderId:             order.Order.Id,
		PaymentId:           payment.Id,
		PaymentInstrumentId: payment.PaymentInstrumentId,
	}, nil
}

func (server CheckoutServer) apple(ctx context.Context, receipt string, userId int64, product *product.Product, price *price.Price) error {
	// execute validation for the IAP from AppleStore
	verified, err := server.iap.VerifyReceipt(ctx, &applerpc.VerifyReceiptParams{
		Receipt: receipt,
	})
	if err != nil {
		log.Printf("iap validation failed: "+err.Error())
		return err
	}
	log.Print(verified)
	amount := int64(0)
	if price.Price >= price.Discount {
		amount = price.Price
	} else {
		amount = price.Discount
	}
	//product, err := server.products.GetProductByID(productId)
	// == LIST
	//switch product.Types {
	//case "LIST":
	//	if verified.Type != "CONSUMABLE" {
	//		log.Printf("FATAL: wrong product type set to the apple iap! Should be CONSUMABLE receipt=%v, productId=%v",receipt, productId)
		//}
		//price, err := server.prices.GetProductPriceByID(productId)
		//if err != nil {
		//	return err
		//}
		// fromAmount :=      => ProductPrice
		// toUserID   :=      => how 'owns' the product
		// toAmount   :=
		// cutAmount  := 	  => int64(60) //=> cut 60% fromAmount
		toAmount 	:= amount //=> cut 60% fromAmount
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
	//}
	return nil
}

func NewCheckoutServer(accounts accountsrpc.AccountService,
	iap applerpc.AppleIAPService,
	orders orderrpc.OrderService,
	payments paymentrpc.PaymentService,
	products product.ProductService,
	prices   price.PriceService) CheckoutServer {
	return CheckoutServer {
		validator: 	validator.NewCheckoutServerValidator(),
		iap:       	iap,
		accounts:  	accounts,
		payments:	payments,
		orders:		orders,
		products:	products,
		prices:  	prices,
	}
}