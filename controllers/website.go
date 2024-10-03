package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func GetWebsiteData(c *gin.Context) {
	var website models.WebsiteData

	if err := database.DB.First(&website).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(website))
}

func UpdateWebsiteData(c *gin.Context) {
	var website models.WebsiteData
	if err := helpers.BindValidateJSON(c, &website); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Updates(&website).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "")
}

func PublicWebsite(c *gin.Context) {
	var website models.WebsiteData
	var products []models.Product
	var collections []models.Collection

	if err := database.DB.First(&website).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := database.DB.Find(&products).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := database.DB.Find(&collections).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transformedResponse := map[string]interface{}{
		"website_name":      website.WebsiteName,
		"promo_text":        website.PromoText,
		"marquee_text":      website.MarqueeText,
		"landing_image1":    website.LandingImage1,
		"landing_image2":    website.LandingImage2,
		"landing_image3":    website.LandingImage3,
		"news_text":         website.NewsText,
		"product_length":    len(products),
		"default_pageSize":  website.DefaultPageSize,
		"collection_length": len(collections),
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedResponse))
}
