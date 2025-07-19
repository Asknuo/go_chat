package initialize

import (
	"context"
	"fmt"
	"gochat/config"
	"gochat/global"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitMongo() *mongo.Client {
	mongoConfig := global.Config.Mongo
	mongoCfg := &config.MongoDB{
		Host: mongoConfig.Host,
		Port: mongoConfig.Port,
	}
	fmt.Println("MongoDB配置:", mongoCfg.Host, mongoCfg.Port)
	clientOptions := options.Client().ApplyURI("mongodb://" + mongoCfg.Host + ":" + strconv.Itoa(mongoCfg.Port))
	mongoDBClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		global.Log.Error("MongoDB连接失败", zap.Error(err))

	}
	err = mongoDBClient.Ping(context.TODO(), nil)
	if err != nil {
		global.Log.Error("MongoDB Ping失败", zap.Error(err))
	}
	global.Log.Info("MongoDB连接成功")
	return mongoDBClient
}
