package api

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/global"
	"gochat/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type HandleFriendReqApi struct{}

var handleupgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 根据需要配置跨域检查
	},
}

func (handle *HandleFriendReqApi) HandleFriendReq(c *gin.Context) {
	token := c.GetHeader("Sec-WebSocket-Protocol")
	if token == "" {
		c.Error(errors.New("未提供Token"))
		return
	}
	ws, err := handleupgrader.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-WebSocket-Protocol": []string{token},
	})

	if err != nil {
		global.Log.Error("WebSocket 升级失败", zap.Error(err))
		return
	}
	defer ws.Close()
	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("BMvfawCaCjlDOzLAoYDUxLZWGIzerY53VeIm03Fy6uE="), nil
	})
	if err != nil {
		c.Error(errors.New("无效token"))
		return
	}
	key := "jwt_blacklist:" + token
	exists, _ := global.RedisClient.Exists(context.Background(), key).Result()
	if exists > 0 {
		global.Log.Warn("token已存在于黑名单")
		return
	}
	claims := tokenString.Claims.(jwt.MapClaims)
	var userid string
	switch v := claims["userid"].(type) {
	case float64: // 当userid是JSON数字时
		userid = strconv.Itoa(int(v))
	case string: // 当userid是字符串时
		userid = v
	default:
		global.Log.Error("userid类型无效")
		return
	}
	for {
		// 读取客户端发送的 JSON 消息
		_, message, err := ws.ReadMessage()
		if err != nil {
			global.Log.Error("读取WebSocket消息失败", zap.Error(err))
			break
		}

		// 解析客户端发送的请求数据
		var request struct {
			Response string `json:"response"`
			Touserid string `json:"touserid"`
			Note     string `json:"note"`
		}
		if err := json.Unmarshal(message, &request); err != nil {
			global.Log.Error("消息解析失败", zap.Error(err))
			writeWSMessage(ws, gin.H{"error": "消息格式错误"})
			continue
		}
		status := request.Response
		friendid := request.Touserid
		// 验证 friendid 是否为空
		if status == "" {
			global.Log.Error("请拒接或者接受!")
			writeWSMessage(ws, gin.H{"error": "请拒接或者接受!"})
			continue
		}
		switch status {
		case "accepted":
			notifyMsg := map[string]interface{}{
				"type":         "friend_accepted",
				"from_user_id": userid,
				"to_user_id":   friendid,
				"note":         request.Note,
				"code":         200,
			}
			msgBytes, err := json.Marshal(notifyMsg)
			if err != nil {
				global.Log.Error("消息编码失败", zap.Error(err))
				writeWSMessage(ws, gin.H{"message": "好友请求接受成功，但通知编码失败"})
				continue
			}
			// 广播通知
			client := &service.Client{
				SenderID:   service.CreateID(userid, friendid),
				ReveiverID: service.CreateID(friendid, userid),
				Socket:     ws,
				Send:       make(chan []byte),
			}
			service.Manager.Boardcast <- &service.Boardcast{
				Client:  client,
				Message: msgBytes,
				Type:    "friend_accepted",
			}
			err = service.HandleFriendRequest("accepted", userid, friendid)
			if err != nil {
				global.Log.Error("接受好友操作异常", zap.Error(err))
				continue
			}
		case "rejected":
			notifyMsg := map[string]interface{}{
				"type":         "friend_rejected",
				"from_user_id": userid,
				"to_user_id":   friendid,
				"note":         request.Note,
				"code":         200,
			}
			msgBytes, err := json.Marshal(notifyMsg)
			if err != nil {
				global.Log.Error("消息编码失败", zap.Error(err))
				writeWSMessage(ws, gin.H{"message": "好友请求拒绝成功,但通知编码失败"})
				continue
			}
			// 广播通知
			client := &service.Client{
				SenderID:   service.CreateID(userid, friendid),
				ReveiverID: service.CreateID(friendid, userid),
				Socket:     ws,
				Send:       make(chan []byte),
			}
			service.Manager.Boardcast <- &service.Boardcast{
				Client:  client,
				Message: msgBytes,
				Type:    "friend_rejected",
			}
			err = service.HandleFriendRequest("rejected", userid, friendid)
			if err != nil {
				global.Log.Error("拒绝好友操作异常", zap.Error(err))
				continue
			}
		}
	}
}
