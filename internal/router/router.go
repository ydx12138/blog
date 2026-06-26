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
	r.Use(middleware.CorsMiddleware())

	// 静态文件存放路径 - 上传图片
	r.Static("/uploads", "./uploads")

	// ========== 公开API ==========
	api := r.Group("/api")
	{
		public := api.Group("")

		// 文章列表
		public.GET("/articles", user.GetArticles)
		//文章详情
		public.GET("/articles/detail", user.GetArticle)
		//搜索文章
		public.GET("/articles/search", user.SearchArticle)

		// 分类
		public.GET("/categories", user.GetCategories)
		// xx分类的文章
		public.GET("/categories/articles", user.GetCategoryArticles)

		// 评论（公开）
		public.GET("/comments", user.GetComments)

		// 用户认证（公开）
		//注册
		public.POST("/register", user.Register)
		//登录
		public.POST("/login", user.Login)
		//点赞·
		public.POST("/articles/like", user.LikeArticle)
		//标签
		public.GET("/tags", user.GetTags)
	}
	{
		// 用户认证路由（需JWT）
		apiAuth := api.Group("")
		apiAuth.Use(middleware.JWTAuth())
		//增加评论，更新文章评论数
		apiAuth.POST("/comments", user.CreateComment)
	}

	{
		// ========== 管理员API ==========
		adminGroup := r.Group("/api/admin")
		//管理员登录
		adminGroup.POST("/login", admin.Login)
		adminAuth := adminGroup.Group("")
		adminAuth.Use(middleware.JWTAuthForAdmin())
		adminAuth.GET("/dashboard", admin.GetDashboard)
		adminAuth.GET("/articles", admin.GetArticles)
		adminAuth.GET("/articles/:id", admin.GetArticle)
		//创建文章
		adminAuth.POST("/articles", admin.CreateArticle)
		adminAuth.PUT("/articles/:id", admin.UpdateArticle)
		adminAuth.DELETE("/articles/:id", admin.DeleteArticle)
		adminAuth.GET("/drafts", admin.GetDrafts)
		adminAuth.PUT("/articles/:id/publish", admin.PublishArticle)
		adminAuth.POST("/upload", admin.UploadImage)
		// 评论管理
		adminAuth.GET("/comments", admin.GetAllComments)
		adminAuth.GET("/comments/pending", admin.GetPendingComments)
		adminAuth.PUT("/comments/:id/approve", admin.ApproveComment)
		adminAuth.PUT("/comments/:id/reject", admin.RejectComment)
		adminAuth.DELETE("/comments/:id", admin.DeleteComment)
		adminAuth.PUT("/comments/:id/status", admin.SetCommentStatus)
		// 用户管理
		adminAuth.GET("/users", admin.GetUsers)
		adminAuth.PUT("/users/:id/ban", admin.BanUser)
		adminAuth.PUT("/users/:id/unban", admin.UnbanUser)
		adminAuth.DELETE("/users/:id", admin.DeleteUser)
	}
	return r
}
