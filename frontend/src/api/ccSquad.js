import request from '@/utils/request'

// ==================== CC战队管理 API ====================

// 战队列表
export function listCCSquad(params) {
    return request({
        url: '/cc/squad',
        method: 'get',
        params
    })
}

// 获取全部战队
export function listAllCCSquad() {
    return request({
        url: '/cc/squad/all',
        method: 'get'
    })
}

// 根据团队获取战队列表
export function listCCSquadByTeam(teamId) {
    return request({
        url: '/cc/squad/byTeam',
        method: 'get',
        params: { teamId }
    })
}

// 获取战队详情
export function getCCSquad(id) {
    return request({
        url: `/cc/squad/${id}`,
        method: 'get'
    })
}

// 创建战队
export function createCCSquad(data) {
    return request({
        url: '/cc/squad',
        method: 'post',
        data
    })
}

// 更新战队
export function updateCCSquad(id, data) {
    return request({
        url: `/cc/squad/${id}`,
        method: 'put',
        data
    })
}

// 获取战队修改记录
export function getCCSquadLogs(id) {
    return request({
        url: `/cc/squad/${id}/logs`,
        method: 'get'
    })
}

// 获取战队资金信息
export function getCCSquadFund(id) {
    return request({
        url: `/cc/squad/${id}/fund`,
        method: 'get'
    })
}

// 编辑战队余额
export function editCCSquadFund(id, data) {
    return request({
        url: `/cc/squad/${id}/fund`,
        method: 'put',
        data
    })
}

// 战队充值
export function rechargeCCSquad(id, data) {
    return request({
        url: `/cc/squad/${id}/recharge`,
        method: 'post',
        data
    })
}

// 战队转账
export function transferCCSquad(id, data) {
    return request({
        url: `/cc/squad/${id}/transfer`,
        method: 'post',
        data
    })
}

// 获取战队账单明细
export function getCCSquadBills(id, billType) {
    return request({
        url: `/cc/squad/${id}/bills`,
        method: 'get',
        params: { billType }
    })
}
