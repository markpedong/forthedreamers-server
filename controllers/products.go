package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetProducts(ctx *gin.Context) {
	var products []models.Product
	if err := database.DB.Order("created_at DESC").Find(&products).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&products))
}

func AddProducts(ctx *gin.Context) {

}

func UpdateProducts(ctx *gin.Context) {

}
