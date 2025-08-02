package api

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/global"
	"gochat/models"
	"gochat/service"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type AddFriendApi struct{}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 根据需要配置跨域检查
	},
}

func (addfriend *AddFriendApi) Addfriend(c *gin.Context) {
	token := c.GetHeader("Sec-WebSocket-Protocol")
	if token == "" {
		c.Error(errors.New("未提供Token"))
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{
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
			FriendID string `json:"friendid"`
			Note     string `json:"note"`
		}
		if err := json.Unmarshal(message, &request); err != nil {
			global.Log.Error("消息解析失败", zap.Error(err))
			writeWSMessage(ws, gin.H{"error": "消息格式错误"})
			continue
		}

		friendid := request.FriendID
		note := request.Note

		// 验证 friendid 是否为空
		if friendid == "" {
			global.Log.Error("请选择添加的好友!")
			writeWSMessage(ws, gin.H{"error": "请选择添加的好友!"})
			continue
		}

		// 检查是否为自己
		if friendid == userid {
			global.Log.Error("不能添加自己为好友!")
			writeWSMessage(ws, gin.H{"error": "不能添加自己为好友!"})
			continue
		}

		// 转换为整数
		user_id, err := strconv.Atoi(userid)
		if err != nil {
			global.Log.Error("用户ID转换失败", zap.Error(err))
			writeWSMessage(ws, gin.H{"error": "用户ID无效"})
			continue
		}
		friend_id, err := strconv.Atoi(friendid)
		if err != nil {
			global.Log.Error("目标用户ID转换失败", zap.Error(err))
			writeWSMessage(ws, gin.H{"error": "目标用户ID无效"})
			continue
		}

		// 检查目标用户是否存在
		var friend models.User
		if err := global.DB.Where("id = ?", uint(friend_id)).First(&friend).Error; err != nil {
			global.Log.Error("目标用户不存在", zap.Error(err))
			writeWSMessage(ws, gin.H{"error": "目标用户不存在"})
			continue
		}

		// 检查是否已存在好友请求
		var count int64
		global.DB.Model(&models.Friendship{}).
			Where("user_id = ? AND friend_id = ?", uint(user_id), uint(friend_id)).
			Count(&count)
		if count > 0 {
			global.Log.Error("好友请求已存在")
			writeWSMessage(ws, gin.H{"error": "好友请求已存在"})
			continue
		}

		// 创建好友请求记录
		friendship := models.Friendship{
			UserID:    uint(user_id),
			FriendID:  uint(friend_id),
			Status:    "pending",
			Note:      note,
			CreatedAt: time.Now(),
		}

		notifyMsg := map[string]interface{}{
			"type":         "friend_request",
			"from_user_id": userid,
			"to_user_id":   friendid,
			"note":         note,
			"code":         200,
		}
		msgBytes, err := json.Marshal(notifyMsg)
		if err != nil {
			global.Log.Error("消息编码失败", zap.Error(err))
			writeWSMessage(ws, gin.H{"message": "发送好友请求成功，但通知编码失败"})
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
			Type:    "friend_request",
		}

		// 保存好友请求到数据库
		if err = global.DB.Create(&friendship).Error; err != nil {
			global.Log.Error("添加好友失败!", zap.Error(err))
			writeWSMessage(ws, gin.H{"error": "添加好友失败"})
			continue
		}
		err = writeWSMessage(ws, gin.H{"message": "好友请求发送成功"})
		if err != nil {
			global.Log.Error("写入WebSocket失败", zap.Error(err))
		}
	}
}

var writeLock sync.Mutex

func writeWSMessage(ws *websocket.Conn, data interface{}) error {
	writeLock.Lock()
	defer writeLock.Unlock()

	message, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, message)
}
