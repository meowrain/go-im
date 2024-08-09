package conf

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"im/model"
	"im/utils/ini_analyzer"
	"im/utils/meowlog"
)

var DbEngine *xorm.Engine

func init() {
	//日志库加载
	logger := meowlog.NewLogger("console", "fatal", "logs")
	var err error
	//Mysql配置文件加载
	conf := &Config{}
	err = ini_analyzer.LoadIni("conf/conf.ini", conf)
	if err != nil {
		logger.Error(err.Error())
	}
	driverName := "mysql"
	MySqlDbUser := conf.MysqlConfig.UserName
	MySqlDbPassword := conf.MysqlConfig.Password
	MysqlDbHost := conf.MysqlConfig.Host
	MysqlDbPort := conf.MysqlConfig.Port
	MysqlDBName := conf.MysqlConfig.DBName
	DsName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", MySqlDbUser, MySqlDbPassword, MysqlDbHost, MysqlDbPort, MysqlDBName)
	//DsName := "root:123456@tcp(127.0.0.1:3306)/chat?charset=utf8"
	DbEngine, err = xorm.NewEngine(driverName, DsName)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("数据库连接成功！")
	//是否显示SQL
	DbEngine.ShowSQL(true)
	//设置最大连接数
	DbEngine.SetMaxOpenConns(20)

	//自动创建表
	DbEngine.Sync2(new(model.User), new(model.Contact), new(model.Community))
}
