package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func AddTestimonials(ctx *gin.Context) {
	var body models.Testimonials
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	newTestimonial := models.Testimonials{
		ID:        helpers.NewUUID(),
		Title:     body.Title,
		Author:    body.Author,
		Status:    body.Status,
		ProductID: body.ProductID,
	}
	if err := helpers.CreateNewData(ctx, &newTestimonial); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}

func GetTestimonials(ctx *gin.Context) {
	var testimonials []models.Testimonials
	if err := database.DB.Find(&testimonials).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(testimonials))
}

func PublicTestimonials(ctx *gin.Context) {
	var testimonials []models.Testimonials
	if err := database.DB.Where("status = ?", 1).Find(&testimonials).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	transformedTestimonials := []map[string]interface{}{}
	for _, v := range testimonials {
		newTestimonial := map[string]interface{}{
			"id":         v.ID,
			"product_id": v.ProductID,
			"author":     v.Author,
			"title":      v.Title,
		}
		transformedTestimonials = append(transformedTestimonials, newTestimonial)
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(transformedTestimonials))
}
