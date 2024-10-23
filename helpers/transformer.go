package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/models"
	"github.com/forthedreamers-server/tokens"
	"github.com/gin-gonic/gin"
)

func UserGetTokenResponse(c *gin.Context, user *models.Users) models.CredentialResponse {
	token, err := tokens.CreateAndSignJWT(&user.ID)
	if err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return models.CredentialResponse{}
	}
	user.Token = token
	database.DB.Save(&user)

	newUserResponse := models.UsersResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Image:     user.Image,
		Username:  user.Username,
		Phone:     user.Phone,
	}

	partToken := strings.Split(token, ".")
	userRes := models.CredentialResponse{
		UserInfo: newUserResponse,
		Token:    partToken[1],
	}

	return userRes
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

func TransformOrderItem(orderItem models.OrderItem, productsMap map[string]models.Product, variationsMap map[string]models.ProductVariation, currAddress models.AddressItem) models.OrderItemResponse {
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
			CreatedAt:   cartItem.CreatedAt,
			Image:       product.Images[0],
			IsReviewed:  cartItem.IsReviewed,
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

func TransformCartItems(c *gin.Context, status int) ([]models.CartItemResponse, error) {
	userID := GetCurrUserToken(c).ID

	var userCarts []models.UserCart
	if err := database.DB.Where("user_id = ?", userID).Find(&userCarts).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	var cartItemIDs []string
	for _, uc := range userCarts {
		cartItemIDs = append(cartItemIDs, uc.CartItemID)
	}

	var cartItems []models.CartItem
	if err := database.DB.
		Where("id IN ?", cartItemIDs).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&cartItems).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	var transformedCartItems []models.CartItemResponse
	for _, v := range cartItems {
		var transformedCartItem models.CartItemResponse
		var product models.Product
		var variation models.ProductVariation
		GetCurrentByID(c, &product, v.ProductID)

		if v.VariationID != "" {
			GetCurrentByID(c, &variation, v.VariationID)
			transformedCartItem.Size = variation.Size
			transformedCartItem.Color = variation.Color
			transformedCartItem.Price = variation.Price
		}

		transformedCartItem.ID = v.ID
		transformedCartItem.ProductName = product.Name
		transformedCartItem.Quantity = v.Quantity
		transformedCartItem.Image = []string{product.Images[0]}
		transformedCartItem.ProductID = v.ProductID

		transformedCartItems = append(transformedCartItems, transformedCartItem)
	}

	return transformedCartItems, nil
}

func ToJSON(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshalling to JSON:", err)
		return "{}"
	}
	return string(jsonData)
}

func Keys(m map[string]struct{}) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func GetOrderByStatus(c *gin.Context, isReview bool) (interface{}, error) {
	userID := GetCurrUserToken(c).ID

	query := database.DB.Where("user_id = ?", userID)
	if isReview {
		query = query.Where("status = ?", 4)
	}

	var orderItems []models.OrderItem
	if err := query.Preload("Items").Find(&orderItems).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return nil, err
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
	if err := database.DB.Find(&products, "id IN ?", Keys(productIDs)).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	var variations []models.ProductVariation
	if err := database.DB.Find(&variations, "id IN ?", Keys(variationIDs)).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	var addresses []models.AddressItem
	if err := database.DB.Find(&addresses, "id IN ?", Keys(addressIDs)).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return nil, err
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
		transformedResponse = append(transformedResponse, TransformOrderItem(orderItem, productsMap, variationsMap, currAddress))
	}

	if isReview {
		var itemResponse []models.ItemResponse
		for _, v := range transformedResponse {
			itemResponse = append(itemResponse, v.Items...)
		}

		return itemResponse, nil
	} else {
		return transformedResponse, nil
	}
}
