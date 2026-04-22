<template>
  <div class="app-container">
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="用户名">
          <el-input
            v-model.trim="params.username"
            clearable
            placeholder="用户名"
            @keyup.enter.native="search"
            @clear="search"
          />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            :disabled="multipleSelection.length === 0"
            :loading="loading"
            icon="el-icon-delete"
            type="danger"
            @click="batchDelete"
          >批量删除</el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            :loading="syncLoading"
            icon="el-icon-refresh"
            type="success"
            @click="syncExistingUsers"
          >同步已有用户</el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip width="80" sortable prop="ID" label="序号" />
        <el-table-column show-overflow-tooltip sortable prop="username" label="用户名" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="nickname" label="昵称" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="mail" label="邮箱" min-width="150" />
        <el-table-column show-overflow-tooltip sortable prop="uidNumber" label="UID" width="80" />
        <el-table-column show-overflow-tooltip sortable prop="gidNumber" label="GID" width="80" />
        <el-table-column show-overflow-tooltip sortable prop="departments" label="部门" min-width="120" />
        <el-table-column show-overflow-tooltip sortable prop="mobile" label="手机" min-width="120" />
        <el-table-column show-overflow-tooltip sortable prop="jobNumber" label="工号" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="position" label="职位" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="isUsed" label="已使用" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.isUsed ? 'success' : 'info'" size="small">
              {{ scope.row.isUsed ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" width="120">
          <template #default="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
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
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        background
        style="margin-top: 10px; float: right; margin-bottom: 10px"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <!-- 新增/编辑对话框 -->
      <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="700px">
        <el-form ref="dialogForm" :model="dialogFormData" :rules="dialogFormRules" label-width="100px" size="small">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="用户名" prop="username">
                <el-input v-model.trim="dialogFormData.username" placeholder="用户名（拼音）" :disabled="isEdit" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="昵称" prop="nickname">
                <el-input v-model.trim="dialogFormData.nickname" placeholder="昵称/中文名" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="邮箱" prop="mail">
                <el-input v-model.trim="dialogFormData.mail" placeholder="邮箱" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="花名" prop="givenName">
                <el-input v-model.trim="dialogFormData.givenName" placeholder="花名/英文名" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="UID" prop="uidNumber">
                <el-input-number v-model="dialogFormData.uidNumber" :min="1000" :max="65534" placeholder="Unix UID" style="width: 100%" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="GID" prop="gidNumber">
                <el-input-number v-model="dialogFormData.gidNumber" :min="0" :max="65534" placeholder="Unix GID (默认=UID)" style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="部门" prop="departmentId">
                <el-select v-model="dialogFormData.departmentId" placeholder="选择部门" clearable style="width: 100%" @change="onDepartmentChange">
                  <el-option v-for="item in groupList" :key="item.ID" :label="item.groupName" :value="String(item.ID)" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="手机" prop="mobile">
                <el-input v-model.trim="dialogFormData.mobile" placeholder="手机号" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="工号" prop="jobNumber">
                <el-input v-model.trim="dialogFormData.jobNumber" placeholder="工号" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="职位" prop="position">
                <el-input v-model.trim="dialogFormData.position" placeholder="职位" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="主目录" prop="homeDirectory">
                <el-input v-model.trim="dialogFormData.homeDirectory" placeholder="/home/用户名" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="登录Shell" prop="loginShell">
                <el-select v-model="dialogFormData.loginShell" placeholder="登录Shell" style="width: 100%">
                  <el-option label="/bin/bash" value="/bin/bash" />
                  <el-option label="/bin/sh" value="/bin/sh" />
                  <el-option label="/bin/zsh" value="/bin/zsh" />
                  <el-option label="/sbin/nologin" value="/sbin/nologin" />
                  <el-option label="/bin/false" value="/bin/false" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="地址" prop="postalAddress">
                <el-input v-model.trim="dialogFormData.postalAddress" placeholder="地址" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="说明" prop="introduction">
                <el-input v-model.trim="dialogFormData.introduction" type="textarea" :rows="2" placeholder="说明/简介" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="备注" prop="remark">
                <el-input v-model.trim="dialogFormData.remark" type="textarea" :rows="2" placeholder="配置备注" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="dialogVisible = false">取消</el-button>
          <el-button size="mini" type="primary" :loading="submitLoading" @click="submitForm">确定</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { preConfigList, preConfigAdd, preConfigUpdate, preConfigDelete, preConfigSyncUsers } from '@/api/personnel/userPreConfig'
import { groupList } from '@/api/personnel/group'

export default {
  name: 'UserPreConfig',
  data() {
    return {
      loading: false,
      submitLoading: false,
      syncLoading: false,
      tableData: [],
      total: 0,
      params: {
        username: '',
        pageNum: 1,
        pageSize: 10
      },
      multipleSelection: [],
      dialogVisible: false,
      dialogTitle: '新增用户配置',
      isEdit: false,
      groupList: [],
      dialogFormData: {
        id: 0,
        username: '',
        nickname: '',
        givenName: '',
        mail: '',
        uidNumber: null,
        gidNumber: null,
        departmentId: '',
        departments: '',
        mobile: '',
        jobNumber: '',
        position: '',
        postalAddress: '',
        introduction: '',
        homeDirectory: '',
        loginShell: '/bin/bash',
        remark: ''
      },
      dialogFormRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getList()
    this.getGroupList()
  },
  methods: {
    async getList() {
      this.loading = true
      try {
        const res = await preConfigList(this.params)
        this.tableData = res.data.list || []
        this.total = res.data.total || 0
      } catch (e) {
        console.error(e)
      } finally {
        this.loading = false
      }
    },
    async getGroupList() {
      try {
        const res = await groupList({ pageNum: 1, pageSize: 1000 })
        this.groupList = res.data.list || []
      } catch (e) {
        console.error(e)
      }
    },
    search() {
      this.params.pageNum = 1
      this.getList()
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getList()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getList()
    },
    create() {
      this.dialogTitle = '新增用户配置'
      this.isEdit = false
      this.dialogFormData = {
        id: 0,
        username: '',
        nickname: '',
        givenName: '',
        mail: '',
        uidNumber: null,
        gidNumber: null,
        departmentId: '',
        departments: '',
        mobile: '',
        jobNumber: '',
        position: '',
        postalAddress: '',
        introduction: '',
        homeDirectory: '',
        loginShell: '/bin/bash',
        remark: ''
      }
      this.dialogVisible = true
    },
    update(row) {
      this.dialogTitle = '编辑用户配置'
      this.isEdit = true
      this.dialogFormData = {
        id: row.ID,
        username: row.username,
        nickname: row.nickname,
        givenName: row.givenName,
        mail: row.mail,
        uidNumber: row.uidNumber || null,
        gidNumber: row.gidNumber || null,
        departmentId: row.departmentId || '',
        departments: row.departments || '',
        mobile: row.mobile,
        jobNumber: row.jobNumber,
        position: row.position,
        postalAddress: row.postalAddress,
        introduction: row.introduction,
        homeDirectory: row.homeDirectory,
        loginShell: row.loginShell || '/bin/bash',
        remark: row.remark
      }
      this.dialogVisible = true
    },
    onDepartmentChange(val) {
      const group = this.groupList.find(g => String(g.ID) === val)
      if (group) {
        this.dialogFormData.departments = group.groupName
      }
    },
    async submitForm() {
      this.$refs.dialogForm.validate(async(valid) => {
        if (!valid) return

        this.submitLoading = true
        try {
          if (this.isEdit) {
            await preConfigUpdate(this.dialogFormData)
            this.$message.success('更新成功')
          } else {
            await preConfigAdd(this.dialogFormData)
            this.$message.success('添加成功')
          }
          this.dialogVisible = false
          this.getList()
        } catch (e) {
          console.error(e)
        } finally {
          this.submitLoading = false
        }
      })
    },
    async singleDelete(id) {
      this.loading = true
      try {
        await preConfigDelete({ ids: [id] })
        this.$message.success('删除成功')
        this.getList()
      } catch (e) {
        console.error(e)
      } finally {
        this.loading = false
      }
    },
    async batchDelete() {
      this.$confirm('确定要批量删除选中的配置吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async() => {
        this.loading = true
        try {
          const ids = this.multipleSelection.map(item => item.ID)
          await preConfigDelete({ ids })
          this.$message.success('删除成功')
          this.getList()
        } catch (e) {
          console.error(e)
        } finally {
          this.loading = false
        }
      }).catch(() => {})
    },
    // 同步已有用户到预配置
    async syncExistingUsers() {
      this.$confirm('确定要将已有用户同步到预配置列表吗？这将把所有用户（除admin外）添加到预配置中。', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }).then(async() => {
        this.syncLoading = true
        try {
          const res = await preConfigSyncUsers()
          if (res.code === 0 || res.code === 200) {
            this.$message.success(`同步完成: 新增 ${res.data.synced} 条, 跳过 ${res.data.skipped} 条`)
            this.getList()
          } else {
            this.$message.error(res.msg || '同步失败')
          }
        } catch (e) {
          console.error(e)
          this.$message.error('同步失败')
        } finally {
          this.syncLoading = false
        }
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.container-card {
  margin: 10px;
}
</style>

