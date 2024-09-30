package controllers

import (
	"errors"
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CheckoutOrder(ctx *gin.Context) {
	var body struct {
		Ids       []string `json:"ids" validate:"required"`
		AddressID string   `json:"address_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, "Internal server error")
		}
	}()

	newOrderItem := models.OrderItem{
		ID:        helpers.NewUUID(),
		AddressID: body.AddressID,
	}

	for _, id := range body.Ids {
		if err := processCartItem(ctx, tx, id, &newOrderItem); err != nil {
			tx.Rollback()
			return
		}
	}

	if err := tx.Create(&newOrderItem).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, "Failed to create new order")
		return
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	helpers.JSONResponse(ctx, "")
}

func processCartItem(ctx *gin.Context, tx *gorm.DB, id string, newOrderItem *models.OrderItem) error {
	var cartItem models.CartItem
	if err := helpers.GetCurrentByID(ctx, &cartItem, id); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusNotFound, "Cart item not found")
		return err
	}

	cartItem.OrderItemID = newOrderItem.ID
	cartItem.Status = 1
	if err := tx.Save(&cartItem).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, "Failed to update cart item")
		return err
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(ctx, &currVariation, cartItem.VariationID); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusNotFound, "Product variation not found")
		return err
	}

	if currVariation.Quantity < cartItem.Quantity {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "Insufficient product quantity")
		return errors.New("insufficient quantity")
	}

	currVariation.Quantity -= cartItem.Quantity
	if err := tx.Save(&currVariation).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, "Failed to update product variation")
		return err
	}

	newOrderItem.Price += currVariation.Price * cartItem.Quantity
	newOrderItem.Items = append(newOrderItem.Items, cartItem)

	return nil
}
