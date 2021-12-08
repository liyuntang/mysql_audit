package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func GetEngine(dataSource string) (engine *xorm.Engine, e error) {
	// 连接数据库
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

