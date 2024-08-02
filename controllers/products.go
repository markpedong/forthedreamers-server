package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func PublicProductDetails(ctx *gin.Context) {
	var body struct {
		ID string `json:"product_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ID, "Variations"); err != nil {
		return
	}

	var productVariations []models.VariationResponse
	for _, v := range currProduct.Variations {
		newProductVariation := models.VariationResponse{
			ID:       v.ID,
			Size:     v.Size,
			Color:    v.Color,
			Price:    v.Price,
			Quantity: v.Quantity,
		}

		productVariations = append(productVariations, newProductVariation)
	}

	transformedProduct := models.ProductResponse{
		ID:           currProduct.ID,
		Name:         currProduct.Name,
		Description:  currProduct.Description,
		CollectionID: currProduct.CollectionID,
		Images:       currProduct.Images,
		Varitions:    productVariations,
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(transformedProduct))
}

func PublicProducts(ctx *gin.Context) {
	var products []models.Product
	helpers.GetTableByModelStatusON(ctx, &products, "Variations")

	var productResponse []models.ProductResponse
	for _, v := range products {
		var variations []models.VariationResponse
		for _, q := range v.Variations {
			newVariation := models.VariationResponse{
				ID:       q.ID,
				Size:     q.Size,
				Color:    q.Color,
				Price:    q.Price,
				Quantity: q.Quantity,
			}
			variations = append(variations, newVariation)
		}

		newProductResponse := models.ProductResponse{
			ID:           v.ID,
			Name:         v.Name,
			CollectionID: v.CollectionID,
			Images:       v.Images,
			Description:  v.Description,
			Varitions:    variations,
		}

		productResponse = append(productResponse, newProductResponse)
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(productResponse))
}

func GetProducts(ctx *gin.Context) {
	var products []models.Product

	//NO NEED TO HANDLE ERROR BECAUSE PRODUCT IS EXISTENT
	helpers.GetTableByModel(ctx, &products)
	helpers.JSONResponse(ctx, "", helpers.DataHelper(products))
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

	var variations []models.ProductVariation
	if err := database.DB.Find(&variations, "product_id = ?", currProduct.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range variations {
		helpers.DeleteByModel(ctx, &v)
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE CURRPRODUCT IS EXISTENT
	helpers.DeleteByModel(ctx, &currProduct)
	helpers.JSONResponse(ctx, "")
}

func ToggleProducts(ctx *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(ctx, &models.Product{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}
