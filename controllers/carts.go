package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetCart(c *gin.Context) {
	transformedCartItems, err := helpers.TransformCartItems(c, 0)
	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedCartItems))
}

func AddCartItem(c *gin.Context) {
	var body struct {
		ProductID   string `json:"product_id" validate:"required"`
		Quantity    int    `json:"quantity" validate:"required"`
		VariationID string `json:"variation_id"`
	}

	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currVariation models.ProductVariation
	var variationID string

	if body.VariationID != "" {
		if err := helpers.GetCurrentByID(c, &currVariation, body.VariationID); err != nil {
			return
		}
		variationID = currVariation.ID
	} else {
		variationID = ""
	}

	newCartItem := models.CartItem{
		ID:          helpers.NewUUID(),
		ProductID:   body.ProductID,
		VariationID: variationID,
		Quantity:    body.Quantity,
	}

	userID := helpers.GetCurrUserToken(c).ID
	if err := CreateNewCartItem(c, userID, &newCartItem); err != nil {
		return
	}

	helpers.JSONResponse(c, "cart item added successfully")
}

func AddCartItemQuantity(c *gin.Context) {
	var body struct {
		CartID   string `json:"cart_id" validate:"required"`
		Quantity int    `json:"quantity" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currCartItem models.CartItem
	if err := helpers.GetCurrentByID(c, &currCartItem, body.CartID); err != nil {
		return
	}

	currCartItem.Quantity = body.Quantity
	if err := database.DB.Save(&currCartItem).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "Failed to update cart item")
		return
	}

	helpers.JSONResponse(c, "")
}

func DeleteCartItem(c *gin.Context) {
	var body struct {
		CartID string `json:"cart_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currCartItem models.CartItem
	if err := helpers.GetCurrentByID(c, &currCartItem, body.CartID); err != nil {
		return
	}

	helpers.DeleteByModel(c, &currCartItem)
	helpers.JSONResponse(c, "")
}

func CreateNewCartItem(c *gin.Context, userID string, newCartItem *models.CartItem) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&newCartItem).Error; err != nil {
		tx.Rollback()
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "failed to create cart item")
		return err
	}

	userCartEntry := models.UserCart{
		UserID:     userID,
		CartItemID: newCartItem.ID,
	}
	if err := tx.Create(&userCartEntry).Error; err != nil {
		tx.Rollback()
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "failed to link cart item to user")
		return err
	}

	return tx.Commit().Error
}
