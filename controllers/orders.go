package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
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

	newOrderItem := models.OrderItem{
		ID: helpers.NewUUID(),
	}

	for _, v := range body.Ids {
		var cartItem models.CartItem
		if err := helpers.GetCurrentByID(ctx, &cartItem, v); err != nil {
			return
		}

		cartItem.OrderItemID = newOrderItem.ID
		if err := database.DB.Save(&cartItem).Error; err != nil {
			helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helpers.JSONResponse(ctx, "")
}
