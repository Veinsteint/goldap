import request from '@/utils/request'

// 获取操作日志列表
export function getOperationLogs(params) {
  return request({
    url: '/api/log/operation/list',
    method: 'get',
    params
  })
}

// 批量删除操作日志
export function batchDeleteOperationLogByIds(data) {
  return request({
    url: '/api/log/operation/delete',
    method: 'post',
    data
  })
}

export function CleanOperationLog() {
  return request({
    url: '/api/log/operation/clean',
    method: 'delete'
  })
}
