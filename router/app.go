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
	r.POST("/problem-create", middlewares.AuthAdminCheck(), service.ProblemCreate)
	return r
}
