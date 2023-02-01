package router

import (
	_ "gin_oj/docs"
	"gin_oj/middlewares"
	"gin_oj/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	//swag配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//problem
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	//user
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/user-login", service.Login)
	r.POST("/user-send_code", service.SendCode)
	r.POST("/user-register", service.Register)
	//排行榜
	r.GET("/rank-list", service.GetRankList)
	//submit
	r.GET("/submit-list", service.GetSubmitList)

	// 管理员私有方法
	authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
	// 问题创建
	authAdmin.POST("/problem-create", service.ProblemCreate)
	// 问题修改
	authAdmin.PUT("/problem-update", service.ProblemUpdate)
	// 分类列表
	authAdmin.GET("/category-list", service.GetCategoryList)
	authAdmin.POST("/category-create", service.CategoryCreate)
	authAdmin.DELETE("/category-delete", service.CategoryDelete)
	authAdmin.PUT("/category-update", service.CategoryUpdate)

	// 用户私有方法
	authUser := r.Group("/me", middlewares.AuthUserCheck())
	// 代码提交
	authUser.POST("/submit", service.Submit)

	return r
}
