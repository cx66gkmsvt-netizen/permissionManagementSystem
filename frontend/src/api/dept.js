import request from '@/utils/request'

// 部门树
export function listDept() {
    return request({
        url: '/system/dept',
        method: 'get'
    })
}

// 获取所有部门
export function listAllDepts() {
    return request({
        url: '/system/dept/all',
        method: 'get'
    })
}

// 获取部门详情
export function getDept(deptId) {
    return request({
        url: `/system/dept/${deptId}`,
        method: 'get'
    })
}

// 创建部门
export function createDept(data) {
    return request({
        url: '/system/dept',
        method: 'post',
        data
    })
}

// 更新部门
export function updateDept(deptId, data) {
    return request({
        url: `/system/dept/${deptId}`,
        method: 'put',
        data
    })
}

// 删除部门
export function deleteDept(deptId) {
    return request({
        url: `/system/dept/${deptId}`,
        method: 'delete'
    })
}
