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
            },
            {
                path: 'system/cc',
                name: 'CC',
                component: () => import('@/views/cc/index.vue'),
                meta: { title: 'CC管理', permission: 'system:cc:list' }
            },
            {
                path: 'cc/legion',
                name: 'Legion',
                component: () => import('@/views/cc/legion/index.vue'),
                meta: { title: '军团管理', permission: 'cc:legion:list' }
            },
            {
                path: 'cc/team',
                name: 'CCTeam',
                component: () => import('@/views/cc/team/index.vue'),
                meta: { title: 'CC团队管理', permission: 'cc:team:list' }
            },
            {
                path: 'cc/squad',
                name: 'CCSquad',
                component: () => import('@/views/cc/squad/index.vue'),
                meta: { title: 'CC战队管理', permission: 'cc:squad:list' }
            },
            {
                path: 'cc/attendance',
                name: 'Attendance',
                component: () => import('@/views/cc/attendance/index.vue'),
                meta: { title: '在班管理', permission: 'cc:attendance:list' }
            },
            {
                path: 'cc/lead-allocation',
                name: 'LeadAllocation',
                component: () => import('@/views/cc/leadAllocation/index.vue'),
                meta: { title: '例子分配', permission: 'cc:lead-allocation:list' }
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

// 无需权限检查的页面
const noPermPages = ['/dashboard', '/profile']

// 获取第一个可访问的菜单路径
function getFirstMenuPath(routes) {
    for (const menu of routes) {
        if (menu.children && menu.children.length > 0) {
            const childPath = menu.children[0].path
            return childPath.startsWith('/') ? childPath : menu.path + '/' + childPath
        }
        if (menu.path) {
            return menu.path
        }
    }
    return '/dashboard'
}

// 检查用户是否有该路由权限
function hasRoutePermission(to, userStore) {
    // 无需权限的页面直接通过
    if (noPermPages.includes(to.path)) return true
    // 没有定义permission的路由直接通过
    if (!to.meta.permission) return true
    // 超管权限直接通过
    if (userStore.permissions.includes('*:*:*')) return true
    // 检查是否有该权限
    return userStore.permissions.includes(to.meta.permission)
}

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
                    // 检查是否有访问目标页面的权限
                    if (!hasRoutePermission(to, userStore)) {
                        const firstPath = getFirstMenuPath(userStore.menuRoutes)
                        next({ path: firstPath, replace: true })
                    } else {
                        next({ ...to, replace: true })
                    }
                } catch (error) {
                    userStore.logout()
                    next('/login')
                }
            } else {
                // 检查权限，无权限则重定向
                if (!hasRoutePermission(to, userStore)) {
                    const firstPath = getFirstMenuPath(userStore.menuRoutes)
                    next({ path: firstPath, replace: true })
                } else {
                    next()
                }
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
