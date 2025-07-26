package service

import (
	"context"
	"gochat/global"
	"time"
)

type Trainer struct {
	Content   string `bson:"content"`   // 内容
	StartTime int64  `bson:"startTime"` // 创建时间
	EndTime   int64  `bson:"endTime"`   // 过期时间
	Read      uint   `bson:"read"`      // 已读
}

func InsertMsg(database, id string, content string, isread uint, expire int64) error {
	collection := global.MongoDBClient.Database(database).Collection(id)
	comment := Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      isread,
	}
	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}
