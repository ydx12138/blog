package middleware

import (
	"blog/internal/utils"
	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 用户token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从context中得到token
		token := utils.GetTokenFromContext(c)
		if token == "" {
			zap.L().Info("从context中没有得到token")
			response.ErrWithMsg(code.Unauthorized, c)
			c.Abort()
			return
		}
		//解析token，得到Data
		data, err := utils.GetDataFromToken(token)
		if data == nil || err != nil {
			zap.L().Info("从token中没有得到Data")
			response.ErrWithMsg(code.Unauthorized, c)
			c.Abort()
			return
		}
		//解析data，得到claim
		var claim *utils.CustomClaims
		if claim = utils.GetClaimFromData(data); claim == nil {
			zap.L().Info("从data中没有得到claim")
			response.ErrWithMsg(code.Unauthorized, c)
			c.Abort()
			return
		}
		//如果type==access,data.Valid==true有效，则通过
		if claim.Type == "access" && data.Valid {
			c.Set("userID", claim.UserID)
			c.Set("role", claim.Role)
			c.Set("type", claim.Type)
			return
		}
		//如果type==access,无效，401，Abort+return
		if claim.Type == "access" && data.Valid == false {
			zap.L().Info("accessToken无效")
			response.ErrWithMsg(code.AccessTokenExpired, c)
			c.Abort()
			return
		}
	}

}

// 检测管理员token
func JWTAuthForAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := utils.ParseClaims(c)
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

//双token
