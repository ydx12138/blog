package router

import (
	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

// Register 注册所有路由
func Register() *gin.Engine {
	r := gin.Default()
	// 全局中间件
	//r.Use(middleware.CORSMiddleware())
	api := r.Group("/api")

	// 第1组：公开接口
	public := api.Group("")

	public.GET("/articles", service.GetArticles)

	return r
}
