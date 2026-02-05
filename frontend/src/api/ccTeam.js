import request from '@/utils/request'

// ==================== CC团队管理 API ====================

// 团队列表
export function listCCTeam(params) {
    return request({
        url: '/cc/team',
        method: 'get',
        params
    })
}

// 获取全部团队
export function listAllCCTeam() {
    return request({
        url: '/cc/team/all',
        method: 'get'
    })
}

// 获取团队详情
export function getCCTeam(id) {
    return request({
        url: `/cc/team/${id}`,
        method: 'get'
    })
}

// 创建团队
export function createCCTeam(data) {
    return request({
        url: '/cc/team',
        method: 'post',
        data
    })
}

// 更新团队
export function updateCCTeam(id, data) {
    return request({
        url: `/cc/team/${id}`,
        method: 'put',
        data
    })
}

// 获取团队修改记录
export function getCCTeamLogs(id) {
    return request({
        url: `/cc/team/${id}/logs`,
        method: 'get'
    })
}

// 获取团队资金信息
export function getCCTeamFund(id) {
    return request({
        url: `/cc/team/${id}/fund`,
        method: 'get'
    })
}

// 编辑团队余额
export function editCCTeamFund(id, data) {
    return request({
        url: `/cc/team/${id}/fund`,
        method: 'put',
        data
    })
}

// 团队充值
export function rechargeCCTeam(id, data) {
    return request({
        url: `/cc/team/${id}/recharge`,
        method: 'post',
        data
    })
}

// 团队转账
export function transferCCTeam(id, data) {
    return request({
        url: `/cc/team/${id}/transfer`,
        method: 'post',
        data
    })
}

// 获取团队账单明细
export function getCCTeamBills(id, billType) {
    return request({
        url: `/cc/team/${id}/bills`,
        method: 'get',
        params: { billType }
    })
}
