<template>
  <div class="password-service-container">
    <div class="service-card">
      <div class="header">
        <h2 class="title">
          <svg-icon icon-class="lock" style="margin-right: 8px;" />
          忘记密码
        </h2>
        <p class="subtitle">忘记密码服务</p>
      </div>

      <el-tabs v-model="activeTab" class="service-tabs">
        <!-- 忘记密码 -->
        <el-tab-pane label="忘记密码" name="forgot">
          <div class="tab-content">
            <el-steps :active="forgotStep" finish-status="success" align-center>
              <el-step title="验证身份" description="通过邮箱验证"></el-step>
              <el-step title="重置密码" description="设置新密码"></el-step>
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
        <el-tab-pane label="修改密码" name="change" :disabled="!isLoggedIn">
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
  </div>
</template>

<script>
import { emailPass, sendCode } from '@/api/system/user'
import { changePwd } from '@/api/system/user'
import { validEmail } from '@/utils/validate'
import { Message } from 'element-ui'
import JSEncrypt from 'jsencrypt'
import store from '@/store'

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
    const validatePassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入密码'))
      } else if (value.length < 6) {
        callback(new Error('密码长度不能少于6位'))
      } else {
        callback()
      }
    }

    return {
      activeTab: 'forgot',
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
      publicKey: process.env.VUE_APP_PUBLIC_KEY
    }
  },
  computed: {
    isLoggedIn() {
      return !!store.getters.token
    }
  },
  watch: {
    activeTab() {
      this.resetAllForms()
    }
  },
  methods: {
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
    // 检查密码强度（仅用于修改密码功能）
    checkPasswordStrength() {
      const password = this.changeForm.newPassword

      if (!password) {
        this.passwordStrength = { level: 'weak', percent: 0, text: '弱' }
        return
      }

      // 检查各项要求
      this.passwordChecks = {
        length: password.length >= 6,
        uppercase: /[A-Z]/.test(password),
        lowercase: /[a-z]/.test(password),
        number: /[0-9]/.test(password),
        special: /[!@#$%^&*(),.?":{}|<>]/.test(password)
      }

      // 计算强度
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
    resetForgotForm() {
      this.$refs.forgotForm?.resetFields()
      this.forgotForm = { mail: '', code: '' }
      this.forgotStep = 0
    },
    resetChangeForm() {
      this.$refs.changeForm?.resetFields()
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
    resetAllForms() {
      this.resetForgotForm()
      this.resetChangeForm()
      this.forgotStep = 0
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
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.service-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 700px;
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

.password-tips {
  margin-top: 10px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;

  p {
    margin: 0 0 8px 0;
    font-weight: 500;
    color: #333;
  }

  ul {
    margin: 0;
    padding-left: 20px;
    color: #666;

    li {
      margin: 4px 0;
      transition: color 0.3s;

      &.valid {
        color: #67c23a;
      }
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

