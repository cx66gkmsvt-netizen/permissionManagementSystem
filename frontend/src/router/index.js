import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 静态路由
const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/index.vue'),
        meta: { title: '登录' }
    },
    {
        path: '/',
        component: () => import('@/views/layout/index.vue'),
        redirect: '/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/dashboard/index.vue'),
                meta: { title: '首页' }
            },
            {
                path: 'profile',
                name: 'Profile',
                component: () => import('@/views/profile/index.vue'),
                meta: { title: '个人中心' }
            },
            {
                path: 'system/user',
                name: 'User',
                component: () => import('@/views/system/user/index.vue'),
                meta: { title: '用户管理', permission: 'system:user:list' }
            },
            {
                path: 'system/role',
                name: 'Role',
                component: () => import('@/views/system/role/index.vue'),
                meta: { title: '角色管理', permission: 'system:role:list' }
            },
            {
                path: 'system/menu',
                name: 'Menu',
                component: () => import('@/views/system/menu/index.vue'),
                meta: { title: '菜单管理', permission: 'system:menu:list' }
            },
            {
                path: 'system/dept',
                name: 'Dept',
                component: () => import('@/views/system/dept/index.vue'),
                meta: { title: '部门管理', permission: 'system:dept:list' }
            }
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/error/404.vue'),
        meta: { title: '404' }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 白名单
const whiteList = ['/login']

// 路由守卫
router.beforeEach(async (to, from, next) => {
    document.title = to.meta.title ? `${to.meta.title} - 用户中心` : '用户中心'

    const userStore = useUserStore()

    if (userStore.token) {
        if (to.path === '/login') {
            next({ path: '/' })
        } else {
            // 如果没有用户信息，获取用户信息
            if (!userStore.userInfo) {
                try {
                    await userStore.fetchUserInfo()
                    await userStore.fetchRoutes()
                    next({ ...to, replace: true })
                } catch (error) {
                    userStore.logout()
                    next('/login')
                }
            } else {
                next()
            }
        }
    } else {
        if (whiteList.includes(to.path)) {
            next()
        } else {
            next('/login')
        }
    }
})

export default router
