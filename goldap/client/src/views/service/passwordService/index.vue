<template>
  <div class="password-service-container" :style="containerStyle">
    <div class="service-card">
      <div class="header">
        <h2 class="title">
          <svg-icon icon-class="lock" style="margin-right: 8px;" />
          自助密码服务
        </h2>
        <p class="subtitle">安全便捷的账户管理服务</p>
      </div>

      <!-- 系统维护提示 -->
      <el-alert
        v-if="systemConfig.systemMaintenanceMode"
        :title="systemConfig.maintenanceMessage"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px;"
      />

      <el-tabs v-model="activeTab" class="service-tabs" :disabled="!systemConfig.passwordSelfServiceEnabled">
        <!-- 账号注册 -->
        <el-tab-pane label="账号注册" name="register" :disabled="systemConfig.registrationDisabled">
          <div class="tab-content">
            <!-- 禁止注册提示 -->
            <el-alert
              v-if="systemConfig.registrationDisabled"
              title="注册已关闭"
              type="error"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              <template slot="default">
                <p>系统已禁止新用户注册，请联系管理员。</p>
              </template>
            </el-alert>
            <!-- 注册模式提示 -->
            <el-alert
              v-else-if="registrationMode === 'preconfig'"
              title="注册限制提示"
              type="warning"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            >
              <template slot="default">
                <p>当前注册模式为<strong>预配置模式</strong>，只允许预配置列表中的用户名进行注册。</p>
                <el-button type="text" icon="el-icon-view" @click="showValidUsernamesDialog">点击查看有效用户名列表</el-button>
              </template>
            </el-alert>

            <el-form ref="registerForm" :model="registerForm" :rules="registerRules" size="medium" class="service-form" autocomplete="off" :disabled="systemConfig.registrationDisabled">
              <el-form-item label="用户名" prop="username">
                <el-input
                  v-model="registerForm.username"
                  placeholder="请输入用户名（3-50个字符）"
                  prefix-icon="el-icon-user"
                  maxlength="50"
                  autocomplete="off"
                  name="register-username"
                >
                  <template v-if="registrationMode === 'preconfig'" slot="append">
                    <el-button icon="el-icon-view" @click="showValidUsernamesDialog">查看有效用户名</el-button>
                  </template>
                </el-input>
                <div class="form-tip">
                  <span v-if="registrationMode === 'preconfig'">
                    <el-tag type="warning" size="mini">注意</el-tag> 只有预配置列表中的用户名才能注册
                  </span>
                  <span v-else>用户名将用于登录，请妥善保管</span>
                </div>
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input
                  v-model="registerForm.email"
                  placeholder="请输入邮箱地址"
                  prefix-icon="el-icon-message"
                  autocomplete="off"
                  name="register-email"
                />
                <div class="form-tip">用于接收验证码和重要通知</div>
              </el-form-item>
              <el-form-item label="真实姓名" prop="realName">
                <el-input
                  v-model="registerForm.realName"
                  placeholder="请输入真实姓名（仅支持英文字母）"
                  prefix-icon="el-icon-user-solid"
                  maxlength="50"
                  autocomplete="off"
                  name="register-realname"
                />
                <div class="form-tip">只能包含英文字母（a-z, A-Z）</div>
              </el-form-item>
              <el-form-item label="密码" prop="password">
                <el-input
                  v-model="registerForm.password"
                  type="password"
                  placeholder="请输入密码（至少6位）"
                  prefix-icon="el-icon-lock"
                  show-password
                  autocomplete="new-password"
                  name="register-password"
                  @input="checkRegisterPasswordStrength"
                />
                <div v-if="registerForm.password" class="password-strength">
                  <div class="strength-label">密码强度：</div>
                  <div class="strength-bar">
                    <div
                      :class="['strength-item', registerPasswordStrength.level]"
                      :style="{ width: registerPasswordStrength.percent + '%' }"
                    />
                  </div>
                  <span :class="['strength-text', registerPasswordStrength.level]">{{ registerPasswordStrength.text }}</span>
                </div>
              </el-form-item>
              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input
                  v-model="registerForm.confirmPassword"
                  type="password"
                  placeholder="请再次输入密码"
                  prefix-icon="el-icon-lock"
                  show-password
                  autocomplete="new-password"
                  name="register-confirm-password"
                />
              </el-form-item>
              <el-form-item label="备注" prop="remark">
                <el-input
                  v-model="registerForm.remark"
                  type="textarea"
                  :rows="3"
                  placeholder="选填：注册原因或备注信息"
                  maxlength="200"
                  show-word-limit
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="registerLoading" :disabled="systemConfig.registrationDisabled" @click="handleRegister">注册</el-button>
                <el-button @click="resetRegisterForm">重置</el-button>
                <el-button type="text" @click="goToLogin">已有账号？去登录</el-button>
              </el-form-item>
            </el-form>
            <el-alert
              v-if="registerSuccess"
              title="注册成功"
              type="success"
              :closable="false"
              show-icon
              style="margin-top: 20px;"
            >
              <template slot="default">
                <p>您的注册申请已提交，等待管理员审核后即可使用。</p>
                <p style="margin-top: 10px;">审核结果将通过邮箱通知您。</p>
              </template>
            </el-alert>
          </div>
        </el-tab-pane>

        <!-- 忘记密码 -->
        <el-tab-pane label="忘记密码" name="forgot" :disabled="!systemConfig.passwordSelfServiceEnabled">
          <div class="tab-content">
            <el-steps :active="forgotStep" finish-status="success" align-center>
              <el-step title="验证身份" description="通过邮箱验证"></el-step>
              <el-step title="确认重置" description="系统生成新密码"></el-step>
              <el-step title="完成" description="重置成功"></el-step>
            </el-steps>

            <!-- 步骤1: 验证身份 -->
            <div v-if="forgotStep === 0" class="step-content">
              <el-form ref="forgotForm" :model="forgotForm" :rules="forgotRules" size="medium" class="service-form">
                <el-form-item label="邮箱地址" prop="mail">
                  <el-input
                    v-model="forgotForm.mail"
                    placeholder="请输入注册时使用的邮箱地址"
                    prefix-icon="el-icon-message"
                  >
                    <template slot="append">
                      <el-button
                        type="primary"
                        :loading="codeLoading"
                        :disabled="codeDisabled"
                        @click="sendEmailCode"
                      >
                        {{ codeButtonText }}
                      </el-button>
                    </template>
                  </el-input>
                </el-form-item>
                <el-form-item label="验证码" prop="code">
                  <el-input
                    v-model="forgotForm.code"
                    placeholder="请输入邮箱收到的验证码"
                    prefix-icon="el-icon-key"
                    maxlength="6"
                  />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="verifyLoading" @click="verifyCode">下一步</el-button>
                  <el-button @click="resetForgotForm">重置</el-button>
                </el-form-item>
              </el-form>
            </div>

            <!-- 步骤2: 确认重置 -->
            <div v-if="forgotStep === 1" class="step-content">
              <el-alert
                title="确认重置密码"
                type="info"
                :closable="false"
                show-icon
                style="margin-bottom: 20px;"
              >
                <template slot="default">
                  <p>系统将自动生成新密码并发送到您的邮箱：<strong>{{ forgotForm.mail }}</strong></p>
                  <p style="margin-top: 10px; color: #e6a23c;">请确保您能访问该邮箱，新密码将在邮件中提供。</p>
                </template>
              </el-alert>
              <el-form ref="resetForm" :model="resetForm" :rules="resetRules" size="medium" class="service-form">
                <el-form-item>
                  <el-button type="primary" :loading="resetLoading" @click="resetPassword">确认重置密码</el-button>
                  <el-button @click="forgotStep = 0">上一步</el-button>
                </el-form-item>
              </el-form>
            </div>

            <!-- 步骤3: 完成 -->
            <div v-if="forgotStep === 2" class="step-content success-content">
              <el-result icon="success" title="密码重置成功！" sub-title="您的密码已成功重置，请使用新密码登录">
                <template slot="extra">
                  <el-button type="primary" @click="goToLogin">前往登录</el-button>
                </template>
              </el-result>
            </div>
          </div>
        </el-tab-pane>

        <!-- 修改密码（已登录用户） -->
        <el-tab-pane label="修改密码" name="change" :disabled="!isLoggedIn || !systemConfig.passwordSelfServiceEnabled">
          <div class="tab-content">
            <el-alert
              v-if="!isLoggedIn"
              title="提示"
              type="warning"
              :closable="false"
              show-icon
            >
              <template slot="default">
                修改密码功能需要先登录，请先 <router-link to="/login">登录</router-link> 后再使用此功能。
                或前往 <router-link to="/profile">个人中心</router-link> 修改密码。
              </template>
            </el-alert>

            <el-form
              v-else
              ref="changeForm"
              :model="changeForm"
              :rules="changeRules"
              size="medium"
              class="service-form"
            >
              <el-form-item label="原密码" prop="oldPassword">
                <el-input
                  v-model="changeForm.oldPassword"
                  type="password"
                  placeholder="请输入当前密码"
                  prefix-icon="el-icon-lock"
                  show-password
                />
              </el-form-item>
              <el-form-item label="新密码" prop="newPassword">
                <el-input
                  v-model="changeForm.newPassword"
                  type="password"
                  placeholder="请输入新密码（至少6位）"
                  prefix-icon="el-icon-lock"
                  show-password
                  @input="checkPasswordStrength"
                />
                <div v-if="changeForm.newPassword" class="password-strength">
                  <div class="strength-label">密码强度：</div>
                  <div class="strength-bar">
                    <div
                      :class="['strength-item', passwordStrength.level]"
                      :style="{ width: passwordStrength.percent + '%' }"
                    />
                  </div>
                  <span :class="['strength-text', passwordStrength.level]">{{ passwordStrength.text }}</span>
                </div>
              </el-form-item>
              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input
                  v-model="changeForm.confirmPassword"
                  type="password"
                  placeholder="请再次输入新密码"
                  prefix-icon="el-icon-lock"
                  show-password
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="changeLoading" @click="changePassword">确认修改</el-button>
                <el-button @click="resetChangeForm">重置</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- SSH公钥管理 -->
        <el-tab-pane label="SSH公钥" name="ssh" :disabled="!isLoggedIn">
          <div class="tab-content">
            <el-alert
              v-if="!isLoggedIn"
              title="提示"
              type="warning"
              :closable="false"
              show-icon
            >
              <template slot="default">
                SSH公钥管理功能需要先登录，请先 <router-link to="/login">登录</router-link> 后再使用此功能。
              </template>
            </el-alert>

            <div v-else>
              <!-- 添加SSH公钥 -->
              <el-card shadow="never" style="margin-bottom: 20px;">
                <div slot="header">
                  <span>添加SSH公钥</span>
                </div>
                <el-form ref="sshForm" :model="sshForm" :rules="sshRules" size="medium">
                  <el-form-item label="标题" prop="title">
                    <el-input
                      v-model="sshForm.title"
                      placeholder="为这个密钥起个名字，例如：cmplab-wwtian-ssh"
                      maxlength="100"
                    />
                  </el-form-item>
                  <el-form-item label="公钥" prop="key">
                    <el-input
                      v-model="sshForm.key"
                      type="textarea"
                      :rows="6"
                      placeholder="粘贴您的SSH公钥内容，例如：ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
                      maxlength="5000"
                    />
                    <div class="form-tip">
                      <p>如何生成SSH密钥？</p>
                      <p>在终端运行：<code>ssh-keygen -t rsa -b 4096 -C "your_email@cimrbj.ac.cn"</code></p>
                      <p>然后复制 <code>~/.ssh/id_rsa.pub</code> 文件的内容</p>
                    </div>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" :loading="sshAddLoading" @click="addSSHKey">添加公钥</el-button>
                    <el-button @click="resetSSHForm">重置</el-button>
                  </el-form-item>
                </el-form>
              </el-card>

              <!-- SSH公钥列表 -->
              <el-card shadow="never">
                <div slot="header">
                  <span>已添加的SSH公钥</span>
                  <el-button style="float: right; padding: 3px 0" type="text" @click="loadSSHKeys">刷新</el-button>
                </div>
                <el-table
                  v-loading="sshListLoading"
                  :data="sshKeys"
                  style="width: 100%"
                  empty-text="暂无SSH公钥"
                >
                  <el-table-column prop="title" label="标题" width="200" />
                  <el-table-column prop="key" label="公钥" show-overflow-tooltip>
                    <template slot-scope="scope">
                      <code style="font-size: 12px;">{{ scope.row.key.substring(0, 50) }}...</code>
                    </template>
                  </el-table-column>
                  <el-table-column prop="createdAt" label="添加时间" width="180" />
                  <el-table-column label="操作" width="100" align="center">
                    <template slot-scope="scope">
                      <el-button
                        type="danger"
                        size="mini"
                        :loading="scope.row.deleting"
                        @click="deleteSSHKey(scope.row.id)"
                      >
                        删除
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </el-card>
            </div>
          </div>
        </el-tab-pane>

        <!-- 密码安全提示 -->
        <el-tab-pane label="安全提示" name="tips">
          <div class="tab-content tips-content">
            <el-card shadow="never">
              <div slot="header">
                <span>密码安全建议</span>
              </div>
              <div class="tips-list">
                <el-alert
                  title="创建强密码"
                  type="info"
                  :closable="false"
                  show-icon
                >
                  <ul>
                    <li>使用至少8-12个字符的组合</li>
                    <li>包含大写字母、小写字母、数字和特殊字符</li>
                    <li>避免使用个人信息（如姓名、生日）</li>
                    <li>不要使用常见的密码模式（如123456、password等）</li>
                  </ul>
                </el-alert>

                <el-alert
                  title="密码管理"
                  type="warning"
                  :closable="false"
                  show-icon
                  style="margin-top: 20px;"
                >
                  <ul>
                    <li>定期更换密码（建议每3-6个月）</li>
                    <li>不要在多个账户使用相同密码</li>
                    <li>不要将密码告诉他人或写在明显的地方</li>
                    <li>使用密码管理器来安全存储密码</li>
                  </ul>
                </el-alert>

                <el-alert
                  title="SSH密钥安全"
                  type="success"
                  :closable="false"
                  show-icon
                  style="margin-top: 20px;"
                >
                  <ul>
                    <li>为每个设备生成独立的SSH密钥</li>
                    <li>使用强密码保护您的私钥</li>
                    <li>定期轮换SSH密钥</li>
                    <li>不要将私钥分享给他人或上传到公共仓库</li>
                    <li>删除不再使用的SSH密钥</li>
                  </ul>
                </el-alert>

                <el-alert
                  title="账户安全"
                  type="success"
                  :closable="false"
                  show-icon
                  style="margin-top: 20px;"
                >
                  <ul>
                    <li>启用双因素认证（如果可用）</li>
                    <li>定期检查账户登录记录</li>
                    <li>发现异常活动立即修改密码</li>
                    <li>保持邮箱和手机号信息更新</li>
                  </ul>
                </el-alert>
              </div>
            </el-card>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 有效用户名列表弹窗 -->
    <el-dialog
      title="有效用户名列表"
      :visible.sync="validUsernamesDialogVisible"
      width="600px"
      append-to-body
    >
      <el-alert
        title="注册说明"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 20px;"
      >
        <template slot="default">
          <p>只有以下列表中状态为"未注册"的用户名才能进行注册。</p>
        </template>
      </el-alert>

      <el-table
        v-loading="validUsernamesLoading"
        :data="validUsernames"
        style="width: 100%"
        max-height="400"
        empty-text="暂无预配置用户"
      >
        <el-table-column prop="username" label="用户名" width="180" />
        <el-table-column prop="nickname" label="姓名" width="150" />
        <el-table-column prop="status" label="状态" width="120">
          <template slot-scope="scope">
            <el-tag
              :type="getStatusTagType(scope.row.status)"
              size="small"
            >
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
      </el-table>

      <span slot="footer" class="dialog-footer">
        <el-button @click="validUsernamesDialogVisible = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { emailPass, sendCode, register, getSSHKeys, addSSHKey, deleteSSHKey, getRegistrationMode, getValidUsernames } from '@/api/system/user'
