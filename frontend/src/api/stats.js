import request from '@/utils/request'

// 获取系统统计数据
export function getStats() {
    return request({
        url: '/system/stats',
        method: 'get'
    })
}
