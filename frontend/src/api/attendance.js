import request from '@/utils/request'

// ==================== 在班管理 API ====================

// 获取在班列表
export function getAttendanceList(params) {
    return request({
        url: '/cc/attendance',
        method: 'get',
        params,
        paramsSerializer: params => {
            const searchParams = new URLSearchParams()
            for (const key in params) {
                const val = params[key]
                if (Array.isArray(val)) {
                    val.forEach(v => searchParams.append(key, v))
                } else if (val !== null && val !== undefined) {
                    searchParams.append(key, val)
                }
            }
            return searchParams.toString()
        }
    })
}

// 更新在班状态
export function updateAttendance(ccId, date, data) {
    return request({
        url: `/cc/attendance/${ccId}/${date}`,
        method: 'put',
        data
    })
}

// 批量更新在班状态
export function batchUpdateAttendance(data) {
    return request({
        url: '/cc/attendance/batch',
        method: 'post',
        data
    })
}

// 获取在班统计
export function getAttendanceStats(date) {
    return request({
        url: '/cc/attendance/stats',
        method: 'get',
        params: { date }
    })
}

// 获取CC在班历史
export function getCCAttendanceHistory(ccId, startDate, endDate) {
    return request({
        url: `/cc/attendance/${ccId}/history`,
        method: 'get',
        params: { startDate, endDate }
    })
}

// 导出在班记录
export function exportAttendance(params) {
    return request({
        url: '/cc/attendance/export',
        method: 'get',
        params,
        responseType: 'blob',
        paramsSerializer: params => {
            const searchParams = new URLSearchParams()
            for (const key in params) {
                const val = params[key]
                if (Array.isArray(val)) {
                    val.forEach(v => searchParams.append(key, v))
                } else if (val !== null && val !== undefined) {
                    searchParams.append(key, val)
                }
            }
            return searchParams.toString()
        }
    })
}
