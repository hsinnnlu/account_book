let expenses = [];

// 更新總消費
function updateTotalAmount() {
    let total = expenses.reduce((sum, expense) => sum + parseFloat(expense.amount), 0);
    $('#total-amount').text(total);
}

// 根據日期排序支出
function sortExpensesByDate() {
    expenses.sort((a, b) => new Date(b.date) - new Date(a.date));
}

// 計算每日消費總額
function calculateDailyTotal(date) {
    return expenses.filter(expense => expense.date === date).reduce((sum, exp) => sum + parseFloat(exp.amount), 0);
}

// 渲染支出紀錄，並根據日期分界
function renderExpenses(category = 'all') {
    const tableBody = $('#expense-table');
    tableBody.empty();

    sortExpensesByDate();

    let lastDate = '';
    const filteredExpenses = category === 'all' ? expenses : expenses.filter(exp => exp.category === category);
    filteredExpenses.forEach((expense, index) => {
        // 新日期時添加分界線及當日總消費
        if (expense.date !== lastDate) {
            let dailyTotal = calculateDailyTotal(expense.date);
            tableBody.append(`
                <tr class="table-secondary">
                    <td colspan="6">
                        <strong>${expense.date}</strong> - 總消費: $${dailyTotal}
                    </td>
                </tr>
            `);
            lastDate = expense.date;
        }

        tableBody.append(`
            <tr>
                <td>${expense.category}</td>
                <td>${expense.name}</td>
                <td>${expense.date}</td>
                <td>$${expense.amount}</td>
                <td>${expense.payment}</td>
                <td>
                    <button class="btn btn-sm btn-danger delete-expense" data-index="${index}">刪除</button>
                </td>
            </tr>
        `);
    });

    updateTotalAmount();
}

// 新增支出
$('#expense-form').on('submit', function(event) {
    event.preventDefault();

    const name = $('#expense-name').val();
    const category = $('#expense-category').val();
    const amount = $('#expense-amount').val();
    const date = $('#expense-date').val();
    const payment = $('#expense-payment').val();

    expenses.push({ name, category, amount, date, payment });
    renderExpenses();

    // 清空表單
    $('#expense-form')[0].reset();
    $('#expenseModal').modal('hide');
});

// 選擇分類過濾支出
$('#category-menu').on('click', 'a', function() {
    const category = $(this).data('category');
    renderExpenses(category);
});

// 刪除支出
$('#expense-table').on('click', '.delete-expense', function() {
    const index = $(this).data('index');
    expenses.splice(index, 1);
    renderExpenses();
});

// 開啟新增支出Modal
$('#add-expense').on('click', function() {
    $('#expenseModal').modal('show');
});

// 預設渲染全部支出
renderExpenses();
