package model

type MysqlConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	UserName string `ini:"username"`
	Password string `ini:"password"`
}

type RedisConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	Password string `ini:"password"`
	Db       int    `ini:"db"`
}
type Config struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}
