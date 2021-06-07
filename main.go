package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Todo struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// 数据库
var db *sqlx.DB

//gin的引擎
var r *gin.Engine

// 从配置文件中读取数据库配置信息，返回dsn
func readConf(path string) (string, error) {
	//从配置文件中读取数据库配置信息
	cfg, err := ini.Load(path)
	if err != nil {
		return "", err
	}
	user := cfg.Section("mysql").Key("user").String()
	passwd := cfg.Section("mysql").Key("passwd").String()
	host := cfg.Section("mysql").Key("host").String()
	port := cfg.Section("mysql").Key("port").String()
	database := cfg.Section("mysql").Key("database").String()
	dsn := fmt.Sprint(user, ":", passwd, "@tcp", "(", host, ":", port, ")", "/", database)
	log.Println("parse config files success")
	return dsn, nil
}

func initDB(dsn string) (err error) {
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	log.Println("init database success")
	return nil
}

func startEngine() {
	r = gin.Default()
	//加载静态文件
	r.Static("/static", "static")
	//加载模板文件夹
	r.LoadHTMLGlob("template/*")
	//注册路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// v1路由组
	//待办事项
	v1Group := r.Group("v1")
	{
		//添加
		v1Group.POST("/todo", postTitle)
		//查看所有待办事项
		v1Group.GET("/todo", getTitle)
		//修改某个待办事项
		v1Group.PUT("/todo/:id", putTitle)
		//删除
		v1Group.DELETE("/todo/:id", deleteTitle)
	}
	r.Run(":40000")
}

/**********************************HTTP处理函数****************************/
func postTitle(c *gin.Context) {
	//读取前端发送过来的数据
	var t Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		log.Println("get json failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	//将数据写入数据库
	sqlStr := `insert into task(title,status) values(?,?);`
	if _, err := db.Exec(sqlStr, t.Title, 0); err != nil {
		log.Println("insert into db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	//log.Println("insert into db success")
	c.JSON(http.StatusOK, t)
}

func getTitle(c *gin.Context) {
	var todoList []Todo
	sqlStr := `select * from task;`
	if err := db.Select(&todoList, sqlStr); err != nil {
		log.Println("query db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, todoList)
}

func putTitle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("parse url param failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 先select，再取反再返回
	var t Todo
	sqlStrQuery := `select * from task where id=?`
	if err := db.Get(&t, sqlStrQuery, id); err != nil {
		log.Println("update db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	t.Status = !t.Status
	sqlStrUpdate := `update task set status=? where id = ?;`
	if _, err := db.Exec(sqlStrUpdate, t.Status, id); err != nil {
		log.Println("update db failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("update status,id:", id)
	c.JSON(http.StatusOK, t)
}

func deleteTitle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("parse url param failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	sqlStr := `delete from task where id = ?;`
	if _, err := db.Exec(sqlStr, id); err != nil {
		log.Println("delete title failed,err:", err)
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

/***********************************main函数**************************/
func main() {
	// 读取配置信息
	dsn, err := readConf("ini/conf.ini")
	if err != nil {
		log.Fatal("read config files failed,err:", err)
	}
	//连接数据库
	if err = initDB(dsn); err != nil {
		log.Panic("connect to mysql failed,err:", err)
	}
	defer db.Close()
	//开启web服务
	startEngine()
}
