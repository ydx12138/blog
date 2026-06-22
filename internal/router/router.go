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

		// 文章
		//用户获取文章列表
		public.GET("/articles", user.GetArticles)
		//用户获取文章详情
		public.GET("/articles/detail", user.GetArticle)
		//用户搜索文章
		public.GET("/articles/search", user.SearchArticle)

		// 分类
		//用户获取所有分类
		public.GET("/categories", user.GetCategories)
		//用户根据分类获取所有文章
		public.GET("/categories/articles", user.GetCategoryArticles)

		// 评论（公开）
		// 用户查看某篇文章的评论
		public.GET("/comments", user.GetComments)

		// 用户认证（公开）
		// 用户注册
		public.POST("/register", user.Register)
		//用户登录
		public.POST("/login", user.Login)
		//点赞文章（无需登录）
		public.POST("/articles/like", user.LikeArticle)
		//获取所有标签
		public.GET("/tags", user.GetTags)
	}
	{
		// 用户认证路由（需JWT）
		apiAuth := api.Group("")
		//验证token的中间件，
		apiAuth.Use(middleware.JWTAuth())
		//用户评论，需要登录
		apiAuth.POST("/comments", user.CreateComment)
	}

	{
		// ========== 管理员API（/api/admin） ==========
		adminGroup := r.Group("/api/admin")
		//管理员登录
		adminGroup.POST("/login", admin.Login)
		adminAuth := adminGroup.Group("")
		adminAuth.Use(middleware.JWTAuth())
		//获取数据面板
		adminAuth.GET("/dashboard", admin.GetDashboard)
		//管理员获取文章列表
		adminAuth.GET("/articles", admin.GetArticles)
		//管理员获取文章详情
		adminAuth.GET("/articles/:id", admin.GetArticle)
		//管理员新建文章
		adminAuth.POST("/articles", admin.CreateArticle)
		//管理员修改文章
		adminAuth.PUT("/articles/:id", admin.UpdateArticle)
		//管理员删除文章
		adminAuth.DELETE("/articles/:id", admin.DeleteArticle)
		//管理员查看草稿文章
		adminAuth.GET("/drafts", admin.GetDrafts)
		//管理员发布文章(将草稿文章发布出去)
		adminAuth.PUT("/articles/:id/publish", admin.PublishArticle)
		//上传文件，把文件放在uploads目录下
		adminAuth.POST("/upload", admin.UploadImage)
		//管理员获取所有待审核的评论
		adminAuth.GET("/comments/pending", admin.GetPendingComments)
		//管理员通过某个待审核的评论
		adminAuth.PUT("/comments/:id/approve", admin.ApproveComment)
		//管理员拒绝某个待审核的评论
		adminAuth.PUT("/comments/:id/reject", admin.RejectComment)
		//管理员删除某个评论
		adminAuth.DELETE("/comments/:id", admin.DeleteComment)
		//管理员获取所有用户列表
		adminAuth.GET("/users", admin.GetUsers)
		//管理员封禁某个用户
		adminAuth.PUT("/users/:id/ban", admin.BanUser)
		//管理员解封某个用户
		adminAuth.PUT("/users/:id/unban", admin.UnbanUser)
		//管理员根据用户id删除某个用户
		adminAuth.DELETE("/users/:id", admin.DeleteUser)
	}
	return r
}
