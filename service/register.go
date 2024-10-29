package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/account_book/db"
)

func SignUp(c *gin.Context) {
	user_id, password, confirmpassword := GetUser(c)

	err := CheckUser(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "使用者已經存在",
		})
		return
	}

	// 檢查密碼是否匹配
	if err := validatePasswordMatch(password, confirmpassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "密碼和確認密碼不一致，請重新輸入",
		})
		return
	}

	// 對密碼進行雜湊處理
	hash_password := HashedPassword(password)

	// 插入用戶
	if err := db.InsertUser(user_id, hash_password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "註冊失敗"})
		return
	}

	// 註冊成功的回應
	c.JSON(http.StatusOK, gin.H{"message": "註冊成功，請登入"})
}

func CheckUser(user_id string) error {
	_, err := db.GetUserById(db.DB, user_id)

	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return errors.New("user already exists")
}

func GetUser(c *gin.Context) (user_id, password, compassword string) {
	if in, isExist := c.GetPostForm("user_id"); isExist && in != "" {
		user_id = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入使用者名稱"),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入密碼"),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		compassword = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("必須輸入確認密碼"),
		})
		return
	}
	return user_id, password, compassword
}

func validatePasswordMatch(password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("密碼和確認密碼不一致")
	}
	return nil
}

func HashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
