package api

import (
	"context"
	"errors"
	"gochat/global"
	"gochat/models"
	"gochat/request"
	"gochat/response"
	"gochat/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type LogoutApi struct{}

func (logoutapi *LogoutApi) Logout(c *gin.Context) {
	var req request.LogoutReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tokenString := req.Token
	if tokenString == "" {
		response.FailWithMessage("token为空", c)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("BMvfawCaCjlDOzLAoYDUxLZWGIzerY53VeIm03Fy6uE="), nil
	})
	if err != nil {
		c.Error(errors.New("无效token"))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	var userID string
	switch v := claims["userid"].(type) {
	case float64: // 当userid是JSON数字时
		userID = strconv.Itoa(int(v))
	case string: // 当userid是字符串时
		userID = v
	default:
		global.Log.Error("userid类型无效")
		return
	}
	// 获取过期时间 (exp)
	var expireTime time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expireTime = time.Unix(int64(exp), 0) // 转换为 time.Time
	} else {
		// 如果没有 exp，使用默认过期时间（例如 1 小时）
		expireTime = time.Now().Add(time.Hour)
		global.Log.Warn("token中未找到exp字段,使用默认过期时间")
	}
	if err := addToBlacklist(tokenString, expireTime); err != nil {
		global.Log.Error("token添加黑名单失败", zap.Error(err))
		response.FailWithMessage("token添加黑名单失败", c)
		return
	}
	//修改用户的在线状态
	err = global.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", models.StatusOffline).Error
	if err != nil {
		global.Log.Error("用户状态修改失败：", zap.Error(err))
		response.FailWithMessage("用户状态修改失败", c)
	}
	//断开所有的客户端连接
	id, _ := strconv.ParseUint(userID, 10, 32)
	for _, client := range service.Manager.Clients {
		if client.User != nil && client.User.ID == uint(id) {
			service.Manager.Unregister <- client
			_ = client.Socket.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(1000, "您已下线"),
			)
		}
	}
}

/*
	func extractToken(c *gin.Context) string {
		return strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	}

	func parseToken(token string) (*jwt.RegisteredClaims, error) {
		claims := &jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("BMvfawCaCjlDOzLAoYDUxLZWGIzerY53VeIm03Fy6uE="), nil
		})
		return claims, err
	}
*/
func addToBlacklist(token string, expireTime time.Time) error {
	ttl := time.Until(expireTime)

	return global.RedisClient.Set(
		context.Background(),
		"jwt_blacklist:"+token,
		"1",
		ttl,
	).Err()
}
