package gormInit

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

var DB *gorm.DB
var MySqlLogger = logger2.Default.LogMode(logger2.Info)

func init() {
	username := "root"
	password := "root"
	host := "127.0.0.1"
	port := 3306
	Dbname := "gorm"
	timeout := "10s"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&Local&timeout=%s", username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: mySqlLogger,
	})
	if err != nil {
		panic("连接数据库失败，error =" + err.Error())
	}
	DB = db
	fmt.Println(DB)
}
