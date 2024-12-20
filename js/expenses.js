let expenses = [];

// 從 API 獲取收入資料
function fetchExpenses() {
    fetch('/api/expenses')
        .then(response => response.json())
        .then(data => {
            console.log(data); // 檢查 API 返回的資料是否正確
            if (data.error) {
                console.error('無法獲取資料:', data.error);
                return;
            }

            expenses = data.expenses;
            renderExpenses(); // 渲染收入紀錄
        })
        .catch(error => console.error('獲取支出資料時出錯:', error));
}

function updateTotalExpense() {
    let total = expenses.reduce((sum, expense) => sum + parseFloat(expense.Amount || 0), 0);
    $('#total-amount').text(total);
}

// 根據日期排序支出
function sortExpensesByDate() {
    expenses.sort((a, b) => new Date(b.Date) - new Date(a.Date));
}

// 計算每日總支出
function calculateDailyExpense(date) {
    return expenses
        .filter(expense => expense.Date === date) // 過濾特定日期的支出
        .reduce((sum, expense) => sum + parseFloat(expense.Amount || 0), 0); // 計算總和
}

function renderExpenses(category = 'all') {
    const tableBody = $('#expense-table tbody');
    tableBody.empty(); // 清空現有表格內容

    sortExpensesByDate(); // 先按日期排序

    let lastDate = ''; // 用於記錄上一次渲染的日期
    const filteredExpenses = category === 'all' 
        ? expenses 
        : expenses.filter(expense => expense.Expense_category.trim() === category);

    filteredExpenses.forEach((expense, index) => {
        // 日期分界線與每日總支出
        if (expense.Date !== lastDate) {
            const dailyTotal = calculateDailyExpense(expense.Date); // 計算每日總支出
            tableBody.append(`
                <tr class="table-secondary">
                    <td colspan="6">
                        <strong>${expense.Date}</strong> - 總支出: $${dailyTotal}
                    </td>
                </tr>
            `);
            lastDate = expense.Date; // 更新記錄的日期
        }

        // 添加支出資料
        tableBody.append(`
            <tr>
                <td>${expense.Date}</td>
                <td>${expense.Expense_category.trim()}</td>
                <td>${expense.Item.trim()}</td>
                <td>${expense.Amount}</td>
                <td>${expense.Account.trim()}</td>
                <td>
                <button class="btn btn-sm btn-danger delete-expense" data-id="${expense.Id}">刪除</button>
                </td>
            </tr>
        `);
    });

    updateTotalExpense(); // 更新總支出
}

// 刪除支出
$('#expense-table').on('click', '.delete-expense', function () {
    const $button = $(this);
    const expenseId = $button.data('id');

    // 確認刪除
    if (!confirm('確定要刪除此支出紀錄嗎？')) {
        return;
    }

    // 發送 DELETE 請求
    fetch(`/api/expenses/${expenseId}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })
        .then(response => {
            // 處理非 200-299 的 HTTP 狀態碼
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || '刪除失敗');
                });
            }
            return response.json();
        })
        .then(data => {
            // 顯示後端返回的成功訊息
            alert(data.message || '刪除成功！');
            // 重新渲染表格
            fetchExpenses();
        })
        .catch(error => {
            // 顯示後端返回的錯誤訊息
            console.error('刪除支出時出錯:', error);
            alert('刪除支出失敗: ' + error.message);
        });
});


$('#expense-form').on('submit', function(event) {
    event.preventDefault(); // 阻止表單默認提交

    // 獲取表單欄位的值
    const date = $('#expense-date').val();
    const category = $('#expense-category').val();
    const amount = parseFloat($('#expense-amount').val());
    const payment = $('#expense-payment').val();
    const expenseName = $('#expense-name').val();  // 支出名稱

    // 確保表單資料有效
    if (!date || !category || !amount || !payment || !expenseName) {
        alert('請填寫所有必要欄位');
        return;
    }

    // 構造要提交的資料
    const formData = {
        Date: date,
        Expense_category: category,
        Amount: amount,
        Account: payment,
        Item: expenseName,  // 支出名稱
    };

    console.log('formData: ', formData);  // 用於調試

    // 發送 POST 請求到後端
    fetch('/api/expenses/insertexpense', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',  // 設置請求頭為 JSON
        },
        body: JSON.stringify(formData), // 將資料轉換為 JSON 字符串並傳送
    })
    .then(response => response.json()) // 處理回應
    .then(data => {
        if (data.error) {
            alert('新增支出失敗: ' + data.error);  // 顯示錯誤訊息
            return;
        }

        expenses.push(data.expense); // 更新本地資料（如果後端回傳成功資料）
        fetchExpenses(); // 重新渲染支出資料
        $('#expense-form')[0].reset(); // 清空表單
        $('#expenseModal').modal('hide'); // 關閉 Modal
    })
    .catch(error => console.error('新增支出時出錯:', error)); // 處理錯誤
});

// 開啟新增支出Modal
$('#add-expense').on('click', function() {
    $('#expenseModal').modal('show');
});

// 頁面加載後初始化
$(document).ready(() => {
    fetchExpenses(); // 獲取收入資料並渲染
});

document.getElementById("logoutBtn").addEventListener("click", function(event) {
    event.preventDefault();  // 防止頁面跳轉

    // 發送登出請求給後端
    fetch("/api/logout", {
        method: "POST",  // 使用 POST 請求
        headers: {
            "Content-Type": "application/json",
        },
    })
    .then(response => {
        if (response.ok) {
            // 如果登出成功，重定向到 login 頁面
            window.location.href = "/login";  // 這裡假設登入頁面的 URL 是 /login
        } else {
            console.error("登出失敗");
        }
    })
    .catch(error => console.error("Error logging out:", error));
});

document.getElementById("logoutBtn").addEventListener("click", function(event) {
    event.preventDefault();  // 防止頁面跳轉

    // 發送登出請求給後端
    fetch("/api/logout", {
        method: "POST",  // 使用 POST 請求
        headers: {
            "Content-Type": "application/json",
        },
    })
    .then(response => {
        if (response.ok) {
            // 如果登出成功，重定向到 login 頁面
            window.location.href = "/login";  // 這裡假設登入頁面的 URL 是 /login
        } else {
            console.error("登出失敗");
        }
    })
    .catch(error => console.error("Error logging out:", error));
});

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('submitChangePassword').addEventListener('click', function () {
        const newPassword = document.getElementById('newPassword').value;
        const confirmNewPassword = document.getElementById('confirmNewPassword').value;

        // 確認輸入框內容
        if (!newPassword || !confirmNewPassword) {
            alert('請輸入新密碼和確認密碼。');
            return;
        }

        if (newPassword !== confirmNewPassword) {
            alert('新密碼和確認密碼不一致，請重新輸入。');
            return;
        }

        // 發送 API 請求
        fetch('/api/reset-password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                NewPassword: newPassword,
                ComNewPassword: confirmNewPassword,
            }),
        })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(data => {
                        throw new Error(data.error || '密碼修改失敗');
                    });
                }
                return response.json();
            })
            .then(data => {
                // 成功提示
                alert(data.message || '密碼修改成功！');
                // 關閉模態框
                const modal = bootstrap.Modal.getInstance(document.getElementById('changePasswordModal'));
                modal.hide();
            })
            .catch(error => {
                console.error('修改密碼時出錯:', error);
                alert('密碼修改失敗: ' + error.message);
            });
    });
});