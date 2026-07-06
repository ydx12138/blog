package router

import (
	"blog/internal/handler"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(h *handler.Handler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		public := api.Group("")
		public.GET("/articles", h.User.GetArticles)
		public.GET("/articles/detail", h.User.GetArticle)
		public.GET("/articles/search", h.User.SearchArticle)
		public.GET("/categories", h.User.GetCategories)
		public.GET("/categories/articles", h.User.GetCategoryArticles)
		public.GET("/comments", h.User.GetComments)
		public.POST("/register/code", h.User.SendRegisterCode)
		public.POST("/register", h.User.Register)
		public.POST("/login", h.User.Login)
		public.POST("/articles/like", h.User.LikeArticle)
		public.GET("/tags", h.User.GetTags)
		public.POST("/sendpwdcode", h.User.SendCodeForgetPwd)
		public.POST("/updatePasswordByCode", h.User.UpdatePasswordByCode)
	}
	{
		apiAuth := api.Group("")
		apiAuth.Use(middleware.JWTAuth())
		apiAuth.POST("/comments", h.User.CreateComment)
		apiAuth.POST("/updatephonenumber", h.User.UpdatePhoneNumber)
	}

	adminGroup := r.Group("/api/admin")
	adminGroup.POST("/login", h.Admin.Login)

	adminAuth := adminGroup.Group("")
	adminAuth.Use(middleware.JWTAuthForAdmin())
	adminAuth.GET("/dashboard", h.Admin.GetDashboard)
	adminAuth.GET("/articles", h.Admin.GetArticles)
	adminAuth.GET("/articles/:id", h.Admin.GetArticle)
	adminAuth.POST("/articles", h.Admin.CreateArticle)
	adminAuth.PUT("/articles/:id", h.Admin.UpdateArticle)
	adminAuth.DELETE("/articles/:id", h.Admin.DeleteArticle)
	adminAuth.GET("/drafts", h.Admin.GetDrafts)
	adminAuth.PUT("/articles/:id/publish", h.Admin.PublishArticle)
	adminAuth.POST("/upload", h.Admin.UploadImage)
	adminAuth.GET("/comments", h.Admin.GetAllComments)
	adminAuth.GET("/comments/pending", h.Admin.GetPendingComments)
	adminAuth.PUT("/comments/:id/approve", h.Admin.ApproveComment)
	adminAuth.PUT("/comments/:id/reject", h.Admin.RejectComment)
	adminAuth.DELETE("/comments/:id", h.Admin.DeleteComment)
	adminAuth.PUT("/comments/:id/status", h.Admin.SetCommentStatus)
	adminAuth.GET("/users", h.Admin.GetUsers)
	adminAuth.PUT("/users/:id/ban", h.Admin.BanUser)
	adminAuth.PUT("/users/:id/unban", h.Admin.UnbanUser)
	adminAuth.DELETE("/users/:id", h.Admin.DeleteUser)

	return r
}
