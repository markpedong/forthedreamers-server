package middleware

import (

	// "github.com/forthedreamers-server/helpers"

	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/forthedreamers-server/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication(ctx *gin.Context) {
	tokenStr := ctx.Request.Header.Get("Token")
	if tokenStr == "" {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, "Token is missing")
		ctx.Abort()
		return
	}

	user := helpers.GetCurrUserToken(ctx)
	token, err := jwt.Parse(user.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HMAC_SECRET")), nil
	})
	if err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, err.Error())
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, "JWT Claims failed")
		ctx.Abort()
		return
	}

	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, "JWT token expired")
		ctx.Abort()
		return
	}

	if user.ID == "" {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, "Could not find the User")
		ctx.Abort()
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}
