<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="名称">
          <el-input style="width: 100px;" v-model.trim="params.groupName" clearable placeholder="名称" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input style="width: 100px;" v-model.trim="params.remark" clearable placeholder="描述" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
          <el-form-item label="同步状态">
          <el-select style="width: 110px;" v-model.trim="params.syncState" clearable placeholder="同步状态" @change="search" @clear="search">
            <el-option label="已同步" value="1" />
            <el-option label="未同步" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <!-- <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="resetData">重置</el-button>
        </el-form-item> -->
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
        <el-form-item>
          <el-button  :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-upload2" type="success" @click="batchSync">批量同步</el-button>
        </el-form-item>
        <br>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-download" type="warning" @click="syncOpenLdapDepts">同步原ldap部门</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :default-expand-all="true" :tree-props="{children: 'children', hasChildren: 'hasChildren'}" row-key="ID" :data="infoTableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="groupName" label="名称" />
        <el-table-column show-overflow-tooltip sortable prop="groupType" label="类型" width="80" />
        <el-table-column show-overflow-tooltip sortable prop="gidNumber" label="GID" width="80">
          <template #default="scope">
            <span v-if="scope.row.groupType === 'posix'">{{ scope.row.gidNumber }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="groupDn" label="DN" />
        <el-table-column show-overflow-tooltip sortable prop="remark" label="描述" />
        <el-table-column show-overflow-tooltip sortable prop="ipRanges" label="IP范围" />
        <el-table-column show-overflow-tooltip sortable prop="CreatedAt" label="创建时间" />
        <el-table-column show-overflow-tooltip sortable prop="UpdatedAt" label="更新时间" />
        <el-table-column fixed="right" label="操作" align="center" width="280">
          <template #default="scope">
            <el-tooltip v-if="scope.row.groupName !== 'sudoers'" content="成员管理" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-user" circle type="info" @click="manageMembers(scope.row)" />
            </el-tooltip>
            <el-tooltip v-if="scope.row.groupType === 'ou' && scope.row.groupName !== 'sudoers'" content="管理用户权限" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-user-solid" circle type="warning" @click="managePermissions(scope.row)" />
            </el-tooltip>
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip v-if="scope.row.groupName !== 'sudouser-nopasswd' && scope.row.groupName !== 'sudouser-other' && scope.row.groupName !== 'sudoers' && scope.row.groupName !== 'docker'" class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip v-else content="默认sudoers分组不允许删除" effect="dark" placement="top">
              <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" disabled />
            </el-tooltip>
            <el-tooltip v-if="scope.row.syncState == 2" class="delete-popover" content="同步" effect="dark" placement="top">
              <el-popconfirm title="确定同步吗？" @onConfirm="singleSync(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-upload2" circle type="success" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <!-- 新增 -->
      <el-dialog :title="dialogFormTitle" :visible.sync="updateLoading">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="名称" prop="groupName">
            <el-input v-model.trim="dialogFormData.groupName" placeholder="名称(拼音)" />
          </el-form-item>
          <el-form-item label="分组类型" prop="groupType">
            <el-select v-model.trim="dialogFormData.groupType" placeholder="建议仅第一层为ou，如果不确定，就用cn" style="width:100%">
              <el-option label="cn[分组]" value="cn" />
              <el-option label="ou[组织]" value="ou" />
              <el-option label="posix[Linux系统组]" value="posix" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="dialogFormData.groupType === 'posix'" label="GID号" prop="gidNumber">
            <el-input-number v-model="dialogFormData.gidNumber" :min="100" :max="65535" placeholder="Linux系统组GID" style="width:100%" />
            <div style="font-size: 12px; color: #909399; margin-top: 4px;">
              用于Linux系统组的GID号，如docker组通常为984
            </div>
          </el-form-item>
          <el-form-item label="上级分组" prop="parentId">
            <treeselect
              v-model="dialogFormData.parentId"
              :options="treeselectData"
              :normalizer="normalizer"
              placeholder="请选择上级分组"
              @input="treeselectInput"
            />
          </el-form-item>
          <el-form-item label="描述" prop="remark">
            <el-input v-model.trim="dialogFormData.remark" type="textarea" placeholder="描述" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
          </el-form-item>
          <el-form-item label="IP范围" prop="ipRanges">
            <el-input v-model.trim="dialogFormData.ipRanges" type="textarea" placeholder='请输入IP范围，JSON格式，如：["192.168.11.6-192.168.11.8","192.168.11.2"]' :autosize="{minRows: 3, maxRows: 6}" />
            <div style="font-size: 12px; color: #909399; margin-top: 4px;">
              支持单个IP（如：192.168.11.2）或IP范围（如：192.168.11.6-192.168.11.8），多个IP用JSON数组格式
            </div>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
      <!-- 编辑 -->
      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="名称" prop="groupName">
            <el-input v-model.trim="dialogFormData.groupName" :disabled="true" placeholder="名称" />
          </el-form-item>
          <el-form-item v-if="dialogFormData.groupType === 'posix'" label="GID号" prop="gidNumber">
            <el-input-number v-model="dialogFormData.gidNumber" :min="100" :max="65535" placeholder="Linux系统组GID" style="width:100%" />
            <div style="font-size: 12px; color: #909399; margin-top: 4px;">
              用于Linux系统组的GID号，如docker组通常为984
            </div>
          </el-form-item>
          <el-form-item label="描述" prop="remark">
            <el-input v-model.trim="dialogFormData.remark" type="textarea" placeholder="描述" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
          </el-form-item>
          <el-form-item label="IP范围" prop="ipRanges">
            <el-input v-model.trim="dialogFormData.ipRanges" type="textarea" placeholder='请输入IP范围，JSON格式，如：["192.168.11.6-192.168.11.8","192.168.11.2"]' :autosize="{minRows: 3, maxRows: 6}" />
            <div style="font-size: 12px; color: #909399; margin-top: 4px;">
              支持单个IP（如：192.168.11.2）或IP范围（如：192.168.11.6-192.168.11.8），多个IP用JSON数组格式
            </div>
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
      <!-- 管理分组用户权限 -->
      <el-dialog title="管理分组用户权限" :visible.sync="permissionDialogVisible" width="80%">
        <el-table v-loading="permissionLoading" :data="permissionTableData" border stripe style="width: 100%">
          <el-table-column prop="username" label="用户名" width="120" />
          <el-table-column prop="nickname" label="真实姓名" width="120" />
          <el-table-column label="允许Sudo" width="120" align="center">
            <template #default="scope">
              <el-switch
                v-model="scope.row.allowSudo"
                @change="updatePermission(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="允许SSH密钥" width="140" align="center">
            <template #default="scope">
              <el-switch
                v-model="scope.row.allowSSHKey"
                @change="updatePermission(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="Sudo规则" min-width="200">
            <template #default="scope">
              <el-input
                v-model="scope.row.sudoRules"
                type="textarea"
                :rows="2"
                placeholder="JSON格式的sudo规则，如：!authenticate"
                @blur="updatePermission(scope.row)"
              />
            </template>
          </el-table-column>
        </el-table>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="permissionDialogVisible = false">关 闭</el-button>
        </div>
      </el-dialog>
      <!-- 成员管理对话框 -->
      <el-dialog title="成员管理" :visible.sync="memberDialogVisible" width="70%">
        <el-row :gutter="20">
          <el-col :span="12">
            <div style="margin-bottom: 10px;">
              <el-input
                v-model="memberSearchKeyword"
                placeholder="搜索用户"
                prefix-icon="el-icon-search"
                clearable
                @input="filterUsers"
              />
            </div>
            <div style="border: 1px solid #DCDFE6; border-radius: 4px; padding: 10px; max-height: 400px; overflow-y: auto;">
              <div style="font-weight: bold; margin-bottom: 10px; color: #409EFF;">可选用户（{{ filteredAvailableUsers.length }}）</div>
              <el-table
                v-loading="memberLoading"
                :data="filteredAvailableUsers"
                border
                stripe
                height="350"
                @selection-change="handleAvailableUserSelection"
              >
                <el-table-column type="selection" width="55" />
                <el-table-column prop="userName" label="用户名" width="120" />
                <el-table-column prop="nickName" label="真实姓名" width="120" />
                <el-table-column prop="mail" label="邮箱" />
              </el-table>
            </div>
          </el-col>
          <el-col :span="12">
            <div style="margin-bottom: 10px;">
              <el-button type="primary" size="small" :disabled="selectedAvailableUsers.length === 0" @click="addMembers">
                <i class="el-icon-arrow-right" /> 添加到分组
              </el-button>
              <el-button type="danger" size="small" :disabled="selectedGroupMembers.length === 0" @click="removeMembers">
                <i class="el-icon-arrow-left" /> 从分组移除
              </el-button>
            </div>
            <div style="border: 1px solid #DCDFE6; border-radius: 4px; padding: 10px; max-height: 400px; overflow-y: auto;">
              <div style="font-weight: bold; margin-bottom: 10px; color: #67C23A;">分组成员（{{ groupMembers.length }}）</div>
              <el-table
                v-loading="memberLoading"
                :data="groupMembers"
                border
                stripe
                height="350"
                @selection-change="handleGroupMemberSelection"
              >
                <el-table-column type="selection" width="55" />
                <el-table-column prop="userName" label="用户名" width="120" />
                <el-table-column prop="nickName" label="真实姓名" width="120" />
                <el-table-column prop="mail" label="邮箱" />
              </el-table>
            </div>
          </el-col>
        </el-row>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="memberDialogVisible = false">关 闭</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'
import { getGroupTree,  groupAdd, groupUpdate, groupDel, syncOpenLdapDeptsApi, syncSqlGroups, useGroupList, useGroupRole, groupInfo, delGroup } from '@/api/personnel/group'
import { getGroupUserPermissions, addGroupUserPermission, updateGroupUserPermission } from '@/api/personnel/groupUserPermission'
import { Message } from 'element-ui'

export default {
  name: 'Group',
  components: {
    Treeselect
  },
  filters: {
    methodTagFilter(val) {
      if (val === 'GET') {
        return ''
      } else if (val === 'POST') {
        return 'success'
      } else {
        return 'info'
      }
    }
  },
  data() {
    return {
      // 查询参数
      params: {
        groupName: undefined,
        remark: undefined,
        syncState: undefined,
        pageNum: 1,
        pageSize: 1000// 平常百姓人家应该不会有这么多数据吧,后台限制最大单次获取1000条
      },
      // 表格数据
      tableData: [],
      infoTableData: [],
      total: 0,
      loading: false,
      // 上级目录数据
      treeselectData: [],
      treeselectValue: 0,
      updateLoading: false, // 新增
      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        ID: '',
        groupName: '',
        parentId: 0,
        syncState:1,
        groupType: '',
        gidNumber: 0,
        remark: '',
        ipRanges: ''
      },
      dialogFormRules: {

        groupName: [
          { required: true, message: '请输入所属类别', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        groupType: [
          { required: true, message: '请输入分组类型', trigger: 'blur' },
          { min: 1, max: 50, message: 'ou、cn或者其他', trigger: 'blur' }
        ],
        parentId: [
          { required: true, message: '请选择父级', trigger: 'blur' },
          { validator: (rule, value, callBack) => {
            if (value >= 0) {
              callBack()
            } else {
              callBack('请选择有效的部门')
            }
          } }
        ],
        remark: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 100, message: '长度在 0 到 100 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      dialogTransfer: '', // 穿梭框头部
      dialogTransferVisible: false,

      transParams: {
        groupId: '',
        nickname: ''
      },
      renderFunc(h, option) {
        return <span>{option.key} - {option.label}</span>
      },
      userArrInfo: [], // 初始人员列表数据
      data: [], // 转化后人员列表数据
      value3: [], // 右侧默认人员列表数据
      userId: [], // 送到后台 -> 勾选的数据code数组
      ui: {
        submitLoading: false
      },
      statusTrans: '',
      // 权限管理相关
      permissionDialogVisible: false,
      permissionLoading: false,
      permissionTableData: [],
      currentGroupId: null,
      // 成员管理相关
      memberDialogVisible: false,
      memberLoading: false,
      groupMembers: [],
      availableUsers: [],
      filteredAvailableUsers: [],
      memberSearchKeyword: '',
      selectedAvailableUsers: [],
      selectedGroupMembers: [],
      currentMemberGroupId: null
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    // // 查询
    search() {
      // 初始化表格数据
      this.infoTableData = JSON.parse(JSON.stringify(this.tableData))
      this.infoTableData = this.deal(this.infoTableData, node => node.groupName.includes(this.params.groupName) || node.remark.includes(this.params.remark)  || node.syncState.toString().includes(this.params.syncState))
    },
    resetData() {
      this.infoTableData = JSON.parse(JSON.stringify(this.tableData))
    },
    // 页面数据过滤
    deal(nodes, predicate) {
      // 如果已经没有节点了，结束递归
      if (!(nodes && nodes.length)) {
        return []
      }
      const newChildren = []
      for (const node of nodes) {
        if (predicate(node)) {
          // 如果节点符合条件，直接加入新的节点集
          newChildren.push(node)
          node.children = this.deal(node.children, predicate)
        } else {
          // 如果当前节点不符合条件，递归过滤子节点，
          // 把符合条件的子节点提升上来，并入新节点集
          newChildren.push(...this.deal(node.children, predicate))
        }
      }
      return newChildren
    },
    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getGroupTree(this.params)
        this.tableData = data
        this.infoTableData = JSON.parse(JSON.stringify(data))
        this.treeselectData = [{ ID: 0, groupName: '顶级类目', children: data }]
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增分组'
      this.updateLoading = true // 新增的展示
      this.dialogType = 'create'
    },
    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.groupName = row.groupName
      this.dialogFormData.groupType = row.groupType
      this.dialogFormData.gidNumber = row.gidNumber || 0
      this.dialogFormData.remark = row.remark
      this.dialogFormData.ipRanges = row.ipRanges || ''
      this.dialogFormTitle = '修改分组'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },
    // 穿梭框
    addUp(row) {
      this.dialogTransfer = '用户管理'
      this.dialogTransferVisible = true
      this.transParams.groupId = row.ID
      this.transParams.nickname = row.remark
      this.$router.push({ path: '/userList', query: row })
    },

    // 判断结果
    judgeResult(res){
      if (res.code==0){
          Message({
            showClose: true,
            message: "操作成功",
            type: 'success'
          })
        }
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true
          try {
            if (this.dialogType === 'create') {
              await groupAdd(this.dialogFormData).then(res =>{
                this.judgeResult(res)
              })
            } else {
              await groupUpdate(this.dialogFormData).then(res =>{
                this.judgeResult(res)
              })
            }
          } finally {
            this.submitLoading = false
          }
          this.resetForm()
          this.getTableData()
        } else {
          Message({
            showClose: true,
            message: '表单校验失败',
            type: 'warn'
          })
          return false
        }
      })
    },

    // 提交表单
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.updateLoading = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        ID: '',
        groupName: '',
        parentId: 0,
        syncState:1,
        groupType: '',
        gidNumber: 0,
        remark: '',
        ipRanges: ''
      }
    },

    // 批量删除
    batchDelete() {
      // 检查是否包含默认sudoers分组
      const protectedGroups = ['sudouser-nopasswd', 'sudouser-other', 'sudoers']
      const hasProtectedGroup = this.multipleSelection.some(x => protectedGroups.includes(x.groupName))
      if (hasProtectedGroup) {
        Message({
          showClose: true,
          type: 'warning',
          message: '不允许删除默认的sudoers分组（sudouser-nopasswd、sudouser-other、sudoers）'
        })
        return
      }
      
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const groupIds = []
        this.multipleSelection.forEach(x => {
          groupIds.push(x.ID)
        })
        try {
          await groupDel({ groupIds: groupIds }).then(res => {
            this.judgeResult(res)
          })
        } finally {
          this.loading = false
        }
        this.getTableData()
      }).catch(() => {
        Message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    // 批量同步
    batchSync() {
      this.$confirm('此操作批量同步数据到Ldap, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const groupIds = []
        this.multipleSelection.forEach(x => {
          groupIds.push(x.ID)
        })
        try {
          await syncSqlGroups({ groupIds: groupIds }).then(res => {
            this.judgeResult(res)
          })
        } finally {
          this.loading = false
        }
        this.getTableData()
      }).catch(() => {
        Message({
          showClose: true,
          type: 'info',
          message: '已取消同步'
        })
      })
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      try {
        await groupDel({ groupIds: [Id] }).then(res =>{
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },
    // 单个同步
    async singleSync(Id) {
      this.loading = true
      try {
        await syncSqlGroups({ groupIds: [Id] }).then(res =>{
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    },
    // treeselect
    normalizer(node) {
      return {
        id: node.ID,
        label: node.groupName,
        children: node.children
      }
    },
    treeselectInput(value) {
      this.treeselectValue = value
    },
    syncOpenLdapDepts() {
      this.loading = true
      syncOpenLdapDeptsApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    // 管理分组用户权限
    async managePermissions(row) {
      this.currentGroupId = row.ID
      this.permissionDialogVisible = true
      this.permissionLoading = true
      try {
        // 获取分组内的用户列表
        const userRes = await useGroupList({ groupId: row.ID })
        const users = userRes.data.userList || []
        
        // 获取分组用户权限列表
        const permRes = await getGroupUserPermissions({ groupId: row.ID })
        const permissions = permRes.data || []
        
        // 创建权限映射
        const permMap = {}
        permissions.forEach(perm => {
          permMap[perm.userId] = perm
        })
        
        // 合并用户和权限数据
        this.permissionTableData = users.map(user => {
          const perm = permMap[user.ID] || {}
          return {
            id: perm.id || null,
            userId: user.ID,
            username: user.username,
            nickname: user.nickname,
            allowSudo: perm.allowSudo || false,
            allowSSHKey: perm.allowSSHKey !== undefined ? perm.allowSSHKey : true,
            sudoRules: perm.sudoRules || ''
          }
        })
      } catch (error) {
        Message.error('加载权限数据失败')
      } finally {
        this.permissionLoading = false
      }
    },
    // 更新权限
    async updatePermission(row) {
      try {
        if (row.id) {
          // 更新现有权限
          await updateGroupUserPermission({
            id: row.id,
            allowSudo: row.allowSudo,
            allowSSHKey: row.allowSSHKey,
            sudoRules: row.sudoRules
          })
          Message.success('权限更新成功')
        } else {
          // 创建新权限
          const res = await addGroupUserPermission({
            groupId: this.currentGroupId,
            userId: row.userId,
            allowSudo: row.allowSudo,
            allowSSHKey: row.allowSSHKey,
            sudoRules: row.sudoRules
          })
          row.id = res.data.id
          Message.success('权限创建成功')
        }
      } catch (error) {
        Message.error('更新权限失败')
      }
    },
    // 管理成员
    async manageMembers(row) {
      this.currentMemberGroupId = row.ID
      this.memberDialogVisible = true
      this.memberLoading = true
      this.memberSearchKeyword = ''
      this.selectedAvailableUsers = []
      this.selectedGroupMembers = []
      try {
        // 获取分组内的用户
        const inGroupRes = await useGroupList({ groupId: row.ID })
        this.groupMembers = inGroupRes.data.userList || []
        
        // 获取不在分组内的用户
        const notInGroupRes = await useGroupRole({ groupId: row.ID })
        this.availableUsers = notInGroupRes.data.userList || []
        this.filteredAvailableUsers = [...this.availableUsers]
      } catch (error) {
        Message.error('加载成员数据失败')
      } finally {
        this.memberLoading = false
      }
    },
    // 过滤用户
    filterUsers() {
      if (!this.memberSearchKeyword) {
        this.filteredAvailableUsers = [...this.availableUsers]
        return
      }
      const keyword = this.memberSearchKeyword.toLowerCase()
      this.filteredAvailableUsers = this.availableUsers.filter(user => {
        const userName = (user.userName || user.username || '').toLowerCase()
        const nickName = (user.nickName || user.nickname || '').toLowerCase()
        const mail = (user.mail || '').toLowerCase()
        return userName.includes(keyword) ||
               nickName.includes(keyword) ||
               mail.includes(keyword)
      })
    },
    // 处理可选用户选择
    handleAvailableUserSelection(selection) {
      this.selectedAvailableUsers = selection
    },
    // 处理分组成员选择
    handleGroupMemberSelection(selection) {
      this.selectedGroupMembers = selection
    },
    // 添加成员到分组
    async addMembers() {
      if (this.selectedAvailableUsers.length === 0) {
        Message.warning('请选择要添加的用户')
        return
      }
      try {
        const userIds = this.selectedAvailableUsers.map(user => user.userId || user.ID)
        await groupInfo({
          groupId: this.currentMemberGroupId,
          userIds: userIds
        })
        Message.success('添加成员成功')
        // 刷新数据
        await this.manageMembers({ ID: this.currentMemberGroupId })
      } catch (error) {
        Message.error('添加成员失败')
      }
    },
    // 从分组移除成员
    async removeMembers() {
      if (this.selectedGroupMembers.length === 0) {
        Message.warning('请选择要移除的用户')
        return
      }
      try {
        const userIds = this.selectedGroupMembers.map(user => user.userId || user.ID)
        await delGroup({
          groupId: this.currentMemberGroupId,
          userIds: userIds
        })
        Message.success('移除成员成功')
        // 刷新数据
        await this.manageMembers({ ID: this.currentMemberGroupId })
      } catch (error) {
        Message.error('移除成员失败')
      }
    }
  }
}
</script>

<style scoped>
  .container-card{
    margin: 10px;
    margin-bottom: 100px;
  }

  .delete-popover{
    margin-left: 10px;
  }
   .transfer-footer {
    margin-left: 20px;
    padding: 6px 5px;
  }
</style>
