import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, logout as logoutApi, getUserInfo, getRoutes } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
    const token = ref(localStorage.getItem('token') || '')
    const userInfo = ref(null)
    const roles = ref([])
    const permissions = ref([])
    const menuRoutes = ref([])

    const isLoggedIn = computed(() => !!token.value)

    // 登录
    async function login(loginForm) {
        const res = await loginApi(loginForm)
        token.value = res.data.token
        localStorage.setItem('token', res.data.token)
        return res
    }

    // 获取用户信息
    async function fetchUserInfo() {
        const res = await getUserInfo()
        userInfo.value = res.data.user
        roles.value = res.data.roles || []
        permissions.value = res.data.permissions || []
        return res.data
    }

    // 获取路由
    async function fetchRoutes() {
        const res = await getRoutes()
        menuRoutes.value = res.data || []
        return res.data
    }

    // 登出
    function logout() {
        token.value = ''
        userInfo.value = null
        roles.value = []
        permissions.value = []
        menuRoutes.value = []
        localStorage.removeItem('token')
    }

    // 检查权限
    function hasPermission(perm) {
        if (permissions.value.includes('*:*:*')) return true
        return permissions.value.includes(perm)
    }

    return {
        token,
        userInfo,
        roles,
        permissions,
        menuRoutes,
        isLoggedIn,
        login,
        fetchUserInfo,
        fetchRoutes,
        logout,
        hasPermission
    }
})
