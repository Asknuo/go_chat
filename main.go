package main

import (
	flag "gochat/flag"
	"gochat/global"
	"gochat/initialize"
)

func main() {
	global.Config = initialize.InitConfig() // 初始化配置
	global.Log = initialize.InitLogger()
	global.DB = initialize.InitGorm() // 初始化数据库连接

	flag.InitFlag()
}
