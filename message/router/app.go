package router

import (
	"genchat/message/docs"
	"genchat/message/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	// swagger
	docs.SwaggerInfo.BasePath = ""
	// 首页相关
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	// r.Static("/asset", "asset/")
	// r.LoadHTMLGlob("views/**/*")
	//用户模块
	r.GET("/index", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList)
	r.GET("/user/findUserByName", service.FindUserByName)
	r.POST("/user/login", service.Login)
	r.POST("/user/register", service.Register)
	r.GET("/user/delete", service.DeleteUser) //删除指定用户名的用户
	r.POST("/user/update", service.UpdateUser)

	//发送消息
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.GET("/user/sendMsg", service.SendMsg)
	return r
}
