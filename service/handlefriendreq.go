package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gochat/global"
	"gochat/request"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type FriendRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FromUserID string             `bson:"from_user_id"`
	ToUserID   string             `bson:"to_user_id"`
	Status     string             `bson:"status"` // pending, accepted, rejected
	Note       string             `bson:"note"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func HandleFriendRequest(client *Client, msg *Msg) error {
	var req request.FriendRequestMsg
	if err := json.Unmarshal([]byte(msg.Content), &req); err != nil {
		global.Log.Error("解析好友请求消息失败", zap.Error(err))
		return err
	}
	collection := global.MongoDBClient.Database(global.Config.Mongo.Name).Collection("friend_requests")
	switch req.Type {
	case "friend_request":
		count, err := collection.CountDocuments(context.Background(), bson.M{
			"from_user_id": req.FromUserID,
			"to_user_id":   req.ToUserID,
			"status":       "pending",
		})
		if err != nil {
			global.Log.Error("检查好友请求失败", zap.Error(err))
			return err
		}
		if count > 0 {
			return fmt.Errorf("好友请求已存在")
		}

	}
	return nil
}