import { getContentHeight, HEADER_HEIGHT_WITH_TAGS } from '@/utils/layout'
import { changePwd } from '@/api/system/user'
import { validEmail } from '@/utils/validate'
import { Message } from 'element-ui'
import JSEncrypt from 'jsencrypt'
import store from '@/store'
import { getSystemConfigPublic } from '@/api/system/systemConfig'

export default {
  name: 'PasswordService',
  data() {
    const validateEmail = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入邮箱地址'))
      } else if (!validEmail(value)) {
        callback(new Error('请输入正确的邮箱地址'))
      } else {
        callback()
      }
    }
    const validateConfirmPassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.changeForm.newPassword) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    const validateRegisterConfirmPassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.registerForm.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入密码'))
      } else if (value.length < 6) {
        callback(new Error('密码长度不能少于6位'))
      } else {
        callback()
      }
    }
    const validateRealName = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入真实姓名'))
      } else if (value.length > 50) {
        callback(new Error('姓名长度不能超过50个字符'))
      } else {
        // 只允许英文字母
        const englishOnlyPattern = /^[a-zA-Z]+$/
        if (!englishOnlyPattern.test(value)) {
          callback(new Error('真实姓名只能包含英文字母'))
        } else {
          callback()
        }
      }
    }
    const validateSSHKey = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入SSH公钥'))
      } else if (!value.trim().startsWith('ssh-')) {
        callback(new Error('SSH公钥格式不正确，应以ssh-开头'))
      } else if (value.trim().length < 50) {
        callback(new Error('SSH公钥长度不足'))
      } else {
        callback()
      }
    }

    return {
      activeTab: 'register',
      // 系统配置相关
      systemConfig: {
        passwordSelfServiceEnabled: true,
        registrationDisabled: false,
        systemMaintenanceMode: false,
        maintenanceMessage: '系统正在升级维护中，请稍后再试'
      },
      // 注册模式相关
      registrationMode: 'open',
      validUsernamesDialogVisible: false,
      validUsernamesLoading: false,
      validUsernames: [],
      // 注册相关
      registerForm: {
        username: '',
        email: '',
        realName: '',
        password: '',
        confirmPassword: '',
        remark: ''
      },
      registerRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 50, message: '用户名长度在3到50个字符', trigger: 'blur' }
        ],
        email: [{ required: true, trigger: 'blur', validator: validateEmail }],
        realName: [
          { required: true, trigger: 'blur', validator: validateRealName }
        ],
        password: [
          { required: true, trigger: 'blur', validator: validatePassword }
        ],
        confirmPassword: [
          { required: true, trigger: 'blur', validator: validateRegisterConfirmPassword }
        ],
        remark: [
          { max: 200, message: '备注长度不能超过200个字符', trigger: 'blur' }
        ]
      },
      registerLoading: false,
      registerSuccess: false,
      registerPasswordStrength: { level: 'weak', percent: 0, text: '弱' },
      // 忘记密码相关
      forgotStep: 0,
      forgotForm: {
        mail: '',
        code: ''
      },
      forgotRules: {
        mail: [{ required: true, trigger: 'blur', validator: validateEmail }],
        code: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
      },
      codeLoading: false,
      codeDisabled: false,
      codeButtonText: '发送验证码',
      countdown: 0,
      verifyLoading: false,
      // 重置密码相关
      resetForm: {},
      resetRules: {},
      resetLoading: false,
      // 修改密码相关
      changeForm: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      changeRules: {
        oldPassword: [
          { required: true, message: '请输入原密码', trigger: 'blur' },
          { min: 6, max: 30, message: '长度在 6 到 30 个字符', trigger: 'blur' }
        ],
        newPassword: [
          { required: true, trigger: 'blur', validator: validatePassword }
        ],
        confirmPassword: [
          { required: true, trigger: 'blur', validator: validateConfirmPassword }
        ]
      },
      changeLoading: false,
      // 密码强度
      passwordStrength: {
        level: 'weak',
        percent: 0,
        text: '弱'
      },
      passwordChecks: {
        length: false,
        uppercase: false,
        lowercase: false,
        number: false,
        special: false
      },
      // SSH公钥相关
      sshForm: {
        title: '',
        key: ''
      },
      sshRules: {
        title: [
          { required: true, message: '请输入标题', trigger: 'blur' },
          { max: 100, message: '标题长度不能超过100个字符', trigger: 'blur' }
        ],
        key: [
          { required: true, trigger: 'blur', validator: validateSSHKey }
        ]
      },
      sshAddLoading: false,
      sshKeys: [],
      sshListLoading: false,
      publicKey: process.env.VUE_APP_PUBLIC_KEY
    }
  },
  computed: {
    isLoggedIn() {
      return !!store.getters.token
    },
    containerMinHeight() {
      const hasTagsView = this.$store.state.settings.tagsView
      return getContentHeight(hasTagsView, 0)
    },
    containerStyle() {
      return {
        '--header-height-with-tags': `${HEADER_HEIGHT_WITH_TAGS}px`,
        minHeight: this.containerMinHeight
      }
    }
  },
  watch: {
    activeTab() {
      this.$nextTick(() => {
        this.resetAllForms()
        if (this.activeTab === 'ssh' && this.isLoggedIn) {
          this.loadSSHKeys()
        }
      })
    }
  },
  mounted() {
    this.loadSystemConfig()
    this.loadRegistrationMode()
    if (this.isLoggedIn && this.activeTab === 'ssh') {
      this.loadSSHKeys()
    }
  },
  methods: {
    // 加载系统配置
    async loadSystemConfig() {
      try {
        const res = await getSystemConfigPublic()
        if (res && (res.code === 0 || res.code === 200)) {
          this.systemConfig = {
            passwordSelfServiceEnabled: res.data?.passwordSelfServiceEnabled !== false,
            registrationDisabled: res.data?.registrationDisabled === true,
            systemMaintenanceMode: res.data?.systemMaintenanceMode === true,
            maintenanceMessage: res.data?.maintenanceMessage || '系统正在升级维护中，请稍后再试'
          }
        }
      } catch (error) {
        console.error('Failed to load system config:', error)
      }
    },
    // 加载注册模式
    async loadRegistrationMode() {
      try {
        const res = await getRegistrationMode()
        if (res && (res.code === 0 || res.code === 200)) {
          this.registrationMode = res.data?.mode || 'open'
        }
      } catch (error) {
        // 默认为开放模式
        this.registrationMode = 'open'
      }
    },
    // 显示有效用户名列表弹窗
    async showValidUsernamesDialog() {
      this.validUsernamesDialogVisible = true
      this.validUsernamesLoading = true
      try {
        const res = await getValidUsernames()
        if (res && (res.code === 0 || res.code === 200)) {
          this.validUsernames = res.data || []
        }
      } catch (error) {
        Message({
          showClose: true,
          message: error.message || '获取有效用户名列表失败',
          type: 'error'
        })
      } finally {
        this.validUsernamesLoading = false
      }
    },
    // 获取状态标签类型
    getStatusTagType(status) {
      switch (status) {
        case '已注册-有效':
          return 'success'
        case '已注册-失效':
          return 'danger'
        case '待审核':
          return 'warning'
        case '未注册':
          return 'info'
        default:
          return 'info'
      }
    },
    // 选择用户名并填入
    fillUsername(username) {
      this.registerForm.username = username
      this.validUsernamesDialogVisible = false
      Message({
        showClose: true,
        message: `已选择用户名: ${username}`,
        type: 'success'
      })
    },
    // 注册
    async handleRegister() {
      this.$refs.registerForm.validate(async (valid) => {
        if (valid) {
          this.registerLoading = true
          try {
            const formData = {
              username: this.registerForm.username.trim(),
              email: this.registerForm.email.trim(),
              realName: this.registerForm.realName.trim(),
              remark: this.registerForm.remark.trim()
            }
            // 密码RSA加密处理
            if (this.publicKey) {
              const encryptor = new JSEncrypt()
              encryptor.setPublicKey(this.publicKey)
              formData.password = encryptor.encrypt(this.registerForm.password)
            } else {
              // 如果没有公钥，直接使用明文（不推荐，但为了兼容性）
              formData.password = this.registerForm.password
            }
            const res = await register(formData)
            if (this.judgeResult(res)) {
              this.registerSuccess = true
              this.resetRegisterForm()
            }
          } catch (error) {
            Message({
              showClose: true,
              message: error.message || '注册失败',
              type: 'error'
            })
          } finally {
            this.registerLoading = false
          }
        }
      })
    },
    // 检查注册密码强度
    checkRegisterPasswordStrength() {
      const password = this.registerForm.password
      if (!password) {
        this.registerPasswordStrength = { level: 'weak', percent: 0, text: '弱' }
        return
      }
      const checks = {
        length: password.length >= 6,
        uppercase: /[A-Z]/.test(password),
        lowercase: /[a-z]/.test(password),
        number: /[0-9]/.test(password),
        special: /[!@#$%^&*(),.?":{}|<>]/.test(password)
      }
      let score = 0
      if (checks.length) score++
      if (checks.uppercase) score++
      if (checks.lowercase) score++
      if (checks.number) score++
      if (checks.special) score++
      if (score <= 2) {
        this.registerPasswordStrength = { level: 'weak', percent: 33, text: '弱' }
      } else if (score <= 3) {
        this.registerPasswordStrength = { level: 'medium', percent: 66, text: '中' }
      } else {
        this.registerPasswordStrength = { level: 'strong', percent: 100, text: '强' }
      }
    },
    // 发送邮箱验证码
    async sendEmailCode() {
      this.$refs.forgotForm.validateField('mail', async (valid) => {
        if (!valid) {
          this.codeLoading = true
          try {
            const res = await sendCode({ mail: this.forgotForm.mail })
            if (this.judgeResult(res)) {
              this.startCountdown()
            }
          } catch (error) {
            Message({
              showClose: true,
              message: error.message || '发送验证码失败',
              type: 'error'
            })
          } finally {
            this.codeLoading = false
          }
        }
      })
    },
    // 倒计时
    startCountdown() {
      this.countdown = 60
      this.codeDisabled = true
      const timer = setInterval(() => {
        this.countdown--
        this.codeButtonText = `${this.countdown}秒后重新发送`
        if (this.countdown <= 0) {
          clearInterval(timer)
          this.codeDisabled = false
          this.codeButtonText = '发送验证码'
        }
      }, 1000)
    },
    // 验证验证码
    async verifyCode() {
      this.$refs.forgotForm.validate(async (valid) => {
        if (valid) {
          this.verifyLoading = true
          this.forgotStep = 1
          this.verifyLoading = false
        }
      })
    },
    // 重置密码
    async resetPassword() {
      this.resetLoading = true
      try {
        const formData = {
          mail: this.forgotForm.mail,
          code: this.forgotForm.code
        }
        const res = await emailPass(formData)
        if (this.judgeResult(res)) {
          this.forgotStep = 2
        }
      } catch (error) {
        Message({
          showClose: true,
          message: error.message || '重置密码失败',
          type: 'error'
        })
      } finally {
        this.resetLoading = false
      }
    },
    // 修改密码
    async changePassword() {
      this.$refs.changeForm.validate(async (valid) => {
        if (valid) {
          this.changeLoading = true
          try {
            const formData = { ...this.changeForm }
            // 密码RSA加密处理
            const encryptor = new JSEncrypt()
            encryptor.setPublicKey(this.publicKey)
            formData.oldPassword = encryptor.encrypt(this.changeForm.oldPassword)
            formData.newPassword = encryptor.encrypt(this.changeForm.newPassword)
            formData.confirmPassword = encryptor.encrypt(this.changeForm.confirmPassword)

            const res = await changePwd(formData)
            if (this.judgeResult(res)) {
              Message({
                showClose: true,
                message: '密码修改成功，请重新登录',
                type: 'success'
              })
              setTimeout(() => {
                store.dispatch('user/logout').then(() => {
                  this.$router.push('/login')
                })
              }, 1500)
            }
          } catch (error) {
            Message({
              showClose: true,
              message: error.message || '修改密码失败',
              type: 'error'
            })
          } finally {
            this.changeLoading = false
          }
        }
      })
    },
    // 检查密码强度
    checkPasswordStrength() {
      const password = this.changeForm.newPassword
      if (!password) {
        this.passwordStrength = { level: 'weak', percent: 0, text: '弱' }
        return
      }
      this.passwordChecks = {
        length: password.length >= 6,
        uppercase: /[A-Z]/.test(password),
        lowercase: /[a-z]/.test(password),
        number: /[0-9]/.test(password),
        special: /[!@#$%^&*(),.?":{}|<>]/.test(password)
      }
      let score = 0
      if (this.passwordChecks.length) score++
      if (this.passwordChecks.uppercase) score++
      if (this.passwordChecks.lowercase) score++
      if (this.passwordChecks.number) score++
      if (this.passwordChecks.special) score++
      if (score <= 2) {
        this.passwordStrength = { level: 'weak', percent: 33, text: '弱' }
      } else if (score <= 3) {
        this.passwordStrength = { level: 'medium', percent: 66, text: '中' }
      } else {
        this.passwordStrength = { level: 'strong', percent: 100, text: '强' }
      }
    },
    // SSH公钥管理
    async loadSSHKeys() {
      this.sshListLoading = true
      try {
        const res = await getSSHKeys()
        if (res === false) {
          this.sshKeys = []
          return
        }
        if (res && (res.code === 0 || res.code === 200)) {
          this.sshKeys = res.data || []
        } else {
          // 如果返回了其他code，设置空数组
          this.sshKeys = []
        }
      } catch (error) {
        // 只有在网络错误或其他异常时才显示错误消息
        // 业务错误（如权限问题）已经在响应拦截器中处理
        if (error.response && error.response.status >= 500) {
          Message({
            showClose: true,
            message: error.message || '加载SSH公钥列表失败',
            type: 'error'
          })
        }
        this.sshKeys = []
      } finally {
        this.sshListLoading = false
      }
    },
    async addSSHKey() {
      this.$refs.sshForm.validate(async (valid) => {
        if (valid) {
          this.sshAddLoading = true
          try {
            const formData = {
              title: this.sshForm.title.trim(),
              key: this.sshForm.key.trim()
            }
            const res = await addSSHKey(formData)
            if (this.judgeResult(res)) {
              this.resetSSHForm()
              this.loadSSHKeys()
            }
          } catch (error) {
            Message({
              showClose: true,
              message: error.message || '添加SSH公钥失败',
              type: 'error'
            })
          } finally {
            this.sshAddLoading = false
          }
        }
      })
    },
    async deleteSSHKey(id) {
      this.$confirm('确定要删除这个SSH公钥吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const res = await deleteSSHKey(id)
          if (this.judgeResult(res)) {
            this.loadSSHKeys()
          }
        } catch (error) {
          Message({
            showClose: true,
            message: error.message || '删除SSH公钥失败',
            type: 'error'
          })
        }
      }).catch(() => {})
    },
    // 判断结果
    judgeResult(res) {
      if (res.code === 0 || res.code === 200) {
        Message({
          showClose: true,
          message: res.msg || '操作成功',
          type: 'success'
        })
        return true
      } else {
        Message({
          showClose: true,
          message: res.msg || '操作失败',
          type: 'error'
        })
        return false
      }
    },
    // 重置表单
    resetRegisterForm() {
      if (this.$refs.registerForm) {
        this.$refs.registerForm.clearValidate()
        this.$nextTick(() => {
          this.$refs.registerForm?.resetFields()
        })
      }
      this.registerForm = {
        username: '',
        email: '',
        realName: '',
        password: '',
        confirmPassword: '',
        remark: ''
      }
      this.registerSuccess = false
      this.registerPasswordStrength = { level: 'weak', percent: 0, text: '弱' }
    },
    resetForgotForm() {
      if (this.$refs.forgotForm) {
        this.$refs.forgotForm.clearValidate()
        this.$nextTick(() => {
          this.$refs.forgotForm?.resetFields()
        })
      }
      this.forgotForm = { mail: '', code: '' }
      this.forgotStep = 0
    },
    resetChangeForm() {
      if (this.$refs.changeForm) {
        this.$refs.changeForm.clearValidate()
        this.$nextTick(() => {
          this.$refs.changeForm?.resetFields()
        })
      }
      this.changeForm = { oldPassword: '', newPassword: '', confirmPassword: '' }
      this.passwordStrength = { level: 'weak', percent: 0, text: '弱' }
      this.passwordChecks = {
        length: false,
        uppercase: false,
        lowercase: false,
        number: false,
        special: false
      }
    },
    resetSSHForm() {
      if (this.$refs.sshForm) {
        this.$refs.sshForm.clearValidate()
        this.$nextTick(() => {
          this.$refs.sshForm?.resetFields()
        })
      }
      this.sshForm = { title: '', key: '' }
    },
    resetAllForms() {
      // 只重置当前可见标签页的表单，避免触发不可见表单的验证
      if (this.activeTab === 'register') {
        this.resetRegisterForm()
      } else if (this.activeTab === 'forgot') {
        this.resetForgotForm()
      } else if (this.activeTab === 'change') {
        this.resetChangeForm()
      } else if (this.activeTab === 'ssh') {
        this.resetSSHForm()
      }
      // 重置忘记密码步骤
      if (this.activeTab !== 'forgot') {
        this.forgotStep = 0
      }
    },
    // 前往登录
    goToLogin() {
      this.$router.push('/login')
    }
  }
}
</script>

