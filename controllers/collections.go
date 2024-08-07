package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func PublicCollections(ctx *gin.Context) {
	var collections []models.Collection

	if err := database.DB.Find(&collections).Where("status = ?", 1).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	transformedCollections := []map[string]interface{}{}
	for _, v := range collections {
		newCollection := map[string]interface{}{
			"id":     v.ID,
			"name":   v.Name,
			"images": v.Images,
		}

		transformedCollections = append(transformedCollections, newCollection)
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(transformedCollections))
}
func AddCollection(ctx *gin.Context) {
	var body models.CollectionPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	if len(body.Images) < 1 {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "images is required")
		return
	}

	newCollection := models.Collection{
		ID:     helpers.NewUUID(),
		Name:   body.Name,
		Images: body.Images,
	}
	if err := helpers.CreateNewData(ctx, &newCollection); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}

func GetCollection(ctx *gin.Context) {
	var collections []models.Collection

	// NO NEED TO HANDLE ERROR HERE BECAUSE COLLECTION IS EXISTENT
	helpers.GetTableByModel(ctx, &collections)
	helpers.JSONResponse(ctx, "", helpers.DataHelper(&collections))
}

func UpdateCollection(ctx *gin.Context) {
	var body models.CollectionPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currCollection models.Collection
	if err := helpers.GetCurrentByID(ctx, &currCollection, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE COLLECTION IS EXISTENT
	helpers.UpdateByModel(ctx, &currCollection, models.Collection{Name: body.Name, Images: body.Images})
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
	if err := helpers.GetCurrentByID(ctx, &currCollection, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE COLLECTION IS EXISTENT
	helpers.DeleteByModel(ctx, &currCollection)
	helpers.JSONResponse(ctx, "")
}
