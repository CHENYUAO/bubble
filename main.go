package main

import (
	"fmt"
	"log"

	"bubble/handle"
	"bubble/mysql"

	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
)

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

/***********************************main函数**************************/
func main() {
	// 读取配置信息
	dsn, err := readConf("ini/conf.ini")
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
