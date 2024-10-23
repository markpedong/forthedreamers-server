package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func AddUsers(c *gin.Context) {
	var body models.UserPayload
	if err := helpers.BindValidateJSON(c, &body); err != nil {
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
	if err := helpers.CreateNewData(c, &newUser); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}

func GetUserInfo(c *gin.Context) {
	user := helpers.GetCurrUserToken(c)

	newUserResponse := models.UsersResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Image:     user.Image,
		Username:  user.Username,
		Phone:     user.Phone,
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(newUserResponse))
}

func GetUsers(c *gin.Context) {
	var users []models.Users

	// NO NEED TO HANDLE ERROR HERE BECAUSE USER IS EXISTENT
	helpers.GetTableByModel(c, &users)
	helpers.JSONResponse(c, "", helpers.DataHelper(&users))
}

func UpdateUsers(c *gin.Context) {
	var body models.UpdateUserPayload
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	currUser := helpers.GetCurrUserToken(c)
	if currUser.Password != body.OldPassword {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "old password is not correct")
		return
	}
	if body.NewPassword == body.OldPassword {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "old password and new password is same")
		return
	}

	helpers.UpdateByModel(c, &currUser, models.Users{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Image:     body.Image,
		Phone:     body.Phone,
		Email:     body.Email,
		Username:  body.Username,
		Password:  body.NewPassword,
	})
	helpers.JSONResponse(c, "Successfully updated!")
}

func DeleteUsers(c *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currUser models.Users
	if err := helpers.GetCurrentByID(c, &currUser, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE USER IS EXISTENT
	helpers.DeleteByModel(c, &currUser)
	helpers.JSONResponse(c, "")
}

func ToggleUsers(c *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(c, &models.Users{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}
