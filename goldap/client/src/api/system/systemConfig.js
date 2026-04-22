import request from '@/utils/request'

// 获取系统配置
export function getSystemConfig() {
  return request({
    url: '/api/system/config/get',
    method: 'get'
  })
}

// 更新系统配置
export function updateSystemConfig(data) {
  return request({
    url: '/api/system/config/update',
    method: 'post',
    data
  })
}

// 获取系统配置（公开接口，无需登录）
export function getSystemConfigPublic() {
  return request({
    url: '/api/base/systemConfig',
    method: 'get'
  })
}

