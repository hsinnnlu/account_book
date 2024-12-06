package service

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/account_book/db"
	"github.com/hsinnnlu/account_book/models"
)

// 渲染 income.html 頁面
func IncomeHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "income.html", nil)
}

// 提供 incomes 的 JSON 資料
func IncomeDataHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	incomes, err := db.GetIncome(db.DB, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"incomes": incomes})
}

func Deleteincomerow(c *gin.Context) {
	incomeID := c.Param("id") // 從路徑參數中獲取 income_id

	println("income_id: ", incomeID)
	str := db.DeleteIncome(db.DB, incomeID)
	if str == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "刪除失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": str,
	})
}

func Insertincomerow(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	var income models.Income
	err := c.ShouldBindJSON(&income)
	fmt.Println("Income data:", income)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = db.InsertIncome(db.DB, income, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收入新增成功", "income": income})
}
