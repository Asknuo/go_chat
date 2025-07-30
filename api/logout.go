package api

import (
	"context"
	"gochat/global"
	"gochat/models"
	"gochat/response"
	"gochat/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type LogoutApi struct{}

func (logoutapi *LogoutApi) Logout(c *gin.Context) {
	//获取token
	tokenString := extractToken(c)
	if tokenString == "" {
		response.FailWithMessage("token为空", c)
		return
	}
	//解析token
	claims, err := parseToken(tokenString)
	if err != nil {
		global.Log.Error("解析错误 :", zap.Error(err))
		response.FailWithMessage("解析错误", c)
		return
	}
	//把token加入黑名单
	userID := claims.Subject
	if userID == "" {
		response.FailWithMessage("无效的token", c)
		return
	}
	if err := addToBlacklist(tokenString, claims.ExpiresAt.Time); err != nil {
		global.Log.Error("token添加黑名单失败", zap.Error(err))
		response.FailWithMessage("token添加黑名单失败", c)
		return
	}
	//修改用户的在线状态
	err = global.DB.Where("id = ?", userID).Update("status", models.StatusOffline).Error
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
func addToBlacklist(token string, expireTime time.Time) error {
	ttl := time.Until(expireTime)

	return global.RedisClient.Set(
		context.Background(),
		"jwt_blacklist:"+token,
		"1",
		ttl,
	).Err()
}
