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

func transformOrderItem(orderItem models.OrderItem, productsMap map[string]models.Product, variationsMap map[string]models.ProductVariation, currAddress models.AddressItem) models.OrderItemResponse {
	var itemsResponse []models.ItemResponse

	for _, cartItem := range orderItem.Items {
		product := productsMap[cartItem.ProductID]
		variation := variationsMap[cartItem.VariationID]

		itemsResponse = append(itemsResponse, models.ItemResponse{
			ID:          cartItem.ID,
			Quantity:    cartItem.Quantity,
			ProductName: product.Name,
			ProductID:   cartItem.ProductID,
			Size:        variation.Size,
			Color:       variation.Color,
			Price:       variation.Price,
			Image:       product.Images[0],
		})
	}

	totalPrice := calculateTotalPrice(orderItem.Items, variationsMap)

	return models.OrderItemResponse{
		ID:            orderItem.ID,
		TotalPrice:    totalPrice,
		PaymentMethod: orderItem.PaymentMethod,
		Items:         itemsResponse,
		Address: models.AddressItemReponse{
			FirstName: currAddress.FirstName,
			LastName:  currAddress.LastName,
			Phone:     currAddress.Phone,
			Address:   currAddress.Address,
			IsDefault: currAddress.IsDefault,
		},
		CreatedAt: orderItem.CreatedAt,
		Status:    orderItem.Status,
	}
}

func calculateTotalPrice(items []models.CartItem, variationsMap map[string]models.ProductVariation) int {
	total := 0
	for _, item := range items {
		if product, exists := variationsMap[item.VariationID]; exists {
			total += product.Price * item.Quantity
		}
	}
	return total
}

func GetOrders(c *gin.Context) {
	userID := helpers.GetCurrUserToken(c).ID

	var orderItems []models.OrderItem
	if err := database.DB.Preload("Items").Find(&orderItems, "user_id = ?", userID).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	productIDs := make(map[string]struct{})
	variationIDs := make(map[string]struct{})
	addressIDs := make(map[string]struct{})

	for _, orderItem := range orderItems {
		addressIDs[orderItem.AddressID] = struct{}{}
		for _, cartItem := range orderItem.Items {
			productIDs[cartItem.ProductID] = struct{}{}
			variationIDs[cartItem.VariationID] = struct{}{}
		}
	}

	var products []models.Product
	if err := database.DB.Find(&products, "id IN ?", helpers.Keys(productIDs)).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var variations []models.ProductVariation
	if err := database.DB.Find(&variations, "id IN ?", helpers.Keys(variationIDs)).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var addresses []models.AddressItem
	if err := database.DB.Find(&addresses, "id IN ?", helpers.Keys(addressIDs)).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	productsMap := make(map[string]models.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	variationsMap := make(map[string]models.ProductVariation)
	for _, variation := range variations {
		variationsMap[variation.ID] = variation
	}

	addressMap := make(map[string]models.AddressItem)
	for _, address := range addresses {
		addressMap[address.ID] = address
	}

	var transformedResponse []models.OrderItemResponse
	for _, orderItem := range orderItems {
		currAddress := addressMap[orderItem.AddressID]
		transformedResponse = append(transformedResponse, transformOrderItem(orderItem, productsMap, variationsMap, currAddress))
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedResponse))
}
