import request from '@/utils/request'

// 用户列表
export function listUser(params) {
    return request({
        url: '/system/user',
        method: 'get',
        params
    })
}

// 获取用户详情
export function getUser(userId) {
    return request({
        url: `/system/user/${userId}`,
        method: 'get'
    })
}

// 创建用户
export function createUser(data) {
    return request({
        url: '/system/user',
        method: 'post',
        data
    })
}

// 更新用户
export function updateUser(userId, data) {
    return request({
        url: `/system/user/${userId}`,
        method: 'put',
        data
    })
}

// 删除用户
export function deleteUser(userId) {
    return request({
        url: `/system/user/${userId}`,
        method: 'delete'
    })
}

// 重置密码
export function resetUserPwd(userId, password) {
    return request({
        url: `/system/user/${userId}/resetPwd`,
        method: 'put',
        data: { password }
    })
}
