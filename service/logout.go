package service

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/account_book/db"
)

func LogoutHandler(c *gin.Context) {
	// 清除 session
	session := sessions.Default(c)
	session.Clear() // 清除所有 session 資料
	session.Save()  // 儲存變更

	// 回傳成功訊息
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func ResetpasswordHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, _ := session.Get("user_id").(string)

	var req struct {
		UserID      	string
		NewPassword 	string
		ComNewPassword  string
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供有效的用戶ID和新密碼"})
		return
	}

	hash := sha256.New()
	hash.Write([]byte(req.NewPassword))
	err := db.UpdatePassword(hex.EncodeToString(hash.Sum(nil)), user)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "資料庫更新失敗"})
		return
	}

	// 回傳成功訊息
	c.JSON(http.StatusOK, gin.H{"message": "密碼修改成功"})
}

