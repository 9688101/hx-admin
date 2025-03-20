package initialize

import (
	"github.com/9688101/hx-admin/core/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var UsingMySQL = false

func openMySQL(dsn string) (*gorm.DB, error) {
	logger.SysLog("using MySQL as database")
	UsingMySQL = true
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
}
