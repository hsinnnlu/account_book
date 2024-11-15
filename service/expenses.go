package service

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ExpenseHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		// 如果 session 中沒有 user，重定向回登入頁面
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "expenses.html", gin.H{
		
	})
}