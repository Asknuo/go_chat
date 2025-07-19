package initialize

import (
	"strconv"

	"context"
	"gochat/config"
	"gochat/global"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis() *redis.Client {
	redisConfig := global.Config.Redis
	redisConf := &config.Redis{
		Host: redisConfig.Host,
		Port: redisConfig.Port,
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConf.Host + ":" + strconv.Itoa(redisConf.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		global.Log.Error("Redis连接失败", zap.Error(err))
	}
	global.Log.Info("Redis连接成功")
	return redisClient
}
