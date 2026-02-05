import request from '@/utils/request'

// ==================== 例子分配 API ====================

// 获取分配列表
export function getLeadAllocationList(params) {
    return request({
        url: '/cc/lead-allocation',
        method: 'get',
        params
    })
}

// 更新分配记录
export function updateLeadAllocation(ccId, date, data) {
    return request({
        url: `/cc/lead-allocation/${ccId}`,
        method: 'put',
        data,
        params: { date }
    })
}

// 批量更新是否分配
export function batchUpdateIsAllocated(data) {
    return request({
        url: '/cc/lead-allocation/batch',
        method: 'post',
        data
    })
}

// 获取分配统计
export function getLeadAllocationStats(date) {
    return request({
        url: '/cc/lead-allocation/stats',
        method: 'get',
        params: { date }
    })
}

// 获取分配详情历史
export function getLeadAllocationDetail(ccId, params) {
    return request({
        url: `/cc/lead-allocation/${ccId}/detail`,
        method: 'get',
        params
    })
}
