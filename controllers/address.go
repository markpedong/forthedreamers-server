package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetAddress(c *gin.Context) {
	userID := helpers.GetCurrUserToken(c).ID

	var address []models.AddressItem
	if err := database.DB.Where("user_id = ?", userID).Find(&address).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var transformedAddress []models.AddressItemReponse
	for _, item := range address {
		transformedAddress = append(transformedAddress, models.AddressItemReponse{
			ID:        item.ID,
			FirstName: item.FirstName,
			LastName:  item.LastName,
			Phone:     item.Phone,
			Address:   item.Address,
			IsDefault: item.IsDefault,
		})
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedAddress))
}

func AddAddress(c *gin.Context) {
	userID := helpers.GetCurrUserToken(c).ID

	var body models.AddressItem
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var existingAddress []models.AddressItem
	if qty := database.DB.Where("user_id = ?", userID).Find(&existingAddress).RowsAffected; qty >= 5 {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "You can only have a maximum of 5 addresses")
		return
	}

	hasDefault := false
	hasReturn := false
	hasPickup := false

	for _, item := range existingAddress {
		switch item.IsDefault {
		case 1:
			hasDefault = true
		case 2:
			hasPickup = true
		case 3:
			hasReturn = true
		}
	}

	switch body.IsDefault {
	case 1:
		if hasDefault {
			helpers.ErrJSONResponse(c, http.StatusBadRequest, "Only one default address is allowed.")
			return
		}
	case 2:
		if hasPickup {
			helpers.ErrJSONResponse(c, http.StatusBadRequest, "Only one pickup address is allowed.")
			return
		}
	case 3:
		if hasReturn {
			helpers.ErrJSONResponse(c, http.StatusBadRequest, "Only one return address is allowed.")
			return
		}
	case 0:
	default:
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "Invalid address type.")
		return
	}

	address := models.AddressItem{
		ID:        helpers.NewUUID(),
		UserID:    userID,
		Phone:     body.Phone,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Address:   body.Address,
		IsDefault: body.IsDefault,
	}

	if err := helpers.CreateNewData(c, &address); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}

func UpdateAddress(c *gin.Context) {
	var body models.AddressItem
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currAddress models.AddressItem
	if err := helpers.GetCurrentByID(c, &currAddress, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE ADDRESS IS EXISTENT
	helpers.UpdateByModel(c, &currAddress, models.AddressItem{
		Phone:     body.Phone,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Address:   body.Address,
		IsDefault: currAddress.IsDefault,
	})

	helpers.JSONResponse(c, "", helpers.DataHelper(&currAddress))
}

func DeleteAddress(c *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currAddress models.AddressItem
	if err := helpers.GetCurrentByID(c, &currAddress, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE ADDRESS IS EXISTENT
	helpers.DeleteByModel(c, &currAddress)
	helpers.JSONResponse(c, "")
}
