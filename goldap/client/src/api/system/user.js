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
// 输入邮箱获取验证码
export function sendCode(data) {
  return request({
    url: '/api/base/sendcode',
    method: 'post',
    data
  })
}
// 邮箱更新用户密码
export function emailPass(data) {
  return request({
    url: '/api/base/changePwd',
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

// 重置用户密码
export function resetPassword(data) {
  return request({
    url: '/api/user/resetPassword',
    method: 'post',
    data
  })
}

// 用户注册
export function register(data) {
  return request({
    url: '/api/base/register',
    method: 'post',
    data
  })
}

// 获取SSH公钥
export function getSSHKeys() {
  return request({
    url: '/api/user/ssh-keys',
    method: 'get'
  })
}

// 添加SSH公钥
export function addSSHKey(data) {
  return request({
    url: '/api/user/ssh-keys',
    method: 'post',
    data
  })
}

// 删除SSH公钥
export function deleteSSHKey(id) {
  return request({
    url: `/api/user/ssh-keys/${id}`,
    method: 'delete'
  })
}

// 获取待审核用户列表
export function getPendingUsers(params) {
  return request({
    url: '/api/user/pending/list',
    method: 'get',
    params
  })
}

// 审核待审核用户
export function reviewPendingUser(data) {
  return request({
    url: '/api/user/pending/review',
    method: 'post',
    data
  })
}

// 删除待审核用户
export function deletePendingUsers(data) {
  return request({
    url: '/api/user/pending/delete',
    method: 'post',
    data
  })
}

// 获取注册模式
export function getRegistrationMode() {
  return request({
    url: '/api/base/registrationMode',
    method: 'get'
  })
}

// 获取有效用户名列表（用于注册）
export function getValidUsernames() {
  return request({
    url: '/api/base/validUsernames',
    method: 'get'
  })
}
