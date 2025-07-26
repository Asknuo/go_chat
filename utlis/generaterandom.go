package utlis

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecretKey() string {
	key := make([]byte, 32) // HS256 推荐32字节
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}
