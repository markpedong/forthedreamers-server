package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetWebsiteData(ctx *gin.Context) {
	var website models.WebsiteData

	if err := database.DB.First(&website).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(website))
}

func UpdateWebsiteData(ctx *gin.Context) {
	var website models.WebsiteData
	if err := helpers.BindValidateJSON(ctx, &website); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Updates(&website).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "")
}

func PublicWebsite(ctx *gin.Context) {
	var website models.WebsiteData

	if err := database.DB.First(&website).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	transformedResponse := map[string]interface{}{
		"website_name":   website.WebsiteName,
		"promo_text":     website.PromoText,
		"marquee_text":   website.MarqueeText,
		"landing_image1": website.LandingImage1,
		"landing_image2": website.LandingImage2,
		"landing_image3": website.LandingImage3,
		"news_text":      website.NewsText,
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(transformedResponse))
}
