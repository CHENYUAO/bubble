package mysql

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库
var DB *sqlx.DB

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
