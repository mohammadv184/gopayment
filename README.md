<p align="center"><img src="golang-logo.png?raw=true"></p>


# GoPayment

[![Go Reference](https://pkg.go.dev/badge/github.com/mohammadv184/gopayment.svg)](https://pkg.go.dev/github.com/mohammadv184/gopayment)
[![GitHub license](https://img.shields.io/github/license/mohammadv184/gopayment)](https://github.com/mohammadv184/gopayment/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/mohammadv184/gopayment)](https://github.com/mohammadv184/gopayment/issues)
[![codecov](https://codecov.io/gh/mohammadv184/gopayment/branch/main/graph/badge.svg?token=VO7KKJTIU2)](https://codecov.io/gh/mohammadv184/gopayment)
[![Go Report Card](https://goreportcard.com/badge/github.com/mohammadv184/gopayment)](https://goreportcard.com/report/github.com/mohammadv184/gopayment)
[![Build Status](https://app.travis-ci.com/mohammadv184/gopayment.svg?branch=main)](https://app.travis-ci.com/mohammadv184/gopayment)

Multi Gateway Payment Package for Golang.

# List of contents

- [GoPayment](#gopayment)
- [List of contents](#list-of-contents)
- [List of available drivers](#list-of-available-drivers)
    - [Installation](#Installation)
    - [How to use](#how-to-use)
        - [Purchase](#purchase)
        - [Pay](#pay)
        - [Verify payment](#verify-payment)
        - [Working with invoices](#working-with-invoices)
        - [Working with receipts](#working-with-receipts)
        - [Example](#example)
- [Contributing](#contributing)
- [Security](#security)
- [Credits](#credits)
- [License](#license)

# List of available drivers
- [zarinpal](https://www.zarinpal.com/) :white_check_mark:
- [payping](https://www.payping.ir/) :white_check_mark:
- [asanpardakht](https://asanpardakht.ir/) :white_check_mark:
- [behpardakht (mellat)](http://www.behpardakht.com/) :x: (will be added soon)
- [digipay](https://www.mydigipay.com/) :x: (will be added soon)
- [etebarino (Installment payment)](https://etebarino.com/) :x: (will be added soon)
- [sepehr (saderat)](https://www.sepehrpay.com/) :x: (will be added soon)
- [idpay](https://idpay.ir/) :white_check_mark:
- [poolam](https://poolam.ir/) :x: (will be added soon)
- [irankish](http://irankish.com/) :x: (will be added soon)
- [nextpay](https://nextpay.ir/) :x: (will be added soon)
- [parsian](https://www.pec.ir/) :x: (will be added soon)
- [pasargad](https://bpi.ir/) :x: (will be added soon)
- [payir](https://pay.ir/) :x: (will be added soon)
- [paystar](http://paystar.ir/) :x: (will be added soon)
- [sadad (melli)](https://sadadpsp.ir/) :x: (will be added soon)
- [saman](https://www.sep.ir) :x: (will be added soon)
- [walleta (Installment payment)](https://walleta.ir/) :x: (will be added soon)
- [yekpay](https://yekpay.com/) :x: (will be added soon)
- [zibal](https://www.zibal.ir/) :white_check_mark:
- [sepordeh](https://sepordeh.com/) :x: (will be added soon)
- [paypal](http://www.paypal.com/) :x: (will be added soon)
- **Help me to add the other gateways by creating `pull requests`**

## Installation

```
go get -u github.com/mohammadv184/gopayment
```

## How to use


#### Purchase
In order to pay the invoice, we need the payment transactionId.
We purchase the invoice to retrieve transaction id:
```go
// Configure the Gateway Driver
gateway:=&payping.Driver{
Callback:    "http://example.test/callback",
Token:       "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
}
// Create new Payment.
payment := gopayment.NewPayment(gateway)
// Set Invoice Amount.
payment.Amount(amountInt)
// Set Invoice Description.
payment.Description("description")
// Purchase the invoice.
err = payment.Purchase()
if err != nil {
    fmt.Println(err)
}
// Get Transaction ID
transactionID := payment.GetTransactionID()
```
#### Pay
After purchasing the invoice, we can redirect the user to the bank payment page:
```go
func pay() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Configure the Gateway Driver
        gateway:=&payping.Driver{
        Callback:    "http://example.test/callback",
        Token:       "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        }
        // Create new Payment.
        payment := gopayment.NewPayment(gateway)
        // Set Invoice Amount.
		payment.Amount(amountInt)
		// Set Invoice Description.
		payment.Description("description")
		// Purchase the invoice.
        err = payment.Purchase()
        if err != nil {
            fmt.Println(err)
        }
        // Get Transaction ID And Save it to the database.
        transactionID := payment.GetTransactionID()
        // Redirect the user to the bank payment page.
        payUrl := payment.PayURL()
		c.Redirect(http.StatusFound, payUrl)
	}
}

```
#### Verify payment

When user has completed the payment, the bank redirects them to your website, then you need to **verify your payment** to make sure that the payment is valid.:
```go
// Configure the Gateway Driver
gateway:=&payping.Driver{
Callback:    "http://example.test/callback",
Token:       "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
}

refId := c.Query("refId")
VerifyRequest:=&payping.VerifyRequest{
Amount: "100",
RefID:  refID,
}

if receipt, err := gateway.Verify(VerifyRequest); err == nil {

c.JSON(200, gin.H{
"status": "success",
"data":   receipt.GetReferenceID(),
"date":   receipt.GetDate().String(),
"card":   receipt.GetDetail("cardNumber"),
"gateway": receipt.GetDriver(),
})
} else {
c.JSON(200, gin.H{
"message": "error " + err.Error(),
})
}
```
#### Working with invoices

When you make a payment, the invoice is automatically generated within the payment


In your code, use it like the below:
```go
// Create new Payment.
payment := gopayment.NewPayment(gateway)
// Get the invoice.
invoice:=payment.GetInvoice()
// Set Invoice Amount.
invoice.SetAmount(1000)
// Set Invoice Description.
invoice.SetDescription("description")
// Set Invoice Deatils.
invoice.Detail("phone","0912345678")
invoice.Detail("email","example@example.com")

```
Available methods:

- `SetUuid`: set the invoice unique id
- `GetUuid`: retrieve the invoice current unique id
- `Detail`: attach some custom details into invoice
- `GetDetail`: retrieve the invoice detail
- `GetDetails`: retrieve all custom details
- `SetAmount`: set the invoice amount
- `GetAmount`: retrieve invoice amount
- `SetTransactionID`: set invoice payment transaction id
- `GetTransactionID`: retrieve payment transaction id
- `SetDescription`: set invoice Description
- `GetDescription`: retrieve payment Description
#### Working with receipts

When you verify a payment, the receipt is automatically generated.


In your code, use it like the below:
```go
// Verify the Payment.
receipt := gateway.verify(...)
// Get the Payment Reference ID.
refId := receipt.GetReferenceID()
// Get the payment date .
paymentDate:=receipt.GetDate()
// Get the payment driver name.
paymentDriver:=receipt.GetDriver()
// Get payment Deatils.
cardNum:=receipt.GetDetail("cardNumber")
cardHNum:=receipt.GetDetail("HashedCardNumber")


```
Available methods:

- `GetReferenceID`: retrieve the payment reference id
- `GetDriver`: retrieve the payment driver name
- `Detail`: attach some custom details into invoice
- `GetDate`: retrieve payment date
- `GetDetail`: retrieve the invoice detail
- `GetDetails`: retrieve all custom details
#### Example
There is an example project using GoPayment you can find at [GoPayment-Example](https://github.com/mohammadv184/gopayment-example) Repo It contains a payment implementation.
## Contributing

Please see [CONTRIBUTING](CONTRIBUTING.md) and [CONDUCT](CONDUCT.md) for details.
## Security

If you discover any security related issues, please email mohammadv184@gmail.com instead of using the issue tracker.

## Credits

- [Mohammad Abbasi](https://github.com/mohammadv184)
- [All Contributors](../../contributors)

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
