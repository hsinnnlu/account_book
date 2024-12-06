package main

import (
	"github.com/gin-gonic/gin"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hsinnnlu/account_book/auth"
	"github.com/hsinnnlu/account_book/db"
	"github.com/hsinnnlu/account_book/service"
)

func main() {
	db.InitDB()

	router := gin.Default()
	router.Use(auth.InitSession("secret"))

	router.Static("/js", "./js")
	router.Static("/css", "./css")
	router.Static("/img", "./img")

	router.LoadHTMLGlob(("webpage/*"))

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)
	router.POST("/register", service.SignUp)

	authorized := router.Group("/")
	authorized.Use(auth.AuthMiddleware())
	{
		// 網頁相關路由
		authorized.GET("/webpage/account_book.html", service.AccountBookHandler)
		authorized.GET("/webpage/expenses.html", service.ExpenseHandler)
		authorized.GET("/webpage/income.html", service.IncomeHandler)
		authorized.GET("/webpage/chart.html", service.ChartHandler)

		// API 路由
		authorized.GET("/api/incomes", service.IncomeDataHandler)
		authorized.DELETE("/api/incomes/:id", service.Deleteincomerow)
		authorized.POST("/api/incomes/insertincome", service.Insertincomerow)

		authorized.GET("/api/expenses", service.ExpenseDataHandler)
		authorized.DELETE("/api/expenses/:id", service.Deleteexpenserow)
		authorized.POST("/api/expenses/insertexpense", service.Insertexpenserow)

		authorized.GET("/api/incomechart", service.ChartIncomeDataHandler)
		authorized.GET("/api/expensechart", service.ChartExpenseDataHandler)

		authorized.POST("/api/logout", service.LogoutHandler)
		authorized.POST("/api/reset-password", service.ResetpasswordHandler)
	}

	router.Run(":8080")
}
