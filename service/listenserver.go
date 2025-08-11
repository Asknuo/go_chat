package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gochat/global"
	"gochat/models"
	"gochat/utlis"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (manager *ClientManager) Start() {
	fmt.Println("----------正在监听连接----------")
	for {
		select {
		case conn := <-manager.Register: // 建立连接
			fmt.Printf("建立新连接: %v", conn.SenderID)
			Manager.Clients[conn.SenderID] = conn
			// 发送连接成功消息
			replyMsg := &Msg{
				Code:    utlis.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			msg, err := json.Marshal(replyMsg)
			if err != nil {
				global.Log.Error("消息编码失败", zap.Error(err))
				continue
			}
			if err := conn.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
				global.Log.Error("发送连接成功消息失败", zap.Error(err))
			}
			userID := strings.Split(conn.SenderID, "->")[0]
			collection1 := global.MongoDBClient.Database(global.Config.Mongo.Name).Collection("messsge")
			cursor, err := collection1.Find(context.Background(), bson.M{"to_user_id": userID, "read": 0})
			if err != nil {
				global.Log.Error("查询未读消息失败", zap.Error(err))
				continue
			}
			for cursor.Next(context.Background()) {
				var mongoMsg struct {
					Content string `bson:"content"`
				}
				if err := cursor.Decode(&mongoMsg); err != nil {
					global.Log.Error("解码消息失败", zap.Error(err))
					continue
				}
				select {
				case conn.Send <- []byte(mongoMsg.Content):
					// 更新消息为已读
					_, err := collection1.UpdateOne(
						context.Background(),
						bson.M{"to_user_id": userID, "content": mongoMsg.Content, "read": 0},
						bson.M{"$set": bson.M{"read": 1, "read_at": time.Now()}},
					)
					if err != nil {
						global.Log.Error("更新消息状态失败", zap.Error(err))
					} else {
						global.Log.Info("推送未读消息", zap.String("user_id", userID))
					}
				default:
					global.Log.Warn("推送未读消息失败，通道阻塞", zap.String("user_id", userID))
				}
			}
			if err := cursor.Close(context.Background()); err != nil {
				global.Log.Error("关闭MongoDB cursor失败", zap.Error(err))
			}
		case conn := <-manager.Unregister: // 断开连接
			if _, ok := Manager.Clients[conn.SenderID]; ok {
				replyMsg := &Msg{
					Code:    utlis.WebsocketEnd,
					Content: "连接已断开",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.SenderID)
				global.Log.Info("客户端注销", zap.String("sender_id", conn.SenderID))
			}
		case broadcast := <-manager.Boardcast:
			switch broadcast.Type {
			case "friend_request":
				message, _ := json.Marshal(broadcast.Message)
				targetUserID := ""
				if broadcast.Client.ReveiverID != "" {
					parts := strings.Split(broadcast.Client.ReveiverID, "->")
					if len(parts) == 2 {
						targetUserID = parts[0] // touid
					}
				}
				if targetUserID == "" {
					global.Log.Warn("无效的ReceiverID", zap.String("receiver_id", broadcast.Client.ReveiverID))
					continue
				}
				flag := false // 是否找到在线客户端
				for senderID, conn := range manager.Clients {
					// 匹配目标用户的客户端
					if senderID == targetUserID || strings.HasPrefix(senderID, targetUserID+"->") {
						select {
						case conn.Send <- message:
							flag = true
							global.Log.Info("消息已发送", zap.String("to", senderID))
						default:
							global.Log.Warn("客户端通道阻塞", zap.String("to", senderID))
							close(conn.Send)
							delete(manager.Clients, senderID)
						}
					}
				}
				// 通知发送者
				replyMsg := &Msg{
					Code:    utlis.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				if !flag {
					replyMsg.Code = utlis.WebsocketOfflineReply
					replyMsg.Content = "对方不在线应答"
				}
				msg, err := json.Marshal(replyMsg)
				if err != nil {
					global.Log.Error("消息编码失败", zap.Error(err))
					continue
				}
				if err := broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
					global.Log.Error("发送应答消息失败", zap.Error(err))
				}
				// 存储消息到 MongoDB
				readStatus := 0
				if flag {
					readStatus = 1 // 在线即标记为已读
				}
				err = InsertFriendReqMsg(global.Config.Mongo.Name, broadcast.Client.SenderID, string(broadcast.Message), readStatus, int64(3*month))
				if err != nil {
					global.Log.Error("消息插入MongoDB失败", zap.Error(err))
				}
			case "friend_accepted":
				message, _ := json.Marshal(broadcast.Message)
				targetUserID := ""
				if broadcast.Client.ReveiverID != "" {
					parts := strings.Split(broadcast.Client.ReveiverID, "->")
					if len(parts) == 2 {
						targetUserID = parts[0] // touid
					}
				}
				if targetUserID == "" {
					global.Log.Warn("无效的ReceiverID", zap.String("receiver_id", broadcast.Client.ReveiverID))
					continue
				}
				flag := false // 是否找到在线客户端
				for senderID, conn := range manager.Clients {
					// 匹配目标用户的客户端
					if senderID == targetUserID || strings.HasPrefix(senderID, targetUserID+"->") {
						select {
						case conn.Send <- message:
							flag = true
							global.Log.Info("消息已发送", zap.String("to", senderID))
						default:
							global.Log.Warn("客户端通道阻塞", zap.String("to", senderID))
							close(conn.Send)
							delete(manager.Clients, senderID)
						}
					}
				}
				// 通知发送者
				replyMsg := &Msg{
					Code:    utlis.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				if !flag {
					replyMsg.Code = utlis.WebsocketOfflineReply
					replyMsg.Content = "对方不在线应答"
				}
				msg, err := json.Marshal(replyMsg)
				if err != nil {
					global.Log.Error("消息编码失败", zap.Error(err))
					continue
				}
				if err := broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
					global.Log.Error("发送应答消息失败", zap.Error(err))
				}
				// 存储消息到 MongoDB
				readStatus := 0
				if flag {
					readStatus = 1 // 在线即标记为已读
				}
				err = InsertFriendReqMsg(global.Config.Mongo.Name, broadcast.Client.SenderID, string(broadcast.Message), readStatus, int64(3*month))
				if err != nil {
					global.Log.Error("消息插入MongoDB失败", zap.Error(err))
				}
			case "friend_rejected":
				message, _ := json.Marshal(broadcast.Message)
				targetUserID := ""
				if broadcast.Client.ReveiverID != "" {
					parts := strings.Split(broadcast.Client.ReveiverID, "->")
					if len(parts) == 2 {
						targetUserID = parts[0] // touid
					}
				}
				if targetUserID == "" {
					global.Log.Warn("无效的ReceiverID", zap.String("receiver_id", broadcast.Client.ReveiverID))
					continue
				}
				flag := false // 是否找到在线客户端
				for senderID, conn := range manager.Clients {
					// 匹配目标用户的客户端
					if senderID == targetUserID || strings.HasPrefix(senderID, targetUserID+"->") {
						select {
						case conn.Send <- message:
							flag = true
							global.Log.Info("消息已发送", zap.String("to", senderID))
						default:
							global.Log.Warn("客户端通道阻塞", zap.String("to", senderID))
							close(conn.Send)
							delete(manager.Clients, senderID)
						}
					}
				}
				// 通知发送者
				replyMsg := &Msg{
					Code:    utlis.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				if !flag {
					replyMsg.Code = utlis.WebsocketOfflineReply
					replyMsg.Content = "对方不在线应答"
				}
				msg, err := json.Marshal(replyMsg)
				if err != nil {
					global.Log.Error("消息编码失败", zap.Error(err))
					continue
				}
				if err := broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
					global.Log.Error("发送应答消息失败", zap.Error(err))
				}
				// 存储消息到 MongoDB
				readStatus := 0
				if flag {
					readStatus = 1 // 在线即标记为已读
				}
				err = InsertFriendReqMsg(global.Config.Mongo.Name, broadcast.Client.SenderID, string(broadcast.Message), readStatus, int64(3*month))
				if err != nil {
					global.Log.Error("消息插入MongoDB失败", zap.Error(err))
				}
			case "private":
				message, _ := json.Marshal(broadcast.Message)
				targetUserID := ""
				if broadcast.Client.ReveiverID != "" {
					parts := strings.Split(broadcast.Client.ReveiverID, "->")
					if len(parts) == 2 {
						targetUserID = parts[0] // touid
					}
				}
				if targetUserID == "" {
					global.Log.Warn("无效的ReceiverID", zap.String("receiver_id", broadcast.Client.ReveiverID))
					continue
				}
				var friendship models.Friendship // 修正变量名
				if err := global.DB.Model(&models.Friendship{}).
					Where("user_id = ? OR friend_id = ?", targetUserID, targetUserID).
					First(&friendship).Error; err != nil {
					// 处理查询错误或非好友关系情况
				}
				if friendship.Status == "accepted" {
					flag := false // 是否找到在线客户端
					for senderID, conn := range manager.Clients {
						// 匹配目标用户的客户端
						if senderID == targetUserID || strings.HasPrefix(senderID, targetUserID+"->") {
							select {
							case conn.Send <- message:
								flag = true
								global.Log.Info("消息已发送", zap.String("to", senderID))
							default:
								global.Log.Warn("客户端通道阻塞", zap.String("to", senderID))
								close(conn.Send)
								delete(manager.Clients, senderID)
							}
						}
					}
					// 通知发送者
					replyMsg := &Msg{
						Code:    utlis.WebsocketOnlineReply,
						Content: "对方在线应答",
					}
					if !flag {
						replyMsg.Code = utlis.WebsocketOfflineReply
						replyMsg.Content = "对方不在线应答"
					}
					msg, err := json.Marshal(replyMsg)
					if err != nil {
						global.Log.Error("消息编码失败", zap.Error(err))
						continue
					}
					if err := broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
						global.Log.Error("发送应答消息失败", zap.Error(err))
					}
					// 存储消息到 MongoDB
					readStatus := 0
					if flag {
						readStatus = 1 // 在线即标记为已读
					}
					err = InsertPrivateMsg(global.Config.Mongo.Name, broadcast.Client.SenderID, string(broadcast.Message), readStatus, int64(3*month))
					if err != nil {
						global.Log.Error("消息插入MongoDB失败", zap.Error(err))
					}
				} else {
					global.Log.Error("ta目前不是你的好友,无法发送")
					continue
				}
			}
		}
	}
}
