package initialize

import (
	"fmt"

	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/utils/env"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var UsingSQLite = false
var SQLitePath = "one-api.db"
var SQLiteBusyTimeout = env.Int("SQLITE_BUSY_TIMEOUT", 3000)

func openSQLite() (*gorm.DB, error) {
	logger.SysLog("SQL_DSN not set, using SQLite as database")
	UsingSQLite = true
	dsn := fmt.Sprintf("%s?_busy_timeout=%d", SQLitePath, SQLiteBusyTimeout)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
}
