
document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById("registerForm");

    form.addEventListener("submit", async function(e) {
        e.preventDefault(); // 阻止默认提交

        // 获取表单数据
        const email = document.getElementById("email").value.trim();
        const password = document.getElementById("password").value;
        const confirmPassword = document.getElementById("confirmPassword").value;

        // 简单前端验证
        if (password !== confirmPassword) {
            showModal("注册失败","两次输入的密码不一致","danger")
            return;
        }

        // 构造请求数据
        const data = {
            email: email,
            password: password
        };

        try {
            const response = await fetch("/api/register", { // 后端接口地址
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            const result = await response.json(); // 解析 JSON

            // 根据后端返回处理
            if (result.code === 200) {
                showModal("注册成功",result.message,"success")
                const modalEl = document.getElementById('messageModal');

                modalEl.addEventListener('hidden.bs.modal', function () {
                    window.location.href = "/login";
                }); // 注册成功确认后跳转登录
            } else {
                // 后端返回错误
                showModal("注册失败",result.message,"danger")
            }
        } catch (err) {
            console.error("注册请求出错:", err);
            showModal("注册失败","请求失败，请检查网络或联系管理员！","danger")
        }
    });
});