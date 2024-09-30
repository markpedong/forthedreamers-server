package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetAddress(ctx *gin.Context) {
	userID := helpers.GetCurrUserToken(ctx).ID

	var address []models.AddressItem
	if err := database.DB.Where("user_id = ?", userID).Find(&address).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
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
		})
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(transformedAddress))
}

func AddAddress(ctx *gin.Context) {
	var body models.AddressItem
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	userID := helpers.GetCurrUserToken(ctx).ID
	address := models.AddressItem{
		ID:        helpers.NewUUID(),
		UserID:    userID,
		Phone:     body.Phone,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Address:   body.Address,
	}
	if err := helpers.CreateNewData(ctx, &address); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}

func UpdateAddress(ctx *gin.Context) {
	var body models.AddressItem
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currAddress models.AddressItem
	if err := helpers.GetCurrentByID(ctx, &currAddress, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE ADDRESS IS EXISTENT
	helpers.UpdateByModel(ctx, &currAddress, models.AddressItem{
		Phone:     body.Phone,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Address:   body.Address,
		IsDefault: currAddress.IsDefault,
	})

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&currAddress))
}

func DeleteAddress(ctx *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currAddress models.AddressItem
	if err := helpers.GetCurrentByID(ctx, &currAddress, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE ADDRESS IS EXISTENT
	helpers.DeleteByModel(ctx, &currAddress)
	helpers.JSONResponse(ctx, "")
}
