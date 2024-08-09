package service

import (
	"github.com/go-xorm/xorm"
	"im/model"
	"log"
)

var DbEngine *xorm.Engine

func init() {
	var err error
	driverName := "mysql"
	DsName := "root:123456@tcp(127.0.0.1:3306)/chat?charset=utf8"
	DbEngine, err = xorm.NewEngine(driverName, DsName)
	if err != nil {
		log.Fatal(err)
	}
	//是否显示SQL
	DbEngine.ShowSQL(true)
	//设置最大连接数
	DbEngine.SetMaxOpenConns(20)

	//自动User
	DbEngine.Sync2(new(model.User))
}