<style lang="scss" scoped>
.password-service-container {
  min-height: calc(100vh - var(--header-height-with-tags, 84px));
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;
}

.service-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 800px;
  max-height: 1100px;
  padding: 30px;

  .header {
    text-align: center;
    margin-bottom: 30px;

    .title {
      font-size: 28px;
      color: #333;
      margin: 0 0 10px 0;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .subtitle {
      color: #666;
      font-size: 14px;
      margin: 0;
    }
  }
}

.service-tabs {
  ::v-deep .el-tabs__header {
    margin-bottom: 30px;
  }

  ::v-deep .el-tabs__item {
    font-size: 16px;
    padding: 0 30px;
  }
}

.tab-content {
  min-height: 400px;
}

.step-content {
  margin-top: 30px;
}

.service-form {
  max-width: 500px;
  margin: 30px auto 0;

  ::v-deep .el-form-item__label {
    font-weight: 500;
  }
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  line-height: 1.5;

  code {
    background: #f5f7fa;
    padding: 2px 6px;
    border-radius: 3px;
    font-family: 'Courier New', monospace;
  }
}

.password-strength {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;

  .strength-label {
    font-size: 12px;
    color: #666;
    white-space: nowrap;
  }

  .strength-bar {
    flex: 1;
    height: 6px;
    background: #f0f0f0;
    border-radius: 3px;
    overflow: hidden;

    .strength-item {
      height: 100%;
      transition: all 0.3s;

      &.weak {
        background: #f56c6c;
      }

      &.medium {
        background: #e6a23c;
      }

      &.strong {
        background: #67c23a;
      }
    }
  }

  .strength-text {
    font-size: 12px;
    font-weight: 500;
    white-space: nowrap;

    &.weak {
      color: #f56c6c;
    }

    &.medium {
      color: #e6a23c;
    }

    &.strong {
      color: #67c23a;
    }
  }
}

.success-content {
  padding: 40px 0;
}

.tips-content {
  padding: 20px 0;

  .tips-list {
    ul {
      margin: 10px 0 0 0;
      padding-left: 20px;

      li {
        margin: 8px 0;
        line-height: 1.6;
      }
    }
  }
}

::v-deep .el-steps {
  margin-bottom: 40px;
}

::v-deep .el-input-group__append {
  padding: 0;
  border: none;

  .el-button {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
}
</style>

