package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammadv184/gopayment"
	"github.com/mohammadv184/gopayment/drivers"
	"github.com/mohammadv184/gopayment/drivers/payping"
	"strconv"
)

var Gateway drivers.Driver

func init() {
	Gateway = &payping.Driver{
		Callback:    "http://localhost:8080/callback",
		Description: "test",
		Token:       "7ce86859fdebcdbcc4a4820b6f8d9752d8bacf4a858d4ed2f6c5f7580215ef44",
	}
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/pay/:amount", pay())
	router.POST("/callback", callback())
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
func pay() gin.HandlerFunc {
	return func(c *gin.Context) {
		amount := c.Param("amount")
		amountInt, _ := strconv.Atoi(amount)
		Payment := gopayment.NewPayment(Gateway)
		err := Payment.Amount(amountInt)
		if err != nil {
			return
		}
		Payment.Purchase()

		c.Redirect(302, Payment.PayUrl())
	}
}
func callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		refID := c.PostForm("refid")
		if rec, ok := Gateway.Verify("100", refID); ok == nil {
			cardNum, _ := rec.GetDetail("cardNumber")
			c.JSON(200, gin.H{
				"status": "success",
				"data":   rec.GetReferenceID(),
				"date":   rec.GetDate().String(),
				"card":   cardNum,
				"driver": rec.GetDriver(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "error " + ok.Error(),
			})
		}
	}
}
