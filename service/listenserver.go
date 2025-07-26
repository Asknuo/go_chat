package service

import (
	"encoding/json"
	"fmt"
	"gochat/global"
	"gochat/utlis"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	fmt.Println("----------正在监听连接----------")
	for {
		select {
		case conn := <-Manager.Register: // 建立连接
			fmt.Printf("建立新连接: %v", conn.SenderID)
			Manager.Clients[conn.SenderID] = conn
			replyMsg := &Msg{
				Code:    utlis.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		case conn := <-Manager.Unregister: // 断开连接
			if _, ok := Manager.Clients[conn.SenderID]; ok {
				replyMsg := &Msg{
					Code:    utlis.WebsocketEnd,
					Content: "连接已断开",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.SenderID)
			}
		case boardcast := <-Manager.Boardcast:
			message := boardcast.Message
			recId := boardcast.Client.ReveiverID
			flag := false //默认不在线
			for id, conn := range Manager.Clients {
				if id != recId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.SenderID)
				}
			}
			id := boardcast.Client.SenderID
			if flag {
				replyMsg := &Msg{
					Code:    utlis.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = boardcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(global.Config.Mongo.Name, id, string(message), 1, int64(3*month)) //1表示在线，在线即是已读
				if err != nil {
					global.Log.Error("消息插入mongodb失败")
				}
			} else {
				fmt.Println("对方不在线")
				replyMsg := &Msg{
					Code:    utlis.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = boardcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(global.Config.Mongo.Name, id, string(message), 0, int64(3*month)) //1表示在线，在线即是已读
				if err != nil {
					global.Log.Error("消息插入mongodb失败")
				}
			}
		}
	}
}
