//使用bootstrap来优化弹窗
function showModal(title, message, type = "primary") {
    // 设置标题和内容
    document.getElementById("modalTitle").innerText = title;
    document.getElementById("modalBody").innerText = message;

    // 修改按钮颜色
    const btn = document.getElementById("modalConfirmBtn");
    if(btn){
        btn.className = "btn btn-" + type;
    }


    // 创建 Modal 实例并显示
    const modal = new bootstrap.Modal(document.getElementById('messageModal'));
    modal.show();
}