package mysql

import (
	"io"

	_ "github.com/go-sql-driver/mysql" // just to register mysql driver
	"github.com/xormplus/xorm"
)

func GetMysqlXORM(dataSourceName string, logger io.Writer, showSQL bool, showExecTime bool) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewMySQL("mysql", dataSourceName)
	if err == nil {
		engine.ShowSQL(showSQL)
		engine.ShowExecTime(showExecTime)

		if logger != nil {
			engine.SetLogger(xorm.NewSimpleLogger(logger))
		}
	}
}


