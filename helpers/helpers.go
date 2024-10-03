package helpers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var Validate = validator.New()

func ErrJSONResponse(ctx *gin.Context, status int, message string, jsonData ...map[string]interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": false,
		"status":  status,
	})
}

func JSONResponse(ctx *gin.Context, optionalMessage string, jsonData ...map[string]interface{}) {
	message := "success"

	if optionalMessage != "" {
		message = optionalMessage
	}

	response := gin.H{
		"message": message,
		"success": true,
		"status":  http.StatusOK,
	}

	if len(jsonData) > 0 && jsonData[0] != nil {
		for key, value := range jsonData[0] {
			response[key] = value
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func DataHelper(data interface{}) map[string]interface{} {
	q := map[string]interface{}{
		"data": data,
	}

	return q

}

func NewUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func BindValidateJSON(ctx *gin.Context, body interface{}) error {
	if err := ctx.ShouldBindJSON(body); err != nil {
		ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return err
	}

	if err := Validate.Struct(body); err != nil {
		ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

func BindValidateQuery(ctx *gin.Context, body interface{}) error {
	if err := ctx.BindQuery(body); err != nil {
		ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return err
	}

	if err := Validate.Struct(body); err != nil {
		ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// func GetTableByModelStatusON(ctx *gin.Context, model interface{}, preload ...string) interface{} {
// 	query := database.DB.Where("status = ?", 1).Order("created_at DESC")

// 	if len(preload) > 0 {
// 		for _, p := range preload {
// 			query = query.Preload(p)
// 		}
// 	}

// 	if err := query.Find(model).Error; err != nil {
// 		ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
// 		return nil
// 	}

// 	return model
// }

func CreateNewData(ctx *gin.Context, model interface{}) error {
	if err := database.DB.Create(model).Error; err != nil {
		ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func GetTableByModel(ctx *gin.Context, model interface{}, preload ...string) error {
	query := database.DB.Order("created_at DESC")

	if len(preload) > 0 {
		for _, p := range preload {
			query = query.Preload(p)
		}
	}

	if err := query.Find(model).Error; err != nil {
		ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func GetCurrentByID(ctx *gin.Context, model interface{}, ID string, preload ...string) error {
	query := database.DB
	if len(preload) > 0 {
		for _, p := range preload {
			query = query.Preload(p)
		}
	}

	result := query.Find(model, "id = ?", ID)
	if result.Error != nil {
		ErrJSONResponse(ctx, http.StatusInternalServerError, result.Error.Error())
		return result.Error
	}

	if result.RowsAffected == 0 {
		ErrJSONResponse(ctx, http.StatusNotFound, "record not found")
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteByModel(ctx *gin.Context, model interface{}) error {
	if err := database.DB.Delete(model).Error; err != nil {
		ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func UpdateByModel(ctx *gin.Context, model interface{}, newValues interface{}) error {
	if err := database.DB.Model(model).Updates(newValues).Error; err != nil {
		ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func ToggleModelByID(ctx *gin.Context, model interface{}, id string) error {
	if err := database.DB.Model(model).
		Where("id = ?", id).
		Update("status", gorm.Expr("1 - status")).
		Error; err != nil {
		ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func GetCurrUserToken(ctx *gin.Context, preload ...string) models.Users {
	token := ctx.GetHeader("Token")

	if token == "" {
		ErrJSONResponse(ctx, http.StatusUnauthorized, "Token is missing")
		ctx.Abort()
		return models.Users{}
	}

	var user models.Users
	query := database.DB
	if len(preload) > 0 {
		for _, p := range preload {
			query = query.Preload(p)
		}
	}
	if err := query.
		Where("SPLIT_PART(token, '.', 2) = ?", token).
		First(&user).Error; err != nil {
		ctx.Abort()
		return models.Users{}
	}

	return user
}
