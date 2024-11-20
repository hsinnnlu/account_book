package service

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/account_book/db"
	"github.com/hsinnnlu/account_book/models"
)

func ExpenseHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "expenses.html", nil)
}

func ExpenseDataHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	expenses, err := db.GetExpenses(db.DB, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func Deleteexpenserow(c *gin.Context) {
	expenseID := c.Param("id") // 從路徑參數中獲取 income_id

	println("expense_id: ", expenseID)
	err := db.DeleteExpense(db.DB, expenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("刪除失敗: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "刪除成功！",
	})
}

func Insertexpenserow(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	var expense models.Expenses
	err := c.ShouldBindJSON(&expense)
	fmt.Println("Expense data:", expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = db.InsertExpense(db.DB, expense, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收入新增成功", "expense": expense})
}
