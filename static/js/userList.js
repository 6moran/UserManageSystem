
let currentPage = 1
let pageSize = 10

// 页面加载时获取数据
document.addEventListener("DOMContentLoaded", () => {
    loadUsers()
})

// 获取用户列表
async function loadUsers() {
    const keyword = document.querySelector('input[placeholder="搜索用户..."]').value
    const status = document.getElementById("statusFilter").value

    try {
        const res = await fetch(`/api/users?page=${currentPage}&size=${pageSize}&status=${status}&keyword=${keyword}`)
        const result = await res.json()
        console.log('后端返回的数据:', result);
        //这里是检验身份出错的状态码
        if (result.code === 401 || result.code === 403 || result.code === 500){
            showModal("请求错误",result.message,"danger")
            window.location.href = "/login"
            return
        }

        //这里是业务
        if (result.code !== 200) {
            showModal("请求错误",result.message,"danger")
            return
        }

        //业务通过
        renderTable(result.data.list)
        renderPagination(result.data.total)

    } catch (err) {
        console.error(err)
        showModal("请求错误", "请求失败，请检查网络或联系管理员！", "danger")
    }
}

// 渲染表格
function renderTable(list) {
    const tbody = document.getElementById("userTableBody")
    tbody.innerHTML = ""

    list.forEach(user => {
    const statusHTML = user.status === 1
        ? `<span class="px-2 py-1 bg-green-100 text-green-800 rounded-full text-sm">启用</span>`
        : `<span class="px-2 py-1 bg-red-100 text-red-800 rounded-full text-sm">禁用</span>`

    const avatar = user.avatar || "https://via.placeholder.com/40"

    const row = `
        <tr class="hover:bg-gray-50">
            <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                    <img src="${avatar}" class="w-8 h-8 rounded-full">
                    <div>
                        <p class="font-medium">${user.username}</p>
                        <p class="text-sm text-gray-500">${user.email}</p>
                    </div>
                </div>
            </td>
            <td class="px-6 py-4">
                <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                    ${user.role}
                </span>
            </td>
            <td class="px-6 py-4">${user.lastLogin || '-'}</td>
            <td class="px-6 py-4">${statusHTML}</td>
            <td class="px-6 py-4">
                <div class="flex gap-2">
                    <button onclick="editUser(${user.id})"
                        class="p-2 hover:bg-gray-100 rounded">
                        <i data-feather="edit" class="w-4 h-4 text-blue-600"></i>
                    </button>
                    <button onclick="deleteUser(${user.id})"
                        class="p-2 hover:bg-gray-100 rounded">
                        <i data-feather="trash-2" class="w-4 h-4 text-red-600"></i>
                    </button>
                </div>
            </td>
        </tr>
        `
    tbody.innerHTML += row
    })

    feather.replace()
}

// 渲染分页
function renderPagination(total) {
    const totalPages = Math.ceil(total / pageSize)
    const container = document.querySelector(".flex.gap-2")
    container.innerHTML = ""

    // 上一页
    const prevBtn = document.createElement("button")
    prevBtn.textContent = "上一页"
    prevBtn.className = "px-3 py-1 rounded hover:bg-gray-100"
    prevBtn.disabled = currentPage === 1
    prevBtn.onclick = () => {
        currentPage--
        loadUsers()
    }
    container.appendChild(prevBtn)

    // 页码
    for (let i = 1; i <= totalPages; i++) {
        const btn = document.createElement("button")
        btn.textContent = i
        btn.className = i === currentPage
            ? "px-3 py-1 bg-blue-100 rounded"
            : "px-3 py-1 rounded hover:bg-gray-100"

        btn.onclick = () => {
            currentPage = i
            loadUsers()
        }
        container.appendChild(btn)
    }

    // 下一页
    const nextBtn = document.createElement("button")
    nextBtn.textContent = "下一页"
    nextBtn.className = "px-3 py-1 rounded hover:bg-gray-100"
    nextBtn.disabled = currentPage === totalPages
    nextBtn.onclick = () => {
        currentPage++
        loadUsers()
    }
    container.appendChild(nextBtn)
}


    // 用户管理页面专属逻辑
    function openUserModal() {
    document.getElementById('userModal').classList.remove('hidden');
}

    function closeUserModal() {
    document.getElementById('userModal').classList.add('hidden');
}

    // 点击模态框外部关闭
    document.getElementById('userModal').addEventListener('click', (e) => {
    if (e.target === document.getElementById('userModal')) {
    closeUserModal();
}
});

    function deleteUser() {
    if (confirm('确定要删除该用户吗？')) {
    // 删除用户逻辑
}
}



    // 全局功能
    function setPageTitle(title) {
    document.getElementById('pageTitle').textContent = title;
    document.title = title + ' - 后台管理系统';
}

    function logout() {
    if (confirm('确定要退出系统吗？')) {
    localStorage.clear();
    window.location.href = '/login';
}
}

    // 初始化图标
    feather.replace();



// // 删除用户
// async function deleteUser(id) {
//     if (!confirm("确定要删除该用户吗？")) return
//
//     const res = await fetch(`/api/users/${id}`, {
//         method: "DELETE",
//         credentials: "include"
//     })
//
//     const result = await res.json()
//
//     if (result.code === 200) {
//         alert("删除成功")
//         loadUsers()
//     } else {
//         alert(result.message)
//     }
// }
//
// // 编辑用户
// function editUser(id) {
//     openUserModal()
//     // 这里可以再请求详情接口
// }
