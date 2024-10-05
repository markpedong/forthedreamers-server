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
	var body models.AddressItem
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	address := models.AddressItem{
		ID:        helpers.NewUUID(),
		UserID:    helpers.GetCurrUserToken(c).ID,
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
