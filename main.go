package main

import (
	"github.com/gin-gonic/gin"

	"github.com/hsinnnlu/account_book/service"
)

func main() {

	router := gin.Default()

	router.Static("/js", "./js")
	router.Static("/css", "./css")
	router.Static("/img", "./img")
	router.Static("/webpage", "./webpage") // 確保此路徑對應到存放 HTML 文件的目錄

	router.LoadHTMLGlob(("webpage/*"))

	router.GET("/login", service.LoginPage)
	router.POST("/login", service.LoginAuth)
	router.POST("/register", service.SignUp)

	// router.POST("/test", service.TestDB)

	router.Run(":8080")
}
