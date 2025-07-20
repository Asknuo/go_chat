package global

import (
	"gochat/config"

	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	//全局变量
	Config        *config.Config // 全局配置对象
	DB            *gorm.DB       // 数据库连接对象
	RedisClient   *redis.Client  // Redis客户端
	MongoDBClient *mongo.Client  // MongoDB客户端
	Log           *zap.Logger    // 日志记录器
)
var Store = base64Captcha.DefaultMemStore
