package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	_, password, confirmpassword := GetUser(c)
	err := validatePasswordMatch(password, confirmpassword)

	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err.Error(),
		})
	}

	// password = HashedPassword(password)

	//把user_id、hashpasswordpassword 加進db中
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
