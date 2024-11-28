package service

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	// 清除 session
	session := sessions.Default(c)
	session.Clear() // 清除所有 session 資料
	session.Save()  // 儲存變更

	// 回傳成功訊息
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}