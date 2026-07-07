package middleware

import (
	"blog/internal/utils"
	"blog/pkg/code"
	"blog/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// 用户token
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

// 管理员token
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

// parseClaims 从请求 Header 中解析 JWT Token，并返回解析后的自定义 Claims
// 返回值：
//   - *utils.CustomClaims：解析成功后的用户信息
//   - bool：是否解析成功（true=成功，false=失败）
func parseClaims(c *gin.Context) (*utils.CustomClaims, bool) {
	// 从请求头中获取 Authorization 字段
	// 一般格式：Bearer xxx.yyy.zzz
	tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
	// 如果没有携带 token，直接返回未授权
	if tokenStr == "" {
		response.ErrWithMsg(code.Unauthorized, c)
		c.Abort()
		return nil, false
	}
	// 兼容 "Bearer xxx" 格式（忽略大小写）
	// 如果有 Bearer 前缀，则截取真实 token
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	}
	// 解析 JWT Token，获取 claims（用户信息 + 自定义字段）
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		// token 无效或过期，返回未授权
		response.ErrWithMsg(code.Unauthorized, c)
		c.Abort()
		return nil, false
	}
	// 解析成功，返回 claims
	return claims, true
}

// 刷新或进入标签页的时候，如果有token，就请求，验证这个token是否还能用
func VerificationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		tokenStr := strings.TrimSpace(c.GetHeader("Authorization"))
		//如果有bearer 前缀，则切割一下
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			tokenStr = strings.TrimSpace(tokenStr[7:])
		}
		//token为空,请先登录，但是前端如果没有token，就不会发这个请求，所以这一步应该用不到
		if tokenStr == "" {
			response.ErrWithMsg(code.Unauthorized, c)
			return
		}
		//

	}
}

//双token
