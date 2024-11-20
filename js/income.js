// 收入資料存放變數
let incomes = [];

// 從 API 獲取收入資料
function fetchIncomes() {
    fetch('/api/incomes')
        .then(response => response.json())
        .then(data => {
            console.log(data); // 檢查 API 返回的資料是否正確
            if (data.error) {
                console.error('無法獲取資料:', data.error);
                return;
            }

            incomes = data.incomes;
            renderIncomes(); // 渲染收入紀錄
        })
        .catch(error => console.error('獲取收入資料時出錯:', error));
}

function updateTotalAmount() {
    let total = incomes.reduce((sum, income) => sum + parseFloat(income.Amount || 0), 0);
    $('#total-amount').text(total);
}

// 根據日期排序收入
function sortIncomesByDate() {
    incomes.sort((a, b) => new Date(b.Date) - new Date(a.Date));
}

// 計算每日總收入
function calculateDailyTotal(date) {
    return incomes
        .filter(income => income.Date === date) // 過濾特定日期的收入
        .reduce((sum, income) => sum + parseFloat(income.Amount || 0), 0); // 計算總和
}

function renderIncomes(category = 'all') {
    const tableBody = $('#income-table tbody');
    tableBody.empty(); // 清空現有表格內容

    sortIncomesByDate(); // 先按日期排序

    let lastDate = ''; // 用於記錄上一次渲染的日期
    const filteredIncomes = category === 'all' 
        ? incomes 
        : incomes.filter(income => income.Income_category.trim() === category);

    filteredIncomes.forEach((income, index) => {
        // 日期分界線與每日總收入
        if (income.Date !== lastDate) {
            const dailyTotal = calculateDailyTotal(income.Date); // 計算每日總收入
            tableBody.append(`
                <tr class="table-secondary">
                    <td colspan="6">
                        <strong>${income.Date}</strong> - 總收入: $${dailyTotal}
                    </td>
                </tr>
            `);
            lastDate = income.Date; // 更新記錄的日期
        }

        // 添加收入資料
        tableBody.append(`
            <tr>
                <td>${income.Date}</td>
                <td>${income.Income_category.trim()}</td>
                <td>${income.Amount}</td>
                <td>${income.Account.trim()}</td>
                <td>${income.Memo.trim() || ''}</td>
                <td>
                <button class="btn btn-sm btn-danger delete-income" data-id="${income.Id}">刪除</button>
                </td>
            </tr>
        `);
    });

    updateTotalAmount(); // 更新總收入
}

$('#income-table').on('click', '.delete-income', function () {
    const $button = $(this);
    const incomeId = $button.data('id');

    if (!confirm('確定要刪除此收入紀錄嗎？')) {
        return;
    }

    fetch(`/api/incomes/${incomeId}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || '刪除失敗');
                });
            }
            return response.json();
        })
        .then(data => {
            fetchIncomes(); // 刪除成功後重新渲染表格
            alert(data.message || '刪除成功！');
        })
        .catch(error => {
            console.error('刪除收入時出錯:', error);
            alert('刪除收入失敗: ' + error.message);
        });
});

$('#income-form').on('submit', function(event) {
    event.preventDefault();

    const date = $('#income-date').val();
    const category = $('#income-category').val();
    const amount = parseFloat($('#income-amount').val());
    const payment = $('#income-payment').val();
    const note = $('#income-note').val();

    // 確保表單資料有效
    if (!date || !category || !amount || !payment) {
        alert('請填寫所有必要欄位');
        return;
    }

    const formData = {
        Date: date,
        Income_category: category,
        Amount: amount,
        Account: payment,
        Memo: note || "",  // 備註可以是空的
    };
    console.log('formData: ', formData)

    fetch('/api/incomes/insertincome', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData), // 這裡傳送的是表單的資料
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert('新增收入失敗: ' + data.error);
            return;
        }

        incomes.push(data.income); // 更新本地資料
        fetchIncomes(); // 新增後重新渲染
        $('#income-form')[0].reset(); // 清空表單
        $('#incomeModal').modal('hide'); // 關閉 Modal
    })
    .catch(error => console.error('新增收入時出錯:', error));
});

// 開啟新增收入 Modal
$('#add-income').on('click', function() {
    $('#incomeModal').modal('show');
});

// 頁面加載後初始化
$(document).ready(() => {
    fetchIncomes(); // 獲取收入資料並渲染
});
