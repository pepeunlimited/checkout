package validator

import (
	"github.com/pepeunlimited/checkout/checkoutrpc"
	"github.com/pepeunlimited/microservice-kit/validator"
	"github.com/twitchtv/twirp"
)

type CheckoutServerValidator struct {}

func (v CheckoutServerValidator) CreateGiftVoucherOrder(params *checkoutrpc.CreateGiftVoucherOrderParams) error {
	if params.UserId == 0 {
		return twirp.RequiredArgumentError("user_id")
	}
	if params.ProductId == 0 {
		return twirp.RequiredArgumentError("product_id")
	}
	if validator.IsEmpty(params.GiftVoucherId) {
		return twirp.RequiredArgumentError("gift_voucher_id")
	}
	return nil
}

func (v CheckoutServerValidator) CreateAppleIAP(params *checkoutrpc.CreateAppleIAPParams) error {
	if params.UserId == 0 {
		return twirp.RequiredArgumentError("user_id")
	}
	if validator.IsEmpty(params.IapReceipt) {
		return twirp.RequiredArgumentError("iap_receipt")
	}
	if params.ProductId == 0 {
		return twirp.RequiredArgumentError("product_id")
	}
	return nil
}

func NewCheckoutServerValidator() CheckoutServerValidator {
	return CheckoutServerValidator{}
}

