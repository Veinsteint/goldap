import request from '@/utils/request'

// 获取分组用户权限列表
export function getGroupUserPermissions(params) {
  return request({
    url: '/api/group-user-permission',
    method: 'get',
    params
  })
}

// 添加分组用户权限
export function addGroupUserPermission(data) {
  return request({
    url: '/api/group-user-permission',
    method: 'post',
    data
  })
}

// 更新分组用户权限
export function updateGroupUserPermission(data) {
  return request({
    url: '/api/group-user-permission',
    method: 'put',
    data
  })
}

// 删除分组用户权限
export function deleteGroupUserPermission(data) {
  return request({
    url: '/api/group-user-permission',
    method: 'delete',
    data
  })
}

