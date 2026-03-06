let currentPage = 1
let pageSize = 5

// 页面加载时获取数据
document.addEventListener("DOMContentLoaded", () => {
    loadUsers()
})


document.getElementById("keywordInput").addEventListener("keydown", e => {
    if (e.key === "Enter") {
        currentPage = 1
        loadUsers()
    }
})


document.getElementById("statusFilter").addEventListener("change", () => {
    currentPage = 1
    loadUsers()
})

// 获取用户列表
async function loadUsers() {
    const keyword = document.getElementById("keywordInput").value
    const status = document.getElementById("statusFilter").value

    try {
        const res = await fetch(`/api/users?page=${currentPage}&size=${pageSize}&status=${status}&keyword=${keyword}`)
        const result = await res.json()
        console.log('后端返回的数据:', result);
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

        //业务通过
        renderTable(result.data.list)
        renderPagination(result.data.total)
        setPageTitle('用户管理');

    } catch (err) {
        console.error(err)
        showModal("请求错误", "请求失败，请检查网络或联系管理员！", "danger")
    }
}

// 渲染表格
function renderTable(list) {
    const tbody = document.getElementById("userTableBody")
    tbody.innerHTML = ""

    if (!Array.isArray(list) || list.length === 0) {
        // 清空表格
        document.getElementById("userTableBody").innerHTML = `
            <tr>
                <td colspan="6" class="text-center py-4 text-gray-500">暂无数据</td>
            </tr>
        `;
        return;
    }


    list.forEach(user => {
    const statusHTML = user.status === 1
        ? `<span class="px-2 py-1 bg-green-100 text-green-800 rounded-full text-sm">启用</span>`
        : `<span class="px-2 py-1 bg-red-100 text-red-800 rounded-full text-sm">禁用</span>`

    const roleHTML = user.role === "管理员"
        ? `<span class="px-2 py-1 bg-orange-100 text-orange-800 rounded-full text-sm">管理员</span>`
        : `<span class="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">用户</span>`

    const avatar = user.avatar || `https://ui-avatars.com/api/?name=${encodeURIComponent(user.username)}&length=2&background=random&color=fff`

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
            <td class="px-6 py-4">${roleHTML}</td>
            <td class="px-6 py-4 ">${user.create_time ? dayjs.utc(user.create_time).format('YYYY-MM-DD HH:mm:ss') : '-'}</td>
            <td class="px-6 py-4 ">${user.last_time ? dayjs.utc(user.last_time).format('YYYY-MM-DD HH:mm:ss') : '-'}</td>
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
    // 更新总条数
    document.getElementById("totalCount").innerText = total
    const container = document.getElementById('paginationButtons')
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

// 删除用户
async function deleteUser(id) {
    if (!confirm("确定要删除该用户吗？")) return

    const res = await fetch(`/api/users/${id}`, {
        method: "DELETE",
        credentials: "include"
    })

    const result = await res.json()

    if (result.code === 200) {
        showModal("删除成功",result.message,"success")
        const modalEl = document.getElementById('messageModal');

        modalEl.addEventListener('hidden.bs.modal', function () {
            loadUsers()
        });
    } else {
        showModal("删除失败",result.message,"danger")
    }
}

// 编辑用户
function editUser(id) {
    openUserModal()
}
