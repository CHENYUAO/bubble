package handle

import (
	"bubble/mysql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

//gin的引擎
var r *gin.Engine

func StartEngine() {
	r = gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("template/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	loginGroup := r.Group("login")
	{
		loginGroup.POST("/", Authority)
		//这里是一个bug，点击进入到清单界面之后，网页url是xxx/login/#/，点击刷新按钮会GET /login 所以临时这样处理，后续想办法解决问题
		loginGroup.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})
	}
	signupGroup := r.Group("/signup")
	{
		signupGroup.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "signup.html", nil)
		})
		signupGroup.POST("/", SignUp)
	}
	//待办事项
	v1Group := r.Group("v1")
	{
		//添加
		v1Group.POST("/todo", PostTitle)
		//查看所有待办事项
		v1Group.GET("/todo", GetTitle)
		//修改某个待办事项
		v1Group.PUT("/todo/:id", PutTitle)
		//删除
		v1Group.DELETE("/todo/:id", DeleteTitle)
	}
	r.Run(":40000")
}

func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	again := c.PostForm("again")
	if password != again {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入密码不一致",
		})
	}
	if err := mysql.InsertUser(username, password); err != nil {
		log.Printf("insert failed,err:%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"error": "更新数据库失败",
		})
		return
	}
	//跳转到登录页面
	c.HTML(http.StatusOK, "login.html", nil)
	return
}

func Authority(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	//验证用户名和密码
	err := mysql.AuthUser(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", nil)
}

func PostTitle(c *gin.Context) {
	tx, err := mysql.DB.Beginx()
	if err != nil {
		tx.Rollback()
		log.Println("begin trans failed,err:", err)
		return
	}
	//读取前端发送过来的数据
	var t Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		tx.Rollback()
		log.Println("get json failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	//将数据写入数据库
	sqlStr := `insert into task(title,status) values(?,?);`
	if _, err := mysql.DB.Exec(sqlStr, t.Title, 0); err != nil {
		tx.Rollback()
		log.Println("insert into db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, t)
	tx.Commit()
}

func GetTitle(c *gin.Context) {
	tx, err := mysql.DB.Beginx()
	if err != nil {
		tx.Rollback()
		log.Println("begin trans failed,err:", err)
		return
	}
	var todoList []Todo
	sqlStr := `select * from task;`
	if err := mysql.DB.Select(&todoList, sqlStr); err != nil {
		tx.Rollback()
		log.Println("query db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, todoList)
	tx.Commit()
}

func PutTitle(c *gin.Context) {
	tx, err := mysql.DB.Beginx()
	if err != nil {
		tx.Rollback()
		log.Println("begin trans failed,err:", err)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tx.Rollback()
		log.Println("parse url param failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	var t Todo
	sqlStrQuery := `select * from task where id=?`
	if err := tx.Get(&t, sqlStrQuery, id); err != nil {
		tx.Rollback()
		log.Println("update db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	t.Status = !t.Status
	sqlStrUpdate := `update task set status=? where id = ?;`
	if _, err := tx.Exec(sqlStrUpdate, t.Status, id); err != nil {
		tx.Rollback()
		log.Println("update db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("update status,id:", id)
	c.JSON(http.StatusOK, t)
	tx.Commit()
}

func DeleteTitle(c *gin.Context) {
	tx, err := mysql.DB.Beginx()
	if err != nil {
		tx.Rollback()
		log.Println("begin trans failed,err:", err)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tx.Rollback()
		log.Println("parse url param failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	sqlStr := `delete from task where id = ?;`
	if _, err := mysql.DB.Exec(sqlStr, id); err != nil {
		tx.Rollback()
		log.Println("delete title failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	tx.Commit()
}
