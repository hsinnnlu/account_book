let incomes = [];

// 更新總收入
function updateTotalAmount() {
    let total = incomes.reduce((sum, income) => sum + parseFloat(income.amount), 0);
    $('#total-amount').text(total);
}

// 根據日期排序收入
function sortIncomesByDate() {
    incomes.sort((a, b) => new Date(b.date) - new Date(a.date));
}

// 計算每日收入總額
function calculateDailyTotal(date) {
    return incomes.filter(income => income.date === date).reduce((sum, inc) => sum + parseFloat(inc.amount), 0);
}

// 渲染收入紀錄，並根據日期分界
function renderIncomes(category = 'all') {
    const tableBody = $('#income-table');
    tableBody.empty();

    sortIncomesByDate();

    let lastDate = '';
    const filteredIncomes = category === 'all' ? incomes : incomes.filter(inc => inc.category === category);
    filteredIncomes.forEach((income, index) => {
        // 新日期時添加分界線及當日總收入
        if (income.date !== lastDate) {
            let dailyTotal = calculateDailyTotal(income.date);
            tableBody.append(`
                <tr class="table-secondary">
                    <td colspan="6">
                        <strong>${income.date}</strong> - 總收入: $${dailyTotal}
                    </td>
                </tr>
            `);
            lastDate = income.date;
        }

        tableBody.append(`
            <tr>
                <td>${income.date}</td>
                <td>${income.category}</td>
                <td>$${income.amount}</td>
                <td>${income.payment}</td>
                <td>${income.note}</td> <!-- 顯示備註 -->
                <td>
                    <button class="btn btn-sm btn-danger delete-income" data-index="${index}">刪除</button>
                </td>
            </tr>
        `);
    });

    updateTotalAmount();
}

// 新增收入
$('#income-form').on('submit', function(event) {
    event.preventDefault();

    const date = $('#income-date').val();
    const category = $('#income-category').val();
    const amount = $('#income-amount').val();
    const payment = $('#income-payment').val();
    const note = $('#income-note').val(); // 取得備註

    incomes.push({ date, category, amount, payment, note }); // 加入備註
    renderIncomes();

    // 清空表單
    $('#income-form')[0].reset();
    $('#incomeModal').modal('hide');
});

// 選擇分類過濾收入
$('#category-menu').on('click', 'a', function() {
    const category = $(this).data('category');
    renderIncomes(category);
});

// 刪除收入
$('#income-table').on('click', '.delete-income', function() {
    const index = $(this).data('index');
    incomes.splice(index, 1);
    renderIncomes();
});

// 開啟新增收入Modal
$('#add-income').on('click', function() {
    $('#incomeModal').modal('show');
});

// 預設渲染全部收入
renderIncomes();
