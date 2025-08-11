package api

import (
	"context"
	"fmt"
	"gochat/global"
	"gochat/request"
	"gochat/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type GethistorymsgApi struct{}

func (history *GethistorymsgApi) GetHistoryMsg(c *gin.Context) {
	var gethistorymsg request.GetHistoryMsgReq
	if err := c.ShouldBindJSON(&gethistorymsg); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 解析 Token
	tokenString := extractToken(c)
	if tokenString == "" {
		global.Log.Error("未解析到token")
		response.FailWithMessage("未解析到token", c)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名算法: %v", token.Header["alg"])
		}
		return []byte("BMvfawCaCjlDOzLAoYDUxLZWGIzerY53VeIm03Fy6uE="), nil
	})
	if err != nil {
		global.Log.Error("验证token错误")
		response.FailWithMessage("验证token错误", c)
		return
	}

	// 检查黑名单
	if isTokenBlacklisted(tokenString) {
		global.Log.Error("token在黑名单中")
		response.FailWithMessage("token在黑名单中", c)
		return
	}

	// 解析 userid
	claims := token.Claims.(jwt.MapClaims)
	var userid string
	switch v := claims["userid"].(type) {
	case float64:
		userid = strconv.Itoa(int(v))
	case string:
		userid = v
	default:
		global.Log.Error("userid类型无效")
		response.FailWithMessage("userid类型无效", c)
		return
	}

	collection := global.MongoDBClient.Database(global.Config.Mongo.Name).Collection("messsge")

	// 1. 批量更新对方发给我的未读消息为已读
	updateFilter := bson.M{
		"from_user_id": gethistorymsg.Touserid,
		"to_user_id":   userid,
		"read":         0,
	}
	_, err = collection.UpdateMany(
		context.Background(),
		updateFilter,
		bson.M{"$set": bson.M{"read": 1, "read_at": time.Now()}},
	)
	if err != nil {
		global.Log.Error("批量更新消息状态失败", zap.Error(err))
	}

	// 2. 查询双方的所有历史消息
	filter := bson.M{
		"$or": []bson.M{
			{"from_user_id": userid, "to_user_id": gethistorymsg.Touserid},
			{"from_user_id": gethistorymsg.Touserid, "to_user_id": userid},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}) // 按时间升序
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		global.Log.Error("查询历史消息失败", zap.Error(err))
		response.FailWithMessage("查询历史消息失败", c)
		return
	}
	defer cursor.Close(context.Background())

	// 3. 读取结果
	var messages []bson.M
	if err := cursor.All(context.Background(), &messages); err != nil {
		global.Log.Error("解码消息失败", zap.Error(err))
		response.FailWithMessage("解码消息失败", c)
		return
	}

	// 4. 返回结果
	response.OkWithData(messages, c)
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func isTokenBlacklisted(tokenString string) bool {
	key := "jwt_blacklist:" + tokenString
	exists, err := global.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		global.Log.Error("查询Redis失败", zap.Error(err))
		return false
	}
	return exists > 0
}
