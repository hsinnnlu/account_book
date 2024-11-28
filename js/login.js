$(document).ready(function() {
    const registerForm = $('#registerForm');
    const successModal = new bootstrap.Modal($('#successModal')[0]);
    const successMessage = $('#success_message');
    const errorMessage = $('#error_message');

    registerForm.on('submit', function(event) {
        event.preventDefault(); // 防止表單的默認提交行為

        const formData = registerForm.serialize(); // 序列化表單數據

        // 清空錯誤消息
        errorMessage.hide().text('');

        // 使用 Fetch API 提交表單
        fetch('/register', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            // 檢查響應狀態
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.message || '未知錯誤');
                });
            }
            return response.json();
        })
        .then(data => {
            const message = data.message || '註冊成功';
            successMessage.text(message); // 更新模態中的消息內容

            // 隱藏註冊模態
            const registerModal = bootstrap.Modal.getInstance($('#registerModal')[0]);
            if (registerModal) {
                registerModal.hide(); // 隱藏註冊模態
            }

            // 顯示成功模態
            successModal.show();
        })
        .catch(error => {
            console.error('錯誤:', error);
            errorMessage.text(error.message).show(); // 顯示錯誤消息

            // 隱藏註冊模態
            const registerModal = bootstrap.Modal.getInstance($('#registerModal')[0]);
            if (registerModal) {
                registerModal.hide(); // 隱藏註冊模態
            }
        });
    });
});

