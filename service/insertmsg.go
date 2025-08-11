package service

import (
	"context"
	"gochat/global"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func InsertFriendReqMsg(dbName, senderID, content string, read int, expire int64) error {
	collection := global.MongoDBClient.Database(dbName).Collection("friend_requests")
	parts := strings.Split(senderID, "->")
	fromUserID := parts[0]
	toUserID := ""
	if len(parts) == 2 {
		toUserID = parts[1]
	}
	_, err := collection.InsertOne(context.Background(), bson.M{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"content":      content,
		"read":         read,
		"created_at":   time.Now(),
		"expire_at":    time.Now().Add(time.Duration(expire)),
	})
	if err != nil {
		global.Log.Error("插入MongoDB消息失败", zap.Error(err))
	}
	return err
}

func InsertPrivateMsg(dbName, senderID, content string, read int, expire int64) error {
	collection := global.MongoDBClient.Database(dbName).Collection("messsge")
	parts := strings.Split(senderID, "->")
	fromUserID := parts[0]
	toUserID := ""
	if len(parts) == 2 {
		toUserID = parts[1]
	}
	_, err := collection.InsertOne(context.Background(), bson.M{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"content":      content,
		"read":         read,
		"created_at":   time.Now(),
		"expire_at":    time.Now().Add(time.Duration(expire)),
	})
	if err != nil {
		global.Log.Error("插入MongoDB消息失败", zap.Error(err))
	}
	return err
}
