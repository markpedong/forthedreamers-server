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

func Authentication(c *gin.Context) {
	tokenStr := c.Request.Header.Get("token")
	if tokenStr == "" {
		helpers.ErrJSONResponse(c, http.StatusUnauthorized, "Token is missing")
		c.Abort()
		return
	}

	user := helpers.GetCurrUserToken(c)
	fmt.Println("USER", user)
	token, err := jwt.Parse(user.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HMAC_SECRET")), nil
	})
	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.ErrJSONResponse(c, http.StatusUnauthorized, "JWT Claims failed")
		c.Abort()
		return
	}

	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		helpers.ErrJSONResponse(c, http.StatusUnauthorized, "JWT token expired")
		c.Abort()
		return
	}

	if user.ID == "" {
		helpers.ErrJSONResponse(c, http.StatusUnauthorized, "Could not find the User")
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
