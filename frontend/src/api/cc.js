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

// 获取CC资金信息
export function getCCFund(id) {
    return request({
        url: `/system/cc/${id}/fund`,
        method: 'get'
    })
}

// 编辑CC余额
export function editCCFund(id, data) {
    return request({
        url: `/system/cc/${id}/fund`,
        method: 'put',
        data
    })
}

// CC转账
export function transferCC(id, data) {
    return request({
        url: `/system/cc/${id}/transfer`,
        method: 'post',
        data
    })
}

// 获取CC账单
export function getCCBills(id, billType) {
    return request({
        url: `/system/cc/${id}/bills`,
        params: { billType },
        method: 'get'
    })
}
