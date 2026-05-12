package auth

import (
	jws "backend/util/jwt/jws"
	rest "backend/util/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			rest.ResponseMessage(ctx, http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		if _, err := jws.ParseToken(auth, false); err != nil {
			rest.ResponseMessage(ctx, http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
