let currentPage = 1
let pageSize = 5
let currentUserID = null
let isAdmin = false
let userList = []
let editingUserID = null
// 页面加载时获取数据
document.addEventListener("DOMContentLoaded", () => {
    getRoleAndAvatar()
})

//搜索功能
document.getElementById("keywordInput").addEventListener("keydown", e => {
    if (e.key === "Enter") {
        currentPage = 1
        loadUsers()
    }
})

//筛选功能
document.getElementById("statusFilter").addEventListener("change", () => {
    currentPage = 1
    loadUsers()
})


//请求当前头像和身份
async function getRoleAndAvatar() {
    try {
        const res = await fetch(`/api/users/raa`)
        const result = await res.json()
        console.log('后端返回的数据:', result);
        //这里是检验身份出错的状态码
        if (result.code === 401 || result.code === 403){
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
        currentUserID=result.data.id
        isAdmin = result.data.role === "管理员"

        loadUsers()
    } catch (err) {
        console.error(err)
        showModal("请求错误", "请求失败，请检查网络或联系管理员！", "danger")
    }
}


// 获取用户列表
async function loadUsers() {
    const keyword = document.getElementById("keywordInput").value
    const status = document.getElementById("statusFilter").value

    try {
        const res = await fetch(`/api/users?page=${currentPage}&size=${pageSize}&status=${status}&keyword=${keyword}`)
        const result = await res.json()
        console.log('后端返回的数据:', result);
        //这里是检验身份出错的状态码
        if (result.code === 401 || result.code === 403){
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
        userList=result.data.list
        renderTable(userList)
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

    const emailName = user.email.split('@')[0]
    const avatarName = emailName.length >= 2 ? emailName.slice(0, 2) : emailName;
    const avatar = user.avatar || `https://ui-avatars.com/api/?name=${encodeURIComponent(avatarName)}&background=random&color=fff`

    let editBtn = ""
    let deleteBtn = ""
    if (isAdmin) {
        // 管理员：可以编辑和删除
        editBtn = `
            <button onclick="editUser(${user.id})" class="p-2 hover:bg-gray-100 rounded">
                <i data-feather="edit" class="w-4 h-4 text-blue-600"></i>
            </button>`

        if (user.id === currentUserID) {
            // 自己 → 禁用按钮
            deleteBtn = `
            <button class="p-2 rounded cursor-not-allowed" disabled>
                <i data-feather="trash-2" class="w-4 h-4 text-gray-400"></i>
            </button>`;
        }else{
            deleteBtn = `
            <button onclick="deleteUser(${user.id})" class="p-2 hover:bg-gray-100 rounded">
                <i data-feather="trash-2" class="w-4 h-4 text-red-600"></i>
            </button>`
        }

    } else {
        // 普通用户：只能编辑自己
        if (user.id === currentUserID) {
            editBtn = `
                <button onclick="editUser(${user.id})" class="p-2 hover:bg-gray-100 rounded">
                    <i data-feather="edit" class="w-4 h-4 text-blue-600"></i>
                </button>`
        }
    }

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
                    ${editBtn}
                    ${deleteBtn}
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
    })

    const result = await res.json()

    //这里是检验身份出错的状态码
    if (result.code === 401 || result.code === 403){
        showModal("删除失败",result.message,"danger")
        const modalEl = document.getElementById('messageModal');

        modalEl.addEventListener('hidden.bs.modal', function () {
            window.location.href = "/login";
        });
        return
    }

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

// 编辑用户弹窗
function editUser(id) {
    openUserModal()
    editingUserID=id
    const user = userList.find(u => u.id === id)
    // 填充数据
    const emailName = user.email.split('@')[0]
    const avatarName = emailName.length >= 2 ? emailName.slice(0, 2) : emailName;
    document.getElementById("avatarPreview").src=user.avatar ||`https://ui-avatars.com/api/?name=${encodeURIComponent(avatarName)}&background=random&color=fff`
    document.getElementById("editUsername").value = user.username
    document.getElementById("userStatus").value = user.status
    document.getElementById("editPassword").value = ""
    const avatarInput = document.getElementById("avatarInput");
    avatarInput.value = ""; // 关键：清空文件选择

    const statusSelect = document.getElementById("userStatus")

    // 权限控制
    if (isAdmin) {
        statusSelect.disabled = user.id === currentUserID;
    } else {
        statusSelect.disabled = true;
    }

}


document.getElementById("avatarInput").addEventListener("change", function () {

    const file = this.files[0];
    if (!file) return;

    document.getElementById("avatarPreview").src = URL.createObjectURL(file);

});



async function saveUser(){

    // 获取表单数据
    const username = document.getElementById("editUsername").value;
    const password = document.getElementById("editPassword").value;
    const status = parseInt(document.getElementById("userStatus").value, 10);

    const avatarFile = document.getElementById("avatarInput").files[0];

    const formData = new FormData();

    formData.append("username", username);
    formData.append("password", password);
    formData.append("status", status);

    // 如果选择了新头像
    if (avatarFile) {
        formData.append("avatar", avatarFile);
    } else {
        formData.append("avatar", "");
    }

    const res = await fetch(`/api/users/${editingUserID}`, {
        method: "PUT",
        body: formData
    });

    const result = await res.json();

    // 身份验证失败
    if (result.code === 401 || result.code === 403){
        showModal("修改失败",result.message,"danger")

        const modalEl = document.getElementById('messageModal');

        modalEl.addEventListener('hidden.bs.modal', function () {
            window.location.href = "/login";
        });

        return
    }

    if (result.code === 200) {
        closeUserModal();
        getRoleAndAvatar()
    } else {
        showModal("修改失败",result.message,"danger")
    }
}