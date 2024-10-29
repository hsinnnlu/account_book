package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hsinnnlu/account_book/models"
)

var DB *sql.DB // 全局變量, 小寫db是外地注入變量

func InitDB() {
	connectionString := "server=127.0.0.1,1433;trusted_connection=yes;database=Accountbook;"

	var err error
	DB, err = sql.Open("mssql", connectionString)
	if err != nil {
		fmt.Println("連接資料庫失敗:", err.Error())
		return
	}

	// 測試實際連線
	if err := DB.Ping(); err != nil {
		fmt.Println("無法連線到資料庫:", err.Error())
		return
	}

	fmt.Println("資料庫連接成功")
}

func GetUserById(db *sql.DB, user_id string) (*models.User, error) {
	user := models.User{}
	query := "SELECT user_id, hash_password FROM users WHERE user_id = ?"
	err := DB.QueryRow(query, user_id).Scan(&user.Id, &user.Hash_password)
	if err != nil {
		// 如果找不到使用者
		if err == sql.ErrNoRows {
			fmt.Println("user not found")
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// 插入使用者帳號與密碼
func InsertUser(user_id, hash_password string) error {
    query := "INSERT INTO users (user_id, hash_password) VALUES (?, ?)"

    // 使用 ExecContext 並傳入位置參數
    _, err := DB.ExecContext(context.Background(), query, user_id, hash_password)
    if err != nil {
        return fmt.Errorf("執行插入語句失敗: %v", err)
    }

    fmt.Println("使用者插入成功")
    return nil
}