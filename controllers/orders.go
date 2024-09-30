package controllers

import (
	"github.com/forthedreamers-server/helpers"
	"github.com/gin-gonic/gin"
)

func CheckoutOrder(ctx *gin.Context) {
	var body struct {
		Ids       []string `json:"ids" validate:"required"`
		AddressID string   `json:"address_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	// newOrder := models.OrderItem{
	// 	ID: helpers.NewUUID(),
	// }

	helpers.JSONResponse(ctx, "")
}
