<template>
  <div>
    <el-card class="container-card" shadow="always">
      <div class="config-header">
        <h2>
          <svg-icon icon-class="setting" style="margin-right: 8px;" />
          系统维护
        </h2>
        <p class="config-subtitle">管理系统功能开关和维护提示信息</p>
      </div>

      <el-form ref="configForm" :model="configForm" label-width="200px" style="margin-top: 30px;">
        <el-form-item label="密码自助服务">
          <el-switch
            v-model="configForm.passwordSelfServiceEnabled"
            active-text="开启"
            inactive-text="关闭"
            active-color="#13ce66"
            inactive-color="#ff4949"
          />
          <div class="config-tip">控制用户是否可以使用密码自助服务功能（忘记密码、修改密码等）</div>
        </el-form-item>

        <el-form-item label="禁止注册">
          <el-switch
            v-model="configForm.registrationDisabled"
            active-text="禁止"
            inactive-text="允许"
            active-color="#ff4949"
            inactive-color="#13ce66"
          />
          <div class="config-tip">开启后将禁止新用户注册，已提交的注册申请将无法通过</div>
        </el-form-item>

        <el-form-item label="系统维护模式">
          <el-switch
            v-model="configForm.systemMaintenanceMode"
            active-text="维护中"
            inactive-text="正常"
            active-color="#e6a23c"
            inactive-color="#13ce66"
          />
          <div class="config-tip">开启后将在密码自助服务页面显示系统维护提示</div>
        </el-form-item>

        <el-form-item label="维护提示信息">
          <el-input
            v-model="configForm.maintenanceMessage"
            type="textarea"
            :rows="3"
            placeholder="请输入系统维护提示信息"
            maxlength="200"
            show-word-limit
            style="width: 500px;"
          />
          <div class="config-tip">当系统维护模式开启时，此信息将显示在密码自助服务页面</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="submitForm">保存配置</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { getSystemConfig, updateSystemConfig } from '@/api/system/systemConfig'
import { Message } from 'element-ui'

export default {
  name: 'SystemConfig',
  data() {
    return {
      configForm: {
        passwordSelfServiceEnabled: true,
        registrationDisabled: false,
        systemMaintenanceMode: false,
        maintenanceMessage: '系统正在升级维护中，请稍后再试'
      },
      submitLoading: false
    }
  },
  created() {
    this.loadConfig()
  },
  methods: {
    async loadConfig() {
      try {
        const res = await getSystemConfig()
        if (res.code === 200 && res.data) {
          this.configForm = {
            passwordSelfServiceEnabled: res.data.passwordSelfServiceEnabled !== false,
            registrationDisabled: res.data.registrationDisabled === true,
            systemMaintenanceMode: res.data.systemMaintenanceMode === true,
            maintenanceMessage: res.data.maintenanceMessage || '系统正在升级维护中，请稍后再试'
          }
        }
      } catch (error) {
        console.error('Failed to load system config:', error)
        Message({
          message: '加载系统配置失败',
          type: 'error'
        })
      }
    },
    async submitForm() {
      this.submitLoading = true
      try {
        const res = await updateSystemConfig(this.configForm)
        if (res.code === 200) {
          Message({
            message: '系统配置保存成功',
            type: 'success'
          })
        }
      } catch (error) {
        Message({
          message: error.message || '保存系统配置失败',
          type: 'error'
        })
      } finally {
        this.submitLoading = false
      }
    },
    resetForm() {
      this.loadConfig()
    }
  }
}
</script>

<style lang="scss" scoped>
.config-header {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;

  h2 {
    margin: 0 0 10px 0;
    font-size: 20px;
    color: #303133;
    display: flex;
    align-items: center;
  }

  .config-subtitle {
    margin: 0;
    color: #909399;
    font-size: 14px;
  }
}

.config-tip {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

::v-deep .el-form-item__label {
  font-weight: 500;
}
</style>

