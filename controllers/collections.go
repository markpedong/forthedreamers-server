package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func AddCollection(ctx *gin.Context) {
	var body models.CollectionPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	newCollection := models.Collection{
		ID:     helpers.NewUUID(),
		Name:   body.Name,
		Images: body.Images,
	}

	if err := database.DB.Create(&newCollection).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}

func GetCollection(ctx *gin.Context) {
	var collections []models.Collection
	if err := database.DB.Order("created_at DESC").Find(&collections).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&collections))
}

func UpdateCollection(ctx *gin.Context) {
	var body models.CollectionPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currCollection models.Collection
	if err := database.DB.Find(&currCollection, "id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := database.DB.Model(&currCollection).Updates(models.Collection{Name: body.Name, Images: body.Images}).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")

}

func DeleteCollection(ctx *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currCollection models.Collection
	if err := database.DB.First(&currCollection, "id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := database.DB.Delete(&currCollection).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}
