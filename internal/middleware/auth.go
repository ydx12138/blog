package middleware

import (
	"blog/internal/utils"
	"blog/pkg/code"
	"blog/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseClaims(c)
		if !ok {
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func JWTAuthForAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseClaims(c)
		if !ok {
			return
		}
		if claims.Role != "admin" {
			response.ErrWithMsg(code.Forbidden, c)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func parseClaims(c *gin.Context) (*utils.CustomClaims, bool) {
	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenStr == "" {
		response.ErrWithMsg(code.Unauthorized, c)
		c.Abort()
		return nil, false
	}
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	}
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		response.ErrWithMsg(code.Unauthorized, c)
		c.Abort()
		return nil, false
	}
	return claims, true
}
