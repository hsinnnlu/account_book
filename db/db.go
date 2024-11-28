package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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

// 取得 DB 裡面的 income 資料
func GetIncome(db *sql.DB, user_id string) ([]models.Income, error) {
	query := `
		SELECT 
			i.income_id,
			CONVERT(varchar, i.income_date, 23) AS income_date, 
			ic.category AS income_category, 
			i.amount, 
			a.account AS account_name, 
			i.memo
		FROM income i
		JOIN income_category ic ON i.income_category = ic.income_category
		JOIN account a ON i.account_id = a.account_id
		WHERE i.user_id = ?
	`

	// 執行查詢
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("查詢失敗: %v", err)
	}
	defer rows.Close()

	// 儲存收入資料
	var incomes []models.Income

	// 逐列處理查詢結果
	for rows.Next() {
		var income models.Income
		var memo sql.NullString // 使用 NullString 處理可能為 NULL 的欄位

		// 扫描資料列
		err := rows.Scan(&income.Id, &income.Date, &income.Income_category, &income.Amount, &income.Account, &memo)
		if err != nil {
			return nil, fmt.Errorf("掃描資料失敗: %v", err)
		}

		// 如果 memo 有值，則設置給 income.Memo，否則設為空字串
		if memo.Valid {
			income.Memo = memo.String
		} else {
			income.Memo = "" // 如果 memo 是 NULL，設為空字串
		}

		// 將資料加入到 incomes 列表
		incomes = append(incomes, income)

		// 查看處理的資料
		log.Printf("收入: %+v\n", income)
	}

	// 檢查查詢結果是否為空
	if len(incomes) == 0 {
		log.Println("查詢結果為空")
	}

	// 返回查詢結果
	return incomes, nil
}

// 刪除一筆收入資料
func DeleteIncome(db *sql.DB, income_id string) error {
	query := "DELETE FROM income WHERE income_id = ?"
	result, err := db.Exec(query, income_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return err
	}
	return nil
}

// InsertIncome 將單筆收入資料插入資料庫
func InsertIncome(db *sql.DB, income models.Income, user_id string) error {
	// 查詢當前最大的 income_id
	var maxIncomeID int
	queryMaxID := "SELECT COALESCE(MAX(CAST(income_id AS INT)), 0) FROM income"
	err := db.QueryRow(queryMaxID).Scan(&maxIncomeID)
	if err != nil {
		return fmt.Errorf("查詢最大收入編號失敗: %w", err)
	}

	// 設定新的 income_id
	newIncomeID := maxIncomeID + 1
	income.Id = fmt.Sprintf("%d", newIncomeID)

	// 查詢 Income_category 對應的 ID
	var categoryID string
	categoryQuery := "SELECT income_category FROM income_category WHERE category = ?"
	err = db.QueryRow(categoryQuery, income.Income_category).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("找不到對應的收入分類: %s", income.Income_category)
		}
		return fmt.Errorf("查詢收入分類失敗: %w", err)
	}

	// 查詢 Account 對應的 ID
	var accountID string
	accountQuery := "SELECT account_id FROM account WHERE account = ?"
	err = db.QueryRow(accountQuery, income.Account).Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("找不到對應的帳戶: %s", income.Account)
		}
		return fmt.Errorf("查詢帳戶失敗: %w", err)
	}

	income.User_id = user_id

	// 插入收入資料
	insertQuery := `
		INSERT INTO income (income_id, user_id, income_category, account_id, income_date, amount, memo)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(
		insertQuery,
		income.Id,
		income.User_id,
		categoryID, // 使用查詢的 categoryID
		accountID,  // 使用查詢的 accountID
		income.Date,
		income.Amount,
		income.Memo,
	)
	if err != nil {
		return fmt.Errorf("插入收入資料失敗: %w", err)
	}

	return nil
}

func InsertIncomeCatagory(db *sql.DB, newCatagory string) error {
	// 查詢目前最大的 income_catagory
	var maxCatagory int
	query := `SELECT CAST(income_catagory AS INT) AS max_catagory 
              FROM income_catagory 
              ORDER BY max_catagory DESC 
              LIMIT 1`
	err := db.QueryRow(query).Scan(&maxCatagory)
	if err != nil {
		// 如果查詢失敗但沒有記錄，表示資料表可能是空的
		if err == sql.ErrNoRows {
			maxCatagory = 0
		} else {
			return err
		}
	}

	// 計算下一個 income_catagory
	newIncomeCatagory := maxCatagory + 1

	// 插入新資料
	insertQuery := `INSERT INTO income_catagory (income_catagory, catagory) VALUES (?, ?)`
	_, err = db.Exec(insertQuery, fmt.Sprintf("%d", newIncomeCatagory), newCatagory)
	if err != nil {
		return err
	}

	fmt.Printf("新增成功，income_catagory: %d, catagory: %s\n", newIncomeCatagory, newCatagory)
	return nil
}

func GetExpenses(db *sql.DB, user_id string) ([]models.Expenses, error) {
	query := `
		SELECT 
			e.expense_id,
			CONVERT(varchar, e.expense_date, 23) AS expense_date, 
			ec.category AS expense_category, 
			e.item,
			e.amount, 
			a.account AS account_name 
		FROM user_expense e
		JOIN expense_category ec ON e.expense_category = ec.expense_category
		JOIN account a ON e.account_id = a.account_id
		WHERE e.user_id = ?
	`

	// 執行查詢
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("查詢失敗: %v", err)
	}
	defer rows.Close()

	// 儲存支出資料
	var expenses []models.Expenses

	// 逐列處理查詢結果
	for rows.Next() {
		var expense models.Expenses
		
		// 扫描資料列
		err := rows.Scan(&expense.Id, &expense.Date, &expense.Expense_category, &expense.Item, &expense.Amount, &expense.Account)
		if err != nil {
			return nil, fmt.Errorf("掃描資料失敗: %v", err)
		}

		// 將資料加入到 expenses 列表
		expenses = append(expenses, expense)

		// 查看處理的資料
		log.Printf("支出: %+v\n", expense)
	}

	// 檢查查詢結果是否為空
	if len(expenses) == 0 {
		log.Println("查詢結果為空")
	}

	// 返回查詢結果
	return expenses, nil
}

// 刪除一筆收入資料
func DeleteExpense(db *sql.DB, expense_id string) error {
	query := "DELETE FROM user_expense WHERE expense_id = ?"
	result, err := db.Exec(query, expense_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return err
	}
	return nil
}

// InsertExpense 將單筆支出資料插入資料庫
func InsertExpense(db *sql.DB, expense models.Expenses, user_id string) error {
	// 查詢當前最大的 expense_id
	var maxExpenseID int
	queryMaxID := "SELECT COALESCE(MAX(CAST(expense_id AS INT)), 0) FROM user_expense"
	err := db.QueryRow(queryMaxID).Scan(&maxExpenseID)
	if err != nil {
		return fmt.Errorf("查詢最大支出編號失敗: %w", err)
	}

	// 設定新的 expense_id
	newExpenseID := maxExpenseID + 1
	expense.Id = fmt.Sprintf("%d", newExpenseID)

	// 查詢 Expense_category 對應的 ID
	var categoryID string
	categoryQuery := "SELECT expense_category FROM expense_category WHERE category = ?"
	err = db.QueryRow(categoryQuery, expense.Expense_category).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("找不到對應的支出分類: %s", expense.Expense_category)
		}
		return fmt.Errorf("查詢支出分類失敗: %w", err)
	}

	// 查詢 Account 對應的 ID
	var accountID string
	accountQuery := "SELECT account_id FROM account WHERE account = ?"
	err = db.QueryRow(accountQuery, expense.Account).Scan(&accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("找不到對應的帳戶: %s", expense.Account)
		}
		return fmt.Errorf("查詢帳戶失敗: %w", err)
	}

	expense.User_id = user_id

	// 插入支出資料
	insertQuery := `
		INSERT INTO user_expense (expense_id, user_id, expense_category, account_id, expense_date, amount, item)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(
		insertQuery,
		expense.Id,
		expense.User_id,
		categoryID, // 使用查詢的 categoryID
		accountID,  // 使用查詢的 accountID
		expense.Date,
		expense.Amount,
		expense.Item, // 支出名稱
	)
	if err != nil {
		return fmt.Errorf("插入支出資料失敗: %w", err)
	}

	return nil
}

