document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById("loginForm");

    form.addEventListener("submit", async function(e) {
        e.preventDefault(); // 阻止默认提交

        // 获取表单数据
        const email = document.getElementById("email").value.trim();
        const password = document.getElementById("password").value;

        // 构造请求数据
        const data = {
            email: email,
            password: password
        };

        try {
            const response = await fetch("/api/login", { // 后端接口地址
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            const result = await response.json(); // 解析 JSON

            // 根据后端返回处理
            if (result.code === 200) {
                //直接跳转
                window.location.href = "/index";
            } else {
                // 后端返回错误
                showModal("登录失败",result.message,"danger")
            }
        } catch (err) {
            console.error("登录请求出错:", err);
            showModal("登录失败","请求失败，请检查网络或联系管理员！","danger")
        }
    });
});