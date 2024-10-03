package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetCart(c *gin.Context) {
	userID := helpers.GetCurrUserToken(c).ID

	var cartItems []models.CartItem
	if err := database.DB.Table("user_cart").
		Select("cart_item.*").
		Joins("JOIN cart_item ON user_cart.cart_item_id = cart_item.id").
		Where("user_cart.user_id = ?", userID).
		Find(&cartItems).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(cartItems))
}

func AddCartItem(c *gin.Context) {
	var body struct {
		ProductID   string `json:"product_id" validate:"required"`
		Quantity    int    `json:"quantity" validate:"required"`
		VariationID string `json:"variation_id" validate:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(c, &currVariation, body.VariationID); err != nil {
		return
	}

	newCartItem := models.CartItem{
		ID:          helpers.NewUUID(),
		ProductID:   body.ProductID,
		VariationID: currVariation.ID,
		Quantity:    body.Quantity,
	}

	userID := helpers.GetCurrUserToken(c).ID
	if err := CreateNewCartItem(c, userID, &newCartItem); err != nil {
		return
	}

	helpers.JSONResponse(c, "cart item added successfully")
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
