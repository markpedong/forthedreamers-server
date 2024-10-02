package controllers

import (
	"net/http"

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

func GetUserInfo(ctx *gin.Context) {
	user := helpers.GetCurrUserToken(ctx)

	newUserResponse := models.UsersResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Image:     user.Image,
		Username:  user.Username,
		Phone:     user.Phone,
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(newUserResponse))
}

func GetUsers(ctx *gin.Context) {
	var users []models.Users

	// NO NEED TO HANDLE ERROR HERE BECAUSE USER IS EXISTENT
	helpers.GetTableByModel(ctx, &users)
	helpers.JSONResponse(ctx, "", helpers.DataHelper(&users))
}

func UpdateUsers(ctx *gin.Context) {
	var body models.UpdateUserPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	currUser := helpers.GetCurrUserToken(ctx)
	if currUser.Password != body.OldPassword {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "old password is not correct")
		return
	}
	if body.NewPassword == body.OldPassword {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "old password and new password is same")
		return
	}

	helpers.UpdateByModel(ctx, &currUser, models.Users{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Image:     body.Image,
		Phone:     body.Phone,
		Email:     body.Email,
		Username:  body.Username,
		Password:  body.NewPassword,
	})
	helpers.JSONResponse(ctx, "Successfully updated!")
}

func DeleteUsers(ctx *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currUser models.Users
	if err := helpers.GetCurrentByID(ctx, &currUser, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE USER IS EXISTENT
	helpers.DeleteByModel(ctx, &currUser)
	helpers.JSONResponse(ctx, "")
}

func ToggleUsers(ctx *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(ctx, &models.Users{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}
