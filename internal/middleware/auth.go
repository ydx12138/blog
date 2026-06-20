package middleware

import (
	"blog/internal/utils"

	"github.com/gin-gonic/gin"
)

// 验证token中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "未登录",
			})
			return
		}
		//解析token
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "Token无效",
			})
			return
		}
		// 把userID和role传下去
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
