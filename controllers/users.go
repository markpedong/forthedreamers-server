package controllers

import (
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func AddUsers(ctx *gin.Context) {
	var body models.UserPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	newUser := models.Users{
		ID:        helpers.NewUUID(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Phone:     body.Phone,
		Address:   body.Address,
		Email:     body.Email,
		Image:     body.Image,
		Username:  body.Username,
		Password:  body.Password,
	}
	if err := helpers.CreateNewData(ctx, &newUser); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}

func GetUsers(ctx *gin.Context) {
	var users []models.Users

	// NO NEED TO HANDLE ERROR HERE BECAUSE USER IS EXISTENT
	helpers.GetTableByModel(ctx, &users)
	helpers.JSONResponse(ctx, "", helpers.DataHelper(&users))
}

func UpdateUsers(ctx *gin.Context) {
	var body models.UserPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currUser models.Users
	if err := helpers.GetCurrentByID(ctx, &currUser, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE USER IS EXISTENT
	helpers.UpdateByModel(ctx, &currUser, models.Users{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Image:     body.Image,
		Phone:     body.Phone,
		Address:   body.Address,
		Email:     body.Email,
		Username:  body.Username,
		Password:  body.Password,
	})
	helpers.JSONResponse(ctx, "")
}
