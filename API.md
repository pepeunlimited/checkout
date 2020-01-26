# cURL

### Install
```$ brew install jq > curl ... | jq```

##### CreateAppleIAP
```
$ curl -H "Content-Type: application/json" \
 -X POST "api.dev.pepeunlimited.com/twirp/pepeunlimited.checkout.CheckoutService/CreateAppleIAP" \
 -d '{"user_id": 1, "product_id": 1, "iap_receipt": "$IAP_RECEIPT"}'
```
##### CreateGiftVoucherOrder
```
$ curl -H "Content-Type: application/json" \
 -X POST "api.dev.pepeunlimited.com/twirp/pepeunlimited.checkout.CheckoutService/CreateGiftVoucherOrder" \
 -d '{"user_id": 1, "product_id": 1, "gift_voucher_id": "gift_voucher_id"}'
```