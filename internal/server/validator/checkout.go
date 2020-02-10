package validator

import (
	"github.com/pepeunlimited/checkout/pkg/checkoutrpc"
	"github.com/twitchtv/twirp"
)

type CheckoutServerValidator struct {}

//func (v CheckoutServerValidator) CreateGiftVoucherOrder(params *checkoutrpc.CreateGiftVoucherOrderParams) error {
//	if params.UserId == 0 {
//		return twirp.RequiredArgumentError("user_id")
//	}
//	if params.ProductId == 0 {
//		return twirp.RequiredArgumentError("product_id")
//	}
//	if validator.IsEmpty(params.GiftVoucherId) {
//		return twirp.RequiredArgumentError("gift_voucher_id")
//	}
//	return nil
//}

func (v CheckoutServerValidator) CreateCheckout(params *checkoutrpc.CreateCheckoutParams) error {
	if params.UserId == 0 {
		return twirp.RequiredArgumentError("user_id")
	}
	if params.PaymentInstrumentId == 0 {
		return twirp.RequiredArgumentError("payment_instrument_id")
	}
	if params.ProductId == 0 {
		return twirp.RequiredArgumentError("product_id")
	}
	return nil
}

func NewCheckoutServerValidator() CheckoutServerValidator {
	return CheckoutServerValidator{}
}