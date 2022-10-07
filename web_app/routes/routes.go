package routes

import (
	"web_app/controller"
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用jwt认证中间件

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.CetPostDetailHandler)
		v1.GET("/posts", controller.CetPostListlHandler)

		//根据时间或分数获取帖子列表
		v1.GET("/posts2", controller.CetPostListlHandler2)

		//投票
		v1.POST("/vote", controller.PostVoteController)
	}

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是登录的用户,判断请求头中是否也有有效的jwt ？
		c.String(200, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "404找不到页面",
		})
	})
	return r
}
