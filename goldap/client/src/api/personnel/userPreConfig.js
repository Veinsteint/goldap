import request from '@/utils/request'

// 用户配置列表
export function preConfigList(params) {
  return request({
    url: '/api/user/preconfig/list',
    method: 'get',
    params
  })
}

// 添加用户配置
export function preConfigAdd(data) {
  return request({
    url: '/api/user/preconfig/add',
    method: 'post',
    data
  })
}

// 更新用户配置
export function preConfigUpdate(data) {
  return request({
    url: '/api/user/preconfig/update',
    method: 'post',
    data
  })
}

// 删除用户配置
export function preConfigDelete(data) {
  return request({
    url: '/api/user/preconfig/delete',
    method: 'post',
    data
  })
}

// 根据用户名获取配置
export function preConfigGetByUsername(params) {
  return request({
    url: '/api/user/preconfig/getByUsername',
    method: 'get',
    params
  })
}

// 同步已有用户到预配置
export function preConfigSyncUsers() {
  return request({
    url: '/api/user/preconfig/syncUsers',
    method: 'post'
  })
}

