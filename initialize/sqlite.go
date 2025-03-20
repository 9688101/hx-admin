package initialize

import (
	"fmt"

	"github.com/9688101/hx-admin/common"
	"github.com/9688101/hx-admin/core/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSQLite() (*gorm.DB, error) {
	logger.SysLog("SQL_DSN not set, using SQLite as database")
	common.UsingSQLite = true
	dsn := fmt.Sprintf("%s?_busy_timeout=%d", common.SQLitePath, common.SQLiteBusyTimeout)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
}
