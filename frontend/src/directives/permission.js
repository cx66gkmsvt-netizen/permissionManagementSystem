import { useUserStore } from '@/stores/user'

export const permission = {
    mounted(el, binding) {
        const userStore = useUserStore()
        const { value } = binding

        if (value) {
            const permissions = userStore.permissions
            const hasPermission = permissions.includes('*:*:*') || permissions.includes(value)

            if (!hasPermission) {
                el.parentNode?.removeChild(el)
            }
        }
    }
}
