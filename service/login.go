// POST: /login 路由
package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登入頁面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
