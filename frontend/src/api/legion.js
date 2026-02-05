import request from '@/utils/request'

// ==================== 军团管理 API ====================

// 军团列表
export function listLegion(params) {
    return request({
        url: '/cc/legion',
        method: 'get',
        params
    })
}

// 获取全部军团
export function listAllLegion() {
    return request({
        url: '/cc/legion/all',
        method: 'get'
    })
}

// 获取军团详情
export function getLegion(id) {
    return request({
        url: `/cc/legion/${id}`,
        method: 'get'
    })
}

// 创建军团
export function createLegion(data) {
    return request({
        url: '/cc/legion',
        method: 'post',
        data
    })
}

// 更新军团
export function updateLegion(id, data) {
    return request({
        url: `/cc/legion/${id}`,
        method: 'put',
        data
    })
}

// 获取军团跟进记录
export function getLegionLogs(id) {
    return request({
        url: `/cc/legion/${id}/logs`,
        method: 'get'
    })
}

// 获取军团资金信息
export function getLegionFund(id) {
    return request({
        url: `/cc/legion/${id}/fund`,
        method: 'get'
    })
}

// 编辑军团余额
export function editLegionFund(id, data) {
    return request({
        url: `/cc/legion/${id}/fund`,
        method: 'put',
        data
    })
}

// 军团充值
export function rechargeLegion(id, data) {
    return request({
        url: `/cc/legion/${id}/recharge`,
        method: 'post',
        data
    })
}

// 军团转账
export function transferLegion(id, data) {
    return request({
        url: `/cc/legion/${id}/transfer`,
        method: 'post',
        data
    })
}

// 获取军团账单明细
export function getLegionBills(id, billType) {
    return request({
        url: `/cc/legion/${id}/bills`,
        method: 'get',
        params: { billType }
    })
}
