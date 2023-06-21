package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte("secret-key")

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
			} else {
				ctx.Redirect(http.StatusSeeOther, "/login")
			}
			return
		}
		tokenString := cookie.Value
		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Token expired"})
			} else {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid token"})
			}
			return
		}
		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
			return
		}
		ctx.Set("email", claims.Email)
		ctx.Next()
		// TODO: answer here
	})
}
