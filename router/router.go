package router

import "gorm.io/gorm"
import "github.com/gin-gonic/gin"

// Register 注册所有路由
func Register(r *gin.Engine, db *gorm.DB) {
	// 全局中间件
	//r.Use(middleware.CORSMiddleware())

	// ---------- 依赖注入 ----------
	// DAO
	//userDAO := dao.NewUserDAO(db)

	// Service 层
	//userService := service.NewUserService(userDAO)

	// Controller 层
	//userCtrl := controller.NewUserController(userService)

	// ---------- 路由分组 ----------
	//api := r.Group("/api/v1")
	//{
	//	// 公开接口（无需认证）
	//	public := api.Group("")
	//	{
	//		public.POST("/user/register", userCtrl.Register)
	//		public.POST("/user/login", userCtrl.Register) // TODO: 替换为 Login
	//	}
	//
	//	// 需要认证的接口
	//	protected := api.Group("")
	//	protected.Use(middleware.Auth())
	//	{
	//		protected.GET("/user/:id", userCtrl.GetUser)
	//		protected.GET("/users", userCtrl.ListUsers)
	//		protected.PUT("/user/:id", userCtrl.UpdateUser)
	//		protected.DELETE("/user/:id", userCtrl.DeleteUser)
	//	}
	//}
}
