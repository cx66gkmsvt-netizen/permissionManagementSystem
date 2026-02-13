/**
 * 全局按钮防连点
 * 在 document 层拦截所有 button / el-button 的 click 事件
 * 点击后在指定时间内禁止再次点击（默认 1500ms）
 */
const THROTTLE_TIME = 1500

export function setupThrottleClick() {
    document.addEventListener('click', (e) => {
        // 向上查找最近的 button 元素（含 el-button）
        const btn = e.target.closest('button, .el-button')
        if (!btn) return

        // 跳过已禁用的按钮
        if (btn.disabled || btn.classList.contains('is-disabled')) return

        // 检查是否在冷却中
        if (btn.__throttled) {
            e.stopImmediatePropagation()
            e.preventDefault()
            return
        }

        // 设置冷却
        btn.__throttled = true
        btn.classList.add('is-disabled')
        btn.style.pointerEvents = 'none'

        setTimeout(() => {
            btn.__throttled = false
            btn.classList.remove('is-disabled')
            btn.style.pointerEvents = ''
        }, THROTTLE_TIME)
    }, true) // 使用捕获阶段，优先于组件事件
}
