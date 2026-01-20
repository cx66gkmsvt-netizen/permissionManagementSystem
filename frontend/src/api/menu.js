import request from '@/utils/request'

// 菜单树
export function listMenu() {
    return request({
        url: '/system/menu',
        method: 'get'
    })
}

// 获取所有菜单
export function listAllMenus() {
    return request({
        url: '/system/menu/all',
        method: 'get'
    })
}

// 获取菜单详情
export function getMenu(menuId) {
    return request({
        url: `/system/menu/${menuId}`,
        method: 'get'
    })
}

// 创建菜单
export function createMenu(data) {
    return request({
        url: '/system/menu',
        method: 'post',
        data
    })
}

// 更新菜单
export function updateMenu(menuId, data) {
    return request({
        url: `/system/menu/${menuId}`,
        method: 'put',
        data
    })
}

// 删除菜单
export function deleteMenu(menuId) {
    return request({
        url: `/system/menu/${menuId}`,
        method: 'delete'
    })
}
