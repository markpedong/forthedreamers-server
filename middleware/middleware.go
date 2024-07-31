package middleware

import (

	// "github.com/forthedreamers-server/helpers"

	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/forthedreamers-server/controllers"
	"github.com/forthedreamers-server/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("Auth")
	if err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "")
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HMAC_SECRET")), nil
	})
	if err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "JWT Claims failed")
		return
	}

	if claims["ttl"].(int64) > time.Now().Unix() {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "JWT token expired")
		return
	}

	user := controllers.GetUserByID(claims["userID"].(string), ctx)
	if user.ID == "" {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, "Could not find the User")
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}
