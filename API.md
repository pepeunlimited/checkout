# cURL

### Install
```$ brew install jq > curl ... | jq```

##### CreateAppleIAP
```
$ curl -H "Content-Type: application/json" \
 -X POST "localhost:8080/twirp/pepeunlimited.checkout.CheckoutService/CreateAppleIAP" \
 -d '{"user_id": 1, "product_id": 1, "iap_receipt": "$IAP_RECEIPT"}'
```