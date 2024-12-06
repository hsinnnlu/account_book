document.addEventListener("DOMContentLoaded", function () {
    // 從 API 獲取收入數據
    fetch("/api/incomechart")
        .then(response => response.json())
        .then(data => {
            const incomeCategories = Object.keys(data.categories);
            const incomeValues = Object.values(data.categories);

            // 渲染收入圓餅圖
            const incomeCtx = document.getElementById("incomeChart").getContext("2d");
            new Chart(incomeCtx, {
                type: "pie",
                data: {
                    labels: incomeCategories,
                    datasets: [
                        {
                            data: incomeValues,
                            backgroundColor: ["#AEDFF7", "#FFEB99", "#F7C4AE"],
                            borderColor: "#ffffff",
                            borderWidth: 2,
                        },
                    ],
                },
                options: {
                    plugins: {
                        legend: {
                            display: false,
                        },
                    },
                },
            });

            // 更新收入列表數據
            $(".income-details").html(`
                <li>總收入：${data.total}</li>
                ${incomeCategories.map((cat, i) => `
                    <li>
                        <span class="color-box" style="background-color: ${["#AEDFF7", "#FFEB99", "#F7C4AE"][i]}"></span>
                        ${cat}：$${incomeValues[i]}
                    </li>
                `).join("")}
            `);
        })
        .catch(error => console.error("Error fetching income data:", error));

    // 從 API 獲取支出數據
    fetch("/api/expensechart")
        .then(response => response.json())
        .then(data => {
            const expenseCategories = Object.keys(data.categories);
            const expenseValues = Object.values(data.categories);

            // 渲染支出圓餅圖
            const expenseCtx = document.getElementById("expenseChart").getContext("2d");
            new Chart(expenseCtx, {
                type: "pie",
                data: {
                    labels: expenseCategories,
                    datasets: [
                        {
                            data: expenseValues,
                            backgroundColor: ["#FFCCCC", "#FFDD99", "#D5E8D4"],
                            borderColor: "#ffffff",
                            borderWidth: 2,
                        },
                    ],
                },
                options: {
                    plugins: {
                        legend: {
                            display: false,
                        },
                    },
                },
            });

            // 更新支出列表數據並添加顏色方塊
            $(".expense-details").html(`
                <li>總花費：${data.total}</li>
                ${expenseCategories.map((cat, i) => `
                    <li>
                        <span class="color-box" style="background-color: ${["#FFCCCC", "#FFDD99", "#D5E8D4"][i]}"></span>
                        ${cat}：$${expenseValues[i]}
                    </li>
                `).join("")}
            `);
        })
        .catch(error => console.error("Error fetching expense data:", error));
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
