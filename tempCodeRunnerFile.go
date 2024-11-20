router.GET("/api/expenses", auth.AuthMiddleware(), service.ExpenseDataHandler)
