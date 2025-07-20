package main

import (
	flag "gochat/flag"
	"gochat/global"
	"gochat/initialize"
	"strconv"
)

func main() {
	global.Config = initialize.InitConfig()       // 初始化配置
	global.Log = initialize.InitLogger()          // 初始化日志记录器
	global.DB = initialize.InitGorm()             // 初始化数据库连接
	global.RedisClient = initialize.InitRedis()   // 初始化 Redis 客户端
	global.MongoDBClient = initialize.InitMongo() // 初始化 MongoDB 客户端（如果需要的话）
	flag.InitFlag()
	router := initialize.InitRouter()                                                     // 初始化路由
	router.Run(global.Config.System.Host + ":" + strconv.Itoa(global.Config.System.Port)) // 启动服务，监听端口

}
