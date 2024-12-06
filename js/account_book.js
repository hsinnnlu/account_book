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