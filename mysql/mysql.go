package mysql

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/ini.v1"
)

// 数据库
var DB *sqlx.DB

type MysqlConfig struct {
	User   string `ini:"user"`
	Passwd string `ini:"passwd"`
	Host   string `ini:"host"`
	Port   string `ini:"port"`
	DB     string `ini:"database"`
}

// 从配置文件中读取数据库配置信息，返回dsn
func ReadConf(path string) (string, error) {
	//从配置文件中读取数据库配置信息
	cfg, err := ini.Load(path)
	if err != nil {
		return "", err
	}
	MysqlCfg := new(MysqlConfig)
	err = cfg.Section("mysql").MapTo(MysqlCfg)
	dsn := fmt.Sprint(MysqlCfg.User, ":", MysqlCfg.Passwd, "@tcp", "(", MysqlCfg.Host, ":", MysqlCfg.Port, ")", "/", MysqlCfg.DB)
	log.Println("parse config files success")
	return dsn, nil
}

func InitDB(dsn string) (err error) {
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)
	log.Println("init database success")
	return nil
}
