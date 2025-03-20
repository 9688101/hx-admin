package initialize

import (
	"github.com/9688101/hx-admin/common"
	"github.com/9688101/hx-admin/core/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openPostgreSQL(dsn string) (*gorm.DB, error) {
	logger.SysLog("using PostgreSQL as database")
	common.UsingPostgreSQL = true
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
}
