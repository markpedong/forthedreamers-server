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
	if err := database.DB.Preload("ProductVariations").Order("created_at DESC").Find(&products).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(&products))
}

func AddProducts(ctx *gin.Context) {
	var body models.ProductPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	newProduct := models.Product{
		ID:           helpers.NewUUID(),
		Name:         body.Name,
		Images:       body.Images,
		Description:  body.Description,
		CollectionID: body.CollectionID,
	}

	var variations []models.ProductVariation
	for _, v := range body.Variations {
		variation := models.ProductVariation{
			Size:      v.Size,
			Color:     v.Color,
			Price:     v.Price,
			Quantity:  v.Quantity,
			ID:        helpers.NewUUID(),
			ProductID: newProduct.ID,
		}

		variations = append(variations, variation)
	}

	newProduct.ProductVariations = variations

	if err := database.DB.Create(&newProduct).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}

func UpdateProducts(ctx *gin.Context) {
	var body models.ProductPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := database.DB.Find(&currProduct, "id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := database.DB.Model(&currProduct).Updates(models.Product{
		Name:         body.Name,
		Images:       body.Images,
		Description:  body.Description,
		CollectionID: body.CollectionID,
	}).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}
