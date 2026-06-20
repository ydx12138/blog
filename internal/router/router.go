package router

import (
	"blog/internal/middleware"
	"blog/internal/service/admin"
	"blog/internal/service/user"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Register 注册所有路由
func Register() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	// 静态文件 - 上传图片
	r.Static("/uploads", "./uploads")

	// ========== 公开API ==========
	api := r.Group("/api")
	public := api.Group("")

	// 文章
	public.GET("/articles", user.GetArticles)
	public.GET("/articles/detail", user.GetArticle)
	public.GET("/articles/search", user.SearchArticle)

	// 分类
	public.GET("/categories", user.GetCategories)
	public.GET("/categories/articles", user.GetCategoryArticles)

	// 评论（公开）
	public.GET("/comments", user.GetComments)

	// 用户认证（公开）
	public.POST("/register", user.Register)
	public.POST("/login", user.Login)

	// 用户认证路由（需JWT）
	apiAuth := api.Group("")
	apiAuth.Use(middleware.JWTAuth())
	apiAuth.POST("/comments", user.CreateComment)

	// ========== 管理员API（/api/admin） ==========
	adminGroup := r.Group("/api/admin")
	adminGroup.POST("/login", admin.Login)

	adminAuth := adminGroup.Group("")
	adminAuth.Use(middleware.JWTAuth())

	adminAuth.GET("/dashboard", admin.GetDashboard)
	adminAuth.GET("/articles", admin.GetArticles)
	adminAuth.GET("/articles/:id", admin.GetArticle)
	adminAuth.POST("/articles", admin.CreateArticle)
	adminAuth.PUT("/articles/:id", admin.UpdateArticle)
	adminAuth.DELETE("/articles/:id", admin.DeleteArticle)
	adminAuth.GET("/drafts", admin.GetDrafts)
	adminAuth.PUT("/articles/:id/publish", admin.PublishArticle)
	adminAuth.POST("/upload", admin.UploadImage)
	adminAuth.GET("/comments/pending", admin.GetPendingComments)
	adminAuth.PUT("/comments/:id/approve", admin.ApproveComment)
	adminAuth.PUT("/comments/:id/reject", admin.RejectComment)
	adminAuth.DELETE("/comments/:id", admin.DeleteComment)
	adminAuth.GET("/users", admin.GetUsers)
	adminAuth.PUT("/users/:id/ban", admin.BanUser)
	adminAuth.PUT("/users/:id/unban", admin.UnbanUser)
	adminAuth.DELETE("/users/:id", admin.DeleteUser)

	// ========== SPA前端兜底 ==========
	// 依次查找可能的前端构建目录
	distDir := findFrontendDist()
	if distDir != "" {
		zap.L().Info("前端静态文件目录: " + distDir)
		r.Static("/assets", filepath.Join(distDir, "assets"))
		r.StaticFile("/favicon.svg", filepath.Join(distDir, "favicon.svg"))
		r.StaticFile("/index.html", filepath.Join(distDir, "index.html"))

		// NoRoute: 所有非API路径返回 index.html（SPA兜底）
		r.NoRoute(func(c *gin.Context) {
			// 已经匹配到API路由的不会走到这里
			c.File(filepath.Join(distDir, "index.html"))
		})
	}

	return r
}

// findFrontendDist 查找前端构建目录
func findFrontendDist() string {
	candidates := []string{
		"../vue6122/dist",
		"dist",
		"./dist",
	}
	for _, d := range candidates {
		abs, err := filepath.Abs(d)
		if err != nil {
			continue
		}
		idx := filepath.Join(abs, "index.html")
		if _, err := os.Stat(idx); err == nil {
			return abs
		}
	}
	return ""
}
