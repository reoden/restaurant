package pkgs

import (
	"fmt"
	"restaurant/common"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db *gorm.DB

func InitDB() {
	// var dsn = fmt.Sprintf("root:123456@tcp(%s:3306)/restaurant?charset=utf8mb4&parseTime=True&loc=Local", "localhost")
	var dsn = fmt.Sprintf("shiflow:Qwert_54321@tcp(%s:3306)/restaurant?charset=utf8mb4&parseTime=True&loc=Local", common.DOMAIN_HOST)
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(50) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(10)
}

func GetDB() *gorm.DB {
	return _db
}
