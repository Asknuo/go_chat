package service

import (
	"context"
	"fmt"
	"gochat/global"
	"gochat/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type FriendRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FromUserID string             `bson:"from_user_id"`
	ToUserID   string             `bson:"to_user_id"`
	Note       string             `bson:"note"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func HandleFriendRequest(req string, FromUserID string, ToUserID string) error {
	collection := global.MongoDBClient.Database(global.Config.Mongo.Name).Collection("friend_requests")
	switch req {
	case "friend_request":
		count, err := collection.CountDocuments(context.Background(), bson.M{
			"from_user_id": FromUserID,
			"to_user_id":   ToUserID,
			"status":       "pending",
		})
		if err != nil {
			global.Log.Error("检查好友请求失败", zap.Error(err))
			return err
		}
		if count > 0 {
			return fmt.Errorf("好友请求已存在")
		}
	case "accepted":
		result := collection.FindOne(
			context.Background(),
			bson.M{
				"from_user_id": FromUserID,
				"to_user_id":   ToUserID,
			},
		)
		var friendRequest FriendRequest
		if err := result.Decode(&friendRequest); err != nil {
			global.Log.Error("未找到匹配的好友请求或已处理", zap.Error(err))
			return fmt.Errorf("好友请求不存在或已处理")
		}
		fromUserID, _ := parseUserID(FromUserID)
		touserid, _ := parseUserID(ToUserID)
		err := global.DB.Model(&models.Friendship{}).Where("user_id = ? AND friend_id = ?", uint(touserid), uint(fromUserID)).Update("status", "accepted").Error
		if err != nil {
			global.Log.Error("好友状态修改失败", zap.Error(err))
			return err
		}
	case "rejected":
		// 查找并更新待处理的好友请求
		result := collection.FindOne(
			context.Background(),
			bson.M{
				"from_user_id": FromUserID,
				"to_user_id":   ToUserID,
			},
		)
		var friendRequest FriendRequest
		if err := result.Decode(&friendRequest); err != nil {
			global.Log.Error("未找到匹配的好友请求或已处理", zap.Error(err))
			return err
		}
		fromUserID, _ := parseUserID(FromUserID)
		touserid, _ := parseUserID(ToUserID)
		err := global.DB.Model(&models.Friendship{}).Where("user_id = ? AND friend_id = ?", uint(touserid), uint(fromUserID)).Update("status", "rejected").Error
		if err != nil {
			global.Log.Error("好友状态修改失败", zap.Error(err))
			return err
		}
	}
	return nil
}

func parseUserID(id string) (int, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("无效的用户ID: %s", id)
	}
	return userID, nil
}
