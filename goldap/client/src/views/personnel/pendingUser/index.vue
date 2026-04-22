<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="用户名">
          <el-input v-model.trim="params.username" style="width: 100px;" clearable placeholder="用户名" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="真实姓名">
          <el-input v-model.trim="params.nickname" style="width: 100px;" clearable placeholder="真实姓名" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model.trim="params.mail" style="width: 150px;" clearable placeholder="邮箱" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model.trim="params.status" style="width: 120px;" clearable placeholder="审核状态" @change="search" @clear="search">
            <el-option label="待审核" :value="0" />
            <el-option label="已通过" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="username" label="用户名" />
        <el-table-column show-overflow-tooltip sortable prop="nickname" label="真实姓名" />
        <el-table-column show-overflow-tooltip sortable prop="mail" label="邮箱" />
        <el-table-column label="审核状态" align="center" width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status === 0" size="small" type="warning">待审核</el-tag>
            <el-tag v-else-if="scope.row.status === 1" size="small" type="success">已通过</el-tag>
            <el-tag v-else-if="scope.row.status === 2" size="small" type="danger">已拒绝</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="remark" label="注册备注" />
        <el-table-column show-overflow-tooltip prop="reviewer" label="审核人" />
        <el-table-column show-overflow-tooltip prop="reviewRemark" label="审核备注" />
        <el-table-column show-overflow-tooltip sortable prop="CreatedAt" label="注册时间" />
        <el-table-column show-overflow-tooltip sortable prop="ReviewedAt" label="审核时间" />
        <el-table-column fixed="right" label="操作" align="center" width="200">
          <template slot-scope="scope">
            <el-tooltip v-if="scope.row.status === 0" content="审核" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-check" circle type="success" @click="review(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <!-- 审核对话框 -->
      <el-dialog title="审核用户注册申请" :visible.sync="reviewDialogVisible" width="60%" :close-on-click-modal="false">
        <el-form ref="reviewForm" size="small" :model="reviewFormData" :rules="reviewFormRules" label-width="120px">
          <el-alert
            title="用户注册信息"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px;"
          />
          <el-row>
            <el-col :span="12">
              <el-form-item label="用户名">
                <el-input v-model.trim="reviewFormData.username" disabled />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="真实姓名">
                <el-input v-model.trim="reviewFormData.nickname" disabled />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="邮箱">
                <el-input v-model.trim="reviewFormData.mail" disabled />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="注册备注">
                <el-input v-model.trim="reviewFormData.remark" type="textarea" :rows="2" disabled />
              </el-form-item>
            </el-col>
          </el-row>

          <el-divider content-position="left">审核设置</el-divider>
          <el-row>
            <el-col :span="12">
              <el-form-item label="审核结果" prop="status">
                <el-radio-group v-model="reviewFormData.status">
                  <el-radio :label="1">通过</el-radio>
                  <el-radio :label="2">拒绝</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="审核备注" prop="reviewRemark">
                <el-input v-model.trim="reviewFormData.reviewRemark" type="textarea" :rows="3" placeholder="请输入审核备注" maxlength="255" show-word-limit />
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 审核通过时的配置 -->
          <div v-if="reviewFormData.status === 1">
            <el-divider content-position="left">用户配置</el-divider>
            <el-row>
              <el-col :span="24">
                <el-form-item label="分配分组" prop="departmentId">
                  <treeselect
                    v-model="reviewFormData.departmentId"
                    :options="departmentsOptions"
                    placeholder="请选择分组（必选）"
                    :normalizer="normalizer"
                    value-consists-of="ALL"
                    :multiple="true"
                    :flat="true"
                    no-children-text="没有更多选项"
                    no-results-text="没有匹配的选项"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="分配角色" prop="roleIds">
                  <el-select v-model.trim="reviewFormData.roleIds" multiple placeholder="请选择角色" style="width:100%">
                    <el-option
                      v-for="item in roles"
                      :key="item.ID"
                      :label="item.name"
                      :value="item.ID"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-divider content-position="left">Unix用户属性</el-divider>
            <el-row>
              <el-col :span="12">
                <el-form-item label="UID号" prop="uidNumber">
                  <el-input-number v-model="reviewFormData.uidNumber" :min="1000" :max="65534" placeholder="Unix用户ID（留空自动分配）" style="width:100%" />
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">留空将自动分配，建议从1000开始</div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="GID号" prop="gidNumber">
                  <el-input-number v-model="reviewFormData.gidNumber" :min="100" :max="65534" placeholder="Unix组ID" style="width:100%" />
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">默认值：UID号</div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="主目录" prop="homeDirectory">
                  <el-input v-model.trim="reviewFormData.homeDirectory" placeholder="用户主目录路径" />
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">默认：/home/用户名</div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="登录Shell" prop="loginShell">
                  <el-select v-model.trim="reviewFormData.loginShell" placeholder="登录Shell" style="width:100%">
                    <el-option label="/bin/bash" value="/bin/bash" />
                    <el-option label="/bin/sh" value="/bin/sh" />
                    <el-option label="/bin/zsh" value="/bin/zsh" />
                    <el-option label="/sbin/nologin" value="/sbin/nologin" />
                    <el-option label="/usr/sbin/nologin" value="/usr/sbin/nologin" />
                  </el-select>
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">默认：/bin/bash</div>
                </el-form-item>
              </el-col>
            </el-row>

            <el-divider content-position="left">权限配置</el-divider>
            <el-row>
              <el-col :span="12">
                <el-form-item label="允许sudo">
                  <el-switch v-model="reviewFormData.allowSudo" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="允许密钥上传">
                  <el-switch v-model="reviewFormData.allowSSHKey" />
                </el-form-item>
              </el-col>
              <el-col v-if="reviewFormData.allowSudo" :span="24">
                <el-form-item label="Sudo规则" prop="sudoRules">
                  <el-input v-model.trim="reviewFormData.sudoRules" type="textarea" :rows="3" placeholder='请输入sudo规则，JSON格式，如：{"option": "!authenticate"}' />
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">留空则使用默认规则（无需密码）</div>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelReview">取 消</el-button>
          <el-button size="mini" :loading="reviewLoading" type="primary" @click="submitReview">确 定</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { getPendingUsers, reviewPendingUser, deletePendingUsers } from '@/api/system/user'
import { getGroupTree } from '@/api/personnel/group'
import { getRoles } from '@/api/system/role'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'

export default {
  name: 'PendingUser',
  components: {
    Treeselect
  },
  data() {
    return {
      loading: false,
      reviewLoading: false,
      params: {
        username: '',
        nickname: '',
        mail: '',
        status: 0, // 默认只显示待审核的
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      multipleSelection: [],
      reviewDialogVisible: false,
      reviewFormData: {
        id: null,
        username: '',
        nickname: '',
        mail: '',
        remark: '',
        status: 1, // 1通过, 2拒绝
        reviewRemark: '',
        departmentId: [],
        roleIds: [],
        uidNumber: 0,
        gidNumber: 0,
        homeDirectory: '',
        loginShell: '/bin/bash',
        allowSudo: false,
        allowSSHKey: true,
        sudoRules: ''
      },
      reviewFormRules: {
        status: [{ required: true, message: '请选择审核结果', trigger: 'change' }],
        departmentId: [
          {
            validator: (rule, value, callback) => {
              if (this.reviewFormData.status === 1 && (!value || value.length === 0)) {
                callback(new Error('审核通过时必须分配分组'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ]
      },
      departmentsOptions: [],
      roles: []
    }
  },
  created() {
    this.search()
    this.loadDepartments()
    this.loadRoles()
  },
  methods: {
    normalizer(node) {
      if (node.children && !node.children.length) {
        delete node.children
      }
      return {
        id: node.ID,
        label: node.groupName,
        // 只禁用root分组，允许选择ou类型的分组（如CMPLabHPC）
        isDisabled: node.groupName === 'root' || (node.groupType === 'T' && node.ID === 0),
        children: node.children
      }
    },
    async loadDepartments() {
      try {
        const res = await getGroupTree({})
        if (res.code === 200) {
          this.departmentsOptions = res.data || []
        }
      } catch (error) {
        this.$message.error('加载分组列表失败')
      }
    },
    async loadRoles() {
      try {
        const res = await getRoles({ pageNum: 1, pageSize: 1000 })
        if (res.code === 200) {
          this.roles = res.data.roles || []
        }
      } catch (error) {
        this.$message.error('加载角色列表失败')
      }
    },
    async search() {
      this.loading = true
      try {
        const res = await getPendingUsers(this.params)
        if (res.code === 200) {
          this.tableData = res.data.pendingUsers || []
          this.total = res.data.total || 0
        }
      } catch (error) {
        this.$message.error('获取待审核用户列表失败')
      } finally {
        this.loading = false
      }
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    handleSizeChange(val) {
      this.params.pageSize = val
      this.search()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.search()
    },
    review(row) {
      this.reviewFormData = {
        id: row.ID,
        username: row.username,
        nickname: row.nickname,
        mail: row.mail,
        remark: row.remark,
        status: 1,
        reviewRemark: '',
        departmentId: [],
        roleIds: [2], // 默认普通用户角色
        uidNumber: 0,
        gidNumber: 0,
        homeDirectory: '',
        loginShell: '/bin/bash',
        allowSudo: false,
        allowSSHKey: true,
        sudoRules: ''
      }
      
      // 查找CMPLabHPC分组并设置为默认值
      const findCMPLabHPC = (groups) => {
        for (const group of groups) {
          if (group.groupName === 'CMPLabHPC') {
            return group.ID
          }
          if (group.children && group.children.length > 0) {
            const found = findCMPLabHPC(group.children)
            if (found) return found
          }
        }
        return null
      }
      const cmplabHpcId = findCMPLabHPC(this.departmentsOptions || [])
      if (cmplabHpcId) {
        this.reviewFormData.departmentId = [cmplabHpcId]
      }
      
      this.reviewDialogVisible = true
    },
    cancelReview() {
      this.reviewDialogVisible = false
      this.$refs.reviewForm.resetFields()
    },
    async submitReview() {
      this.$refs.reviewForm.validate(async (valid) => {
        if (valid) {
          this.reviewLoading = true
          try {
            const res = await reviewPendingUser(this.reviewFormData)
            if (res.code === 200) {
              this.$message.success('审核成功')
              this.reviewDialogVisible = false
              this.search()
            }
          } catch (error) {
            this.$message.error(error.message || '审核失败')
          } finally {
            this.reviewLoading = false
          }
        }
      })
    },
    async singleDelete(id) {
      try {
        const res = await deletePendingUsers({ ids: [id] })
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.search()
        }
      } catch (error) {
        this.$message.error('删除失败')
      }
    },
    async batchDelete() {
      if (this.multipleSelection.length === 0) {
        this.$message.warning('请选择要删除的记录')
        return
      }
      try {
        const ids = this.multipleSelection.map(item => item.ID)
        const res = await deletePendingUsers({ ids })
        if (res.code === 200) {
          this.$message.success('批量删除成功')
          this.search()
        }
      } catch (error) {
        this.$message.error('批量删除失败')
      }
    }
  }
}
</script>

<style scoped>
.container-card {
  margin: 10px;
}
</style>

