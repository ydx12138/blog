package middleware

import (
	"blog/internal/utils"
	"blog/pkg/code"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// 用户token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从context中得到token
		token := utils.GetTokenFromContext(c)
		if token == "" {
			response.ErrWithMsg(code.Unauthorized, c)
			c.Abort()
			return
		}
		//解析token，得到Data
		data, err := utils.GetDataFromToken(token)
		if data == nil || err != nil {
			response.ErrWithMsg(code.Unauthorized, c)
			c.Abort()
			return
		}
		//解析token，得到claim
		var claim *utils.CustomClaims
		if claim = utils.GetClaimFromData(data); claim == nil {
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
			response.ErrWithMsg(code.AccessTokenExpired, c)
			c.Abort()
			return
		}
	}
	//	return func(c *gin.Context) {
	//		//从token中解析出CustomClaims结构体(UserId，Role，Type),同时检测token的有效
	//		claims, ok := utils.ParseClaims(c)
	//		if !ok {
	//			return
	//		}
	//		//获得token
	//		token := utils.GetToken(c)
	//		if token == "" {
	//			response.ErrWithMsg(code.Unauthorized, c)
	//			c.Abort()
	//			return
	//		}
	//		//获得CustomClaims
	//		claim, err := utils.GetClaimFromToken(token)
	//		if err != nil {
	//			response.ErrWithMsg(code.Unauthorized, c)
	//			c.Abort()
	//			return
	//		}
	//		//如果类型是access:
	//		//
	//		c.Set("userID", claims.UserID)
	//		c.Set("role", claims.Role)
	//		c.Next()
	//	}
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
