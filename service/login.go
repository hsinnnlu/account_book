// POST: /login 路由
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/account_book/db"
	"github.com/hsinnnlu/account_book/models"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginAuth(c *gin.Context) {

	input_id, input_password := preProcessingInput(c)
	fmt.Println("Received user_id:", input_id)
	fmt.Println("Received password:", input_password)

	user, err := checkPassword(input_id, input_password)
	log.Print(err)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	c.SetCookie("user_cookie", user.Id, 3600, "/", "localhost", false, true)
	session.Save()

	// 成功處理邏輯，重定向
	c.Redirect(http.StatusFound, "/webpage/account_book.html")
}

func preProcessingInput(c *gin.Context) (user_id, password string) {

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
	return user_id, password
}

// 串資料庫
func checkPassword(user_id, inputPassword string) (*models.User, error) {
	hashedInputPassword := GetHashedPassword(inputPassword)
	
	// 檢查使用者是否存在
	user, err := db.GetUserById(db.DB, user_id)

	if err != nil {
		fmt.Println("error:", err)
		return nil, errors.New("user does not exist")
	}

	storedPasswordHash := strings.TrimSpace(user.Hash_password)
	fmt.Println("db的密碼: ", storedPasswordHash)
	if storedPasswordHash != hashedInputPassword {
		return nil, errors.New("password is incorrect")
	}
	return user, nil
}

// 將密碼使用 SHA256
func GetHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
