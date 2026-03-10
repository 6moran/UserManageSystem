function initPasswordToggle(pwdId, toggleId) {
    const toggle = document.getElementById(toggleId)
    const pwd = document.getElementById(pwdId)
    let toggling = false

    if (!toggle || !pwd || !window.feather || !feather.icons) return

    // 初始化图标
    toggle.innerHTML = feather.icons["eye"].toSvg({ width: 20, height: 20 })

    // 防止点击切换时触发 blur 导致图标瞬间隐藏
    toggle.addEventListener("mousedown", (e) => {
        toggling = true
        e.preventDefault()
    })

    toggle.addEventListener("click", (e) => {
        e.preventDefault()
        if (pwd.type === "password") {
            pwd.type = "text"
            toggle.innerHTML = feather.icons["eye-off"].toSvg({ width: 20, height: 20 })
        } else {
            pwd.type = "password"
            toggle.innerHTML = feather.icons["eye"].toSvg({ width: 20, height: 20 })
        }
        // 点击后保持输入框聚焦以继续显示图标
        pwd.focus()
        toggling = false
        updateToggleVisibility()
    })

    // 更新图标显示状态
    function updateToggleVisibility() {
        const hasFocus = document.activeElement === pwd
        const hasText = pwd.value.length > 0
        const show = (hasFocus || toggling) && hasText

        if (show) {
            toggle.classList.remove("opacity-0", "pointer-events-none", "opacity-70")
            toggle.classList.add("opacity-100")
        } else {
            toggle.classList.add("opacity-0", "pointer-events-none")
            toggle.classList.remove("opacity-100")
        }
    }

    // 输入框事件绑定
    pwd.addEventListener("focus", updateToggleVisibility)
    pwd.addEventListener("blur", () => {
        if (toggling) return
        updateToggleVisibility()
    })
    pwd.addEventListener("input", updateToggleVisibility)

    // 初始化状态
    updateToggleVisibility()
}