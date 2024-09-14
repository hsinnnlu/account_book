// POST: /login 路由
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginAuth(c *gin.Context) {
	input_id, input_password := preProcessingInput(c)
	fmt.Println("Received user_id:", input_id)
	fmt.Println("Received password:", input_password)

	err := checkPassword(input_id, input_password)
	log.Print(err)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 成功處理邏輯，例如重定向
	c.Redirect(http.StatusSeeOther, "/webpage/account_book.html")
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
func checkPassword(user_id, inputPassword string) error {
	hashedInputPassword := GetHashedPassword(inputPassword)
	fmt.Println(hashedInputPassword) // 這裡會印出輸入密碼的 SHA256 雜湊值

	if user_id != "test" {
		return errors.New("user does not exist")
	}

	storedPasswordHash := GetHashedPassword("123")
	userpassword := GetHashedPassword(inputPassword)

	if userpassword != storedPasswordHash {
		return errors.New("password is incorrect")
	}

	// 檢查使用者是否存在
	// user, err := db.GetUserById(DB, user_id)
	// fmt.Println("pass B: ", user.Password_hash) // 這裡會印出資料庫中的密碼雜湊值
	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return nil, errors.New("user does not exist")
	// }

	// storedPasswordHash := user.Password_hash
	// if storedPasswordHash != hashedInputPassword {
	// 	return nil, errors.New("password is incorrect")
	// }
	return nil
}

// 將密碼使用 SHA256
func GetHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
