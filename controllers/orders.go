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

func CheckoutOrder(c *gin.Context) {
	var body struct {
		Ids           []string `json:"ids" validate:"required"`
		AddressID     string   `json:"address_id" validate:"required"`
		PaymentMethod int      `json:"payment_method" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			helpers.ErrJSONResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	newOrderItem := models.OrderItem{
		ID:            helpers.NewUUID(),
		AddressID:     body.AddressID,
		PaymentMethod: body.PaymentMethod,
		UserID:        helpers.GetCurrUserToken(c).ID,
	}

	if err := tx.Create(&newOrderItem).Error; err != nil {
		tx.Rollback()
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "Failed to create new order")
		return
	}

	for _, id := range body.Ids {
		if err := processCartItem(c, tx, id, &newOrderItem); err != nil {
			tx.Rollback()
			return
		}
	}

	if err := tx.Save(&newOrderItem).Error; err != nil {
		tx.Rollback()
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "Failed to update order price")
		return
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	helpers.JSONResponse(c, "Order Placed Successfully")
}

func processCartItem(c *gin.Context, tx *gorm.DB, id string, newOrderItem *models.OrderItem) error {
	var cartItem models.CartItem
	if err := helpers.GetCurrentByID(c, &cartItem, id); err != nil {
		helpers.ErrJSONResponse(c, http.StatusNotFound, "Cart item not found")
		return err
	}

	cartItem.OrderItemID = &newOrderItem.ID
	cartItem.Status = 1
	if err := tx.Save(&cartItem).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(c, &currVariation, cartItem.VariationID); err != nil {
		helpers.ErrJSONResponse(c, http.StatusNotFound, "Product variation not found")
		return err
	}

	if currVariation.Quantity < cartItem.Quantity {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "Insufficient product quantity")
		return errors.New("insufficient quantity")
	}

	currVariation.Quantity -= cartItem.Quantity
	if err := tx.Save(&currVariation).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "Failed to update product variation")
		return err
	}

	newOrderItem.Items = append(newOrderItem.Items, cartItem)

	return nil
}

func GetOrders(c *gin.Context) {
	transformedResponse, err := helpers.GetOrderByStatus(c, false)
	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedResponse))
}

func FinishOrder(c *gin.Context) {
	var body struct {
		OrderID string `json:"order_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var orderItem models.OrderItem
	tx := database.DB.Begin()
	if err := tx.First(&orderItem, "id = ?", body.OrderID).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if orderItem.Status == 3 {
		orderItem.Status = 4
		if err := tx.Save(&orderItem).Error; err != nil {
			tx.Rollback()
			helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		tx.Rollback()
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "invalid order status")
		return
	}

	tx.Commit()
	helpers.JSONResponse(c, "")
}
