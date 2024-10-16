package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var Validate = validator.New()

func ErrJSONResponse(c *gin.Context, status int, message string, jsonData ...map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": false,
		"status":  status,
	})
}

func JSONResponse(c *gin.Context, optionalMessage string, jsonData ...map[string]interface{}) {
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

	c.JSON(http.StatusOK, response)
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

func BindValidateJSON(c *gin.Context, body interface{}) error {
	if err := c.ShouldBindJSON(body); err != nil {
		ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	if err := Validate.Struct(body); err != nil {
		ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

func BindValidateQuery(c *gin.Context, body interface{}) error {
	if err := c.BindQuery(body); err != nil {
		ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	if err := Validate.Struct(body); err != nil {
		ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

// func GetTableByModelStatusON(c *gin.Context, model interface{}, preload ...string) interface{} {
// 	query := database.DB.Where("status = ?", 1).Order("created_at DESC")

// 	if len(preload) > 0 {
// 		for _, p := range preload {
// 			query = query.Preload(p)
// 		}
// 	}

// 	if err := query.Find(model).Error; err != nil {
// 		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
// 		return nil
// 	}

// 	return model
// }

func CreateNewData(c *gin.Context, model interface{}) error {
	if err := database.DB.Create(model).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func GetTableByModel(c *gin.Context, model interface{}, preload ...string) error {
	query := database.DB.Order("created_at DESC")

	if len(preload) > 0 {
		for _, p := range preload {
			query = query.Preload(p)
		}
	}

	if err := query.Find(model).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func GetCurrentByID(c *gin.Context, model interface{}, ID string, preload ...string) error {
	query := database.DB
	if len(preload) > 0 {
		for _, p := range preload {
			query = query.Preload(p)
		}
	}

	result := query.Find(model, "id = ?", ID)
	if result.Error != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, result.Error.Error())
		return result.Error
	}

	if result.RowsAffected == 0 {
		ErrJSONResponse(c, http.StatusNotFound, "record not found")
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteByModel(c *gin.Context, model interface{}) error {
	if err := database.DB.Delete(model).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func UpdateByModel(c *gin.Context, model interface{}, newValues interface{}) error {
	if err := database.DB.Model(model).Updates(newValues).Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func ToggleModelByID(c *gin.Context, model interface{}, id string) error {
	if err := database.DB.Model(model).
		Where("id = ?", id).
		Update("status", gorm.Expr("1 - status")).
		Error; err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}

func GetCurrUserToken(c *gin.Context, preload ...string) models.Users {
	token := c.GetHeader("Token")

	if token == "" {
		ErrJSONResponse(c, http.StatusUnauthorized, "Token is missing")
		c.Abort()
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
		c.Abort()
		return models.Users{}
	}

	return user
}

func ToJSON(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshalling to JSON:", err)
		return "{}"
	}
	return string(jsonData)
}

func Keys(m map[string]struct{}) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
