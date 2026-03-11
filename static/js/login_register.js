document.addEventListener("DOMContentLoaded", function() {
  function setup(btnId, inputId) {
    const btn = document.getElementById(btnId)
    const input = document.getElementById(inputId)
    if (!btn || !input || !window.feather || !feather.icons) return
    btn.innerHTML = feather.icons["eye"].toSvg({ width: 20, height: 20 })
    btn.addEventListener("mousedown", function(e) { e.preventDefault() })
    function update() {
      const hasFocus = document.activeElement === input
      const hasText = input.value.length > 0
      const show = hasFocus && hasText
      if (show) {
        btn.classList.remove("opacity-0","pointer-events-none")
        btn.classList.add("opacity-100")
      } else {
        btn.classList.add("opacity-0","pointer-events-none")
        btn.classList.remove("opacity-100")
      }
    }
    btn.addEventListener("click", function(e) {
      e.preventDefault()
      const start = typeof input.selectionStart === "number" ? input.selectionStart : input.value.length
      const end = typeof input.selectionEnd === "number" ? input.selectionEnd : input.value.length
      const isHidden = input.type === "password"
      input.type = isHidden ? "text" : "password"
      btn.innerHTML = isHidden
        ? feather.icons["eye-off"].toSvg({ width: 20, height: 20 })
        : feather.icons["eye"].toSvg({ width: 20, height: 20 })
      input.focus()
      if (typeof input.setSelectionRange === "function") {
        input.setSelectionRange(start, end)
      }
      update()
    })
    input.addEventListener("focus", update)
    input.addEventListener("blur", update)
    input.addEventListener("input", update)
    update()
  }
  setup("togglePwd","password")
  setup("toggleRegPwd","password")
  setup("toggleRegConfirm","confirmPassword")
})
