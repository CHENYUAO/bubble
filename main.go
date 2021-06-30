package main

import (
	"log"

	"github.com/chenyuao/bubble/database/mysql"
	"github.com/chenyuao/bubble/handle"
)

func main() {
	// 读取配置信息
	dsn, err := mysql.ReadConf("./conf/bubble.ini")
	if err != nil {
		log.Fatal("read config files failed,err:", err)
	}
	//连接数据库
	if err = mysql.InitDB(dsn); err != nil {
		log.Panic("connect to mysql failed,err:", err)
	}
	defer mysql.DB.Close()
	//开启web服务
	handle.StartEngine()
}
