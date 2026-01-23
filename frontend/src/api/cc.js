import request from '@/utils/request'

// CC列表
export function listCC(params) {
    return request({
        url: '/system/cc',
        method: 'get',
        params
    })
}

// 获取CC详情
export function getCC(id) {
    return request({
        url: `/system/cc/${id}`,
        method: 'get'
    })
}

// 创建CC
export function createCC(data) {
    return request({
        url: '/system/cc',
        method: 'post',
        data
    })
}

// 更新CC
export function updateCC(id, data) {
    return request({
        url: `/system/cc/${id}`,
        method: 'put',
        data
    })
}

// 删除CC
export function deleteCC(id) {
    return request({
        url: `/system/cc/${id}`,
        method: 'delete'
    })
}
