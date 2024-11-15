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
	// router.Static("/webpage", "./webpage")

	router.LoadHTMLGlob(("webpage/*"))

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)
	router.POST("/register", service.SignUp)

	router.GET("/webpage/account_book.html", auth.AuthMiddleware(), service.AccountBookHandler)
	router.GET("/webpage/expense.html", auth.AuthMiddleware(), service.ExpenseHandler)
	router.GET("/webpage/income.html", auth.AuthMiddleware(), service.IncomeHandler)

	router.Run(":8080")
}
