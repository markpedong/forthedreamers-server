package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/forthedreamers-server/cloudinary"
	"github.com/forthedreamers-server/helpers"
	"github.com/gin-gonic/gin"
)

func VerifyPassword(expectedHashedPassword, givenPassword string) (bool, string) {
	// err := bcrypt.CompareHashAndPassword([]byte(expectedHashedPassword), []byte(givenPassword))
	err := expectedHashedPassword == givenPassword

	switch {
	case err:
		return false, "Password matched!"
	// case errors.Is(_, bcrypt.ErrMismatchedHashAndPassword):
	// 	return false, "Password is incorrect!"
	case !err:
		return true, "Password is incorrect!"
	default:
		// fmt.Printf("Password verification error: %s\n", err)
		return true, "Failed to verify password"
	}
}

func UploadImage(c *gin.Context) {
	form, err := c.FormFile("file")

	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uploadResult, err := cloudinary.CloudinaryService.Upload.Upload(c, form, uploader.UploadParams{
		Folder:         "forthedreamers",
		Transformation: "f_webp,q_auto:good,fl_lossy,c_fit",
	})

	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	imageRes := map[string]interface{}{
		"url":      uploadResult.URL,
		"fileName": uploadResult.OriginalFilename,
		"size":     uploadResult.Bytes,
	}

	helpers.JSONResponse(c, "upload successful!", helpers.DataHelper(imageRes))
}
