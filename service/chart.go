package service

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsinnnlu/account_book/db"
)

// 渲染 income.html 頁面
func ChartHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "chart.html", nil)
}

func ChartIncomeDataHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	chart_income, err := db.GetIncomeSummary(db.DB, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 構造返回的 JSON 結構
	categories := make(map[string]int)
	total := 0
	for _, income := range chart_income {
		amount, _ := strconv.Atoi(income.Income_amount) // 將金額轉換為整數
		categories[income.Income_catagory] = amount
		total += amount
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      total,
	})
}

func ChartExpenseDataHandler(c *gin.Context) {
	session := sessions.Default(c)
	user, exists := session.Get("user_id").(string)

	if !exists || user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入"})
		return
	}

	chart_expense, err := db.GetExpenseSummary(db.DB, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 構造返回的 JSON 結構
	categories := make(map[string]int)
	total := 0
	for _, expense := range chart_expense {
		amount, _ := strconv.Atoi(expense.Expense_amount) // 將金額轉換為整數
		categories[expense.Expense_catagory] = amount
		total += amount
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      total,
	})
}
