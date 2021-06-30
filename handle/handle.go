package handle

import (
	"net/http"

	"github.com/chenyuao/bubble/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//gin的引擎
var Router *gin.Engine

func StartEngine() {
	Router = gin.Default()
	RouterRegister()
	Router.Run(":40000")
}

func RouterRegister() {
	Router.Static("/static", "static")
	Router.LoadHTMLGlob("template/*")
	Router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })

	//login路由组
	loginGroup := Router.Group("login")
	{
		loginGroup.POST("/", service.Authority)
		/* 这里是一个bug，点击进入到清单界面之后，网页url是xxx/login/#/
		点击刷新按钮会GET /login 所以临时这样处理，后续想办法解决问题*/
		loginGroup.GET("/", service.GetLogin)
	}
	//signup路由组
	signupGroup := Router.Group("/signup")
	{
		signupGroup.GET("/", service.GetSignup)
		signupGroup.POST("/", service.SignUp)
	}
	//v1路由组
	v1Group := Router.Group("v1")
	{
		//添加待办事项
		v1Group.POST("/todo", service.PostTitle)
		//查看所有待办事项
		v1Group.GET("/todo", service.GetTitle)
		//修改某个待办事项
		v1Group.PUT("/todo/:id", service.PutTitle)
		//删除待办事项
		v1Group.DELETE("/todo/:id", service.DeleteTitle)
	}
}
