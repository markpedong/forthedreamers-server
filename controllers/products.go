package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetProducts(ctx *gin.Context) {
	var products []models.Product

	//NO NEED TO HANDLE ERROR BECAUSE PRODUCT IS EXISTENT
	helpers.GetTableByModel(ctx, &products, "ProductVariations")
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
	if err := helpers.CreateNewData(ctx, &newProduct); err != nil {
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
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ID); err != nil {
		return
	}

	//NO NEED TO HANDLE ERROR BECAUSE PRODUCT IS EXISTENT
	helpers.UpdateByModel(ctx, &currProduct, models.Product{
		Name:         body.Name,
		Images:       body.Images,
		Description:  body.Description,
		CollectionID: body.CollectionID,
	})
	helpers.JSONResponse(ctx, "")
}

func DeleteProducts(ctx *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE CURRPRODUCT IS EXISTENT
	helpers.DeleteByModel(ctx, &currProduct)
	helpers.JSONResponse(ctx, "")
}
