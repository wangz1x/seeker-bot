// mariadb - 2024/12/16
// Author: wangzx
// Description:

package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root@tcp(127.0.0.1:3306)/seekerbot?charset=utf8mb4&parseTime=True&loc=Local"

var DB *gorm.DB

func init() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
