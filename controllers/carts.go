package controllers

import (
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetCart(ctx *gin.Context) {
	currUser := helpers.GetCurrUserToken(ctx, "CartItems")

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&currUser.CartItems))
}

func AddToCart(ctx *gin.Context) {
	var body struct {
		ProductID   string `json:"product_id" validate:"required"`
		Quantity    int    `json:"quantity" validate:"required"`
		VariationID string `json:"variation_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ProductID); err != nil {
		return
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(ctx, &currVariation, body.VariationID); err != nil {
		return
	}

	user := helpers.GetCurrUserToken(ctx)
	newCartItem := models.CartItem{
		ID:          helpers.NewUUID(),
		ProductID:   currProduct.ID,
		VariationID: body.VariationID,
		Quantity:    body.Quantity,
		UserID:      user.ID,
	}
	if err := helpers.CreateNewData(ctx, &newCartItem); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}
