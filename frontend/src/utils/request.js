import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        const userStore = useUserStore()
        if (userStore.token) {
            config.headers['Authorization'] = `Bearer ${userStore.token}`
        }
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const res = response.data
        if (res.code !== 200) {
            ElMessage.error(res.msg || '请求失败')
            return Promise.reject(new Error(res.msg || 'Error'))
        }
        return res
    },
    error => {
        if (error.response?.status === 401) {
            ElMessage.error('登录已过期，请重新登录')
            const userStore = useUserStore()
            userStore.logout()
            router.push('/login')
        } else {
            ElMessage.error(error.message || '网络错误')
        }
        return Promise.reject(error)
    }
)

export default request
