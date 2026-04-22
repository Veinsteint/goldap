import request from '@/utils/request'

// 获取当前登录用户信息 
export function getInfo() {
  return request({
    url: '/api/user/info',
    method: 'get'
  })
}

// 获取用户列表 
export function getUsers(params) {
  return request({
    url: '/api/user/list',
    method: 'get',
    params
  })
}

// 更新用户登录密码
export function changePwd(data) {
  return request({
    url: '/api/user/changePwd',
    method: 'post',
    data
  })
}

// 创建用户
export function createUser(data) {
  return request({
    url: '/api/user/add',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUserById(data) {
  return request({
    url: '/api/user/update',
    method: 'post',
    data
  })
}
// 批量删除记录*
export function batchDeleteUserByIds(data) {
  return request({
    url: '/api/user/delete',
    method: 'post',
    data
  })
}

// 更改用户状态
export function changeUserStatus(data) {
  return request({
    url: '/api/user/changeUserStatus',
    method: 'post',
    data
  })
}

// 同步openldap用户信息
export function syncOpenLdapUsersApi(data) {
  return request({
    url: '/api/user/syncOpenLdapUsers',
    method: 'post',
    data
  })
}

// 同步Sql中的用户到ldap
export function syncSqlUsers(data) {
  return request({
    url: '/api/user/syncSqlUsers',
    method: 'post',
    data
  })
}
