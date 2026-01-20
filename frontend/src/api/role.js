import request from '@/utils/request'

// 角色列表
export function listRole(params) {
    return request({
        url: '/system/role',
        method: 'get',
        params
    })
}

// 获取所有角色
export function listAllRoles() {
    return request({
        url: '/system/role/all',
        method: 'get'
    })
}

// 获取角色详情
export function getRole(roleId) {
    return request({
        url: `/system/role/${roleId}`,
        method: 'get'
    })
}

// 创建角色
export function createRole(data) {
    return request({
        url: '/system/role',
        method: 'post',
        data
    })
}

// 更新角色
export function updateRole(roleId, data) {
    return request({
        url: `/system/role/${roleId}`,
        method: 'put',
        data
    })
}

// 删除角色
export function deleteRole(roleId) {
    return request({
        url: `/system/role/${roleId}`,
        method: 'delete'
    })
}
