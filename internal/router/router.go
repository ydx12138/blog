package router

import (
	"blog/internal/middleware"
	"blog/internal/service/admin"
	"blog/internal/service/user"

	"github.com/gin-gonic/gin"
)

// Register 注册所有路由
func Register() *gin.Engine {
	r := gin.Default()
	// 全局中间件
	// 设置 CORS 中间件（全局生效，白名单来自配置）
	r.Use(middleware.CorsMiddleware())

	// 第1组：api公开接口
	api := r.Group("/api")
	public := api.Group("")

	public.GET("/articles", user.GetArticles)
	public.GET("/articles/detail", user.GetArticle)
	public.GET("/articles/search", user.SearchArticle)

	//第二组: admin管理员接口

	notpublic := r.Group("/admin")
	notpublic.POST("/login", admin.Login)

	return r
}