// chart資料
func GetExpenseSummary(db *sql.DB, userID string) ([]models.Chart_expense, error) {
	query := `
		SELECT 
			ec.category AS expense_category, 
			COALESCE(SUM(e.amount), 0) AS total_amount
		FROM expense_category ec
		LEFT JOIN user_expense e 
			ON e.expense_category = ec.expense_category AND e.user_id = ?
		GROUP BY ec.category
		ORDER BY total_amount DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查詢支出分類匯總失敗: %v", err)
	}
	defer rows.Close()

	var expenses []models.Chart_expense
	for rows.Next() {
		var expense models.Chart_expense
		var totalAmount int

		// 扫描資料列
		err := rows.Scan(&expense.Expense_catagory, &totalAmount)
		if err != nil {
			return nil, fmt.Errorf("掃描支出分類匯總資料失敗: %v", err)
		}

		// 將數字轉換為字串並存入結構體
		expense.Expense_amount = fmt.Sprintf("%d", totalAmount)
		expenses = append(expenses, expense)
	}

	for _, expense := range expenses {
		fmt.Printf("分類: %s, 總金額: %s\n", expense.Expense_catagory, expense.Expense_amount)
	}

	if len(expenses) == 0 {
		log.Println("支出分類匯總結果為空")
	}
	return expenses, nil
}

func GetIncomeSummary(db *sql.DB, userID string) ([]models.Chart_income, error) {
	query := `
		SELECT 
			ic.category AS income_category, 
			COALESCE(SUM(i.amount), 0) AS total_amount
		FROM income_category ic
		LEFT JOIN income i 
			ON i.income_category = ic.income_category AND i.user_id = ?
		GROUP BY ic.category
		ORDER BY total_amount DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查詢收入分類匯總失敗: %v", err)
	}
	defer rows.Close()

	var incomes []models.Chart_income
	for rows.Next() {
		var income models.Chart_income
		var totalAmount int

		// 扫描資料列
		err := rows.Scan(&income.Income_catagory, &totalAmount)
		if err != nil {
			return nil, fmt.Errorf("掃描收入分類匯總資料失敗: %v", err)
		}

		// 將數字轉換為字串並存入結構體
		income.Income_amount = fmt.Sprintf("%d", totalAmount)
		incomes = append(incomes, income)
	}

	for _, income := range incomes {
		fmt.Printf("分類: %s, 總金額: %s\n", income.Income_catagory, income.Income_amount)
	}

	if len(incomes) == 0 {
		log.Println("收入分類匯總結果為空")
	}
	return incomes, nil
}
