package middleware

import (

	// "github.com/forthedreamers-server/helpers"

	"net/http"

	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication(ctx *gin.Context) {
	clientToken := ctx.Request.Header.Get("token")
	if clientToken == "" {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, "No Authorization")
		ctx.Abort()
		return
	}

	claims, err := tokens.ValidateToken(clientToken)
	if err != "" {
		helpers.ErrJSONResponse(ctx, http.StatusUnauthorized, err)
		ctx.Abort()
		return
	}

	ctx.Set("email", claims.Email)
	ctx.Set("uid", claims.Uid)
	ctx.Next()
}
