package service

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/global"
	"gochat/models"
	"gochat/utlis"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

type Msg struct {
	Target  string `json:"target"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Code    int    `json:"code"`
}

type Client struct {
	User       *models.User
	SenderID   string
	ReveiverID string
	Socket     *websocket.Conn
	Send       chan []byte
}

type Boardcast struct {
	Client  *Client
	Message []byte
	Type    string
}

type ClientManager struct {
	Clients    map[string]*Client
	Boardcast  chan *Boardcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

const month = time.Hour * 24 * 30 //一个月

var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Boardcast:  make(chan *Boardcast, 60),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handler(c *gin.Context) {
	token := c.GetHeader("Sec-WebSocket-Protocol")
	if token == "" {
		c.Error(errors.New("未提供Token"))
		return
	}
	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("BMvfawCaCjlDOzLAoYDUxLZWGIzerY53VeIm03Fy6uE="), nil // 替换为你的实际密钥
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
	var uid string
	switch v := claims["userid"].(type) {
	case float64: // 当userid是JSON数字时
		uid = strconv.Itoa(int(v))
	case string: // 当userid是字符串时
		uid = v
	default:
		log.Println("userid类型无效")
		return
	}
	var user models.User
	err = global.DB.Where("id = ?", uid).First(&user).Error
	if err != nil {
		global.Log.Error("在数据库未找到该用户!")
	}
	touid := c.Query("touid")
	if touid == "" {
		touid = "0"
	}
	clientProtocols := websocket.Subprotocols(c.Request)
	conn, err := (&websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1014,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: clientProtocols,
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
	}
	//创建用户实例
	client := &Client{
		User:       &user,
		SenderID:   CreateID(uid, touid),
		ReveiverID: CreateID(touid, uid),
		Socket:     conn,
		Send:       make(chan []byte),
	}
	//注册到用户管理
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(Msg)
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			global.Log.Error("读取错误")
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		if sendMsg.Type == "private" {
			r1, _ := global.RedisClient.Get(context.Background(), c.SenderID).Result()
			r2, _ := global.RedisClient.Get(context.Background(), c.ReveiverID).Result()
			if r1 > "3" && r2 == "" {
				replymsg := &Msg{
					Code:    utlis.WebsocketLimit,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replymsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				global.RedisClient.Incr(context.Background(), c.SenderID)
				_, _ = global.RedisClient.Expire(context.Background(), c.SenderID, month*3).Result()
			}
			Manager.Boardcast <- &Boardcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
				Type:    "private",
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	// 使用 for range 自动处理 channel 数据
	for message := range c.Send {
		replyMsg := &Msg{
			Code:    utlis.WebsocketSuccessMessage,
			Content: string(message), // 直接转为 string
		}
		msg, err := json.Marshal(replyMsg)
		if err != nil {
			log.Printf("JSON 编码失败: %v", err)
			continue
		}
		if err := c.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("消息发送失败: %v", err)
			return // 发送失败时直接退出
		}
	}
	// 如果 c.Send 被关闭，发送 WebSocket 关闭消息
	_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
}
