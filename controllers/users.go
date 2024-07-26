package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
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
	if err := database.DB.Create(&newUser).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}

func GetUsers(ctx *gin.Context) {
	var users []models.Users
	if err := database.DB.Order("created_at DESC").Find(&users).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&users))
}

func UpdateUsers(ctx *gin.Context) {
	var body models.UserPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currUser models.Users
	if err := database.DB.Find(&currUser, "id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := database.DB.Model(&currUser).Updates(models.Users{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Image:     body.Image,
		Phone:     body.Phone,
		Address:   body.Address,
		Email:     body.Email,
		Username:  body.Username,
		Password:  body.Password,
	}).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}
