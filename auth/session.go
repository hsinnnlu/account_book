package auth

import (
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var store sessions.Store
var sessionIDs sync.Map // 使用 sync.Map 來儲存所有 sessionID

// 初始化 session 中間件
func InitSession(secret string) gin.HandlerFunc {
	store = cookie.NewStore([]byte(secret))
	return sessions.Sessions("mysession", store)
}

// 設置 session
func SetSession(c *gin.Context, key string, value interface{}) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

// 列出所有 sessionID
func ListAllSessionIDs() []string {
	var ids []string
	sessionIDs.Range(func(key, value interface{}) bool {
		ids = append(ids, key.(string))
		return true
	})
	return ids
}

func ClearSession(c *gin.Context) {
	sessions := sessions.Default(c)
	sessions.Clear()
	sessions.Save()
}
