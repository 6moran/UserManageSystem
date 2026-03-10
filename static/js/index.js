// 页面加载时获取数据
document.addEventListener("DOMContentLoaded", () => {
    getRoleAndAvatar()
})



//请求当前头像
async function getRoleAndAvatar() {
    try {
        const res = await fetch(`/api/users/raa`)
        const result = await res.json()
        //这里是检验身份出错的状态码
        if (result.code === 401 || result.code === 403 || result.code === 500){
            showModal("请求错误",result.message,"danger")
            const modalEl = document.getElementById('messageModal');

            modalEl.addEventListener('hidden.bs.modal', function () {
                window.location.href = "/login";
            });
            return
        }
        //这里是业务
        if (result.code !== 200) {
            showModal("请求错误",result.message,"danger")
            return
        }
        const emailName = result.data.email.split('@')[0]
        const avatarName = emailName.length >= 2 ? emailName.slice(0, 2) : emailName;
        document.getElementById("userAvatar").src = result.data.avatar || `https://ui-avatars.com/api/?name=${encodeURIComponent(avatarName)}&background=random&color=fff`
        setPageTitle('数据概览');
        // 初始化图标
        feather.replace();
    } catch (err) {
        console.error(err)
        showModal("请求错误", "请求失败，请检查网络或联系管理员！", "danger")
    }
}



// 全局功能
function setPageTitle(title) {
    document.getElementById('pageTitle').textContent = title;
    document.title = title + ' - 后台管理系统';
}

function logout() {
    if (confirm('确定要退出系统吗？')) {
        // 清理前端存储
        localStorage.clear();

        // 发送请求给后端销毁 cookie
        fetch('/logout', {
            method: 'POST',            // 后端 logout 接口
        }).finally(() => {
            // 请求完成后跳转登录页
            window.location.href = '/login';
        });
    }
}

