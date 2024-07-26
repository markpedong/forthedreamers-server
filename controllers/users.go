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

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&body))
}

func GetUsers(ctx *gin.Context) {
	var users []models.Users
	if err := database.DB.Find(&users).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&users))
}
