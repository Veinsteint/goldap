<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <!-- 用户管理标签页 -->
        <el-tab-pane label="用户管理" name="users">
          <div class="tab-content">
            <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="用户名">
          <el-input v-model.trim="params.username" style="width: 100px;" clearable placeholder="用户名" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model.trim="params.nickname" style="width: 100px;" clearable placeholder="昵称" @keyup.enter.native="search" @clear="search" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model.trim="params.status" style="width: 100px;" clearable placeholder="状态" @change="search" @clear="search">
            <el-option label="正常" value="1" />
            <el-option label="禁用" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="同步状态">
          <el-select v-model.trim="params.syncState" style="width: 100px;" clearable placeholder="同步状态" @change="search" @clear="search">
            <el-option label="已同步" value="1" />
            <el-option label="未同步" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-upload2" type="success" @click="batchSync">批量同步</el-button>
        </el-form-item>
        <br>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-download" type="warning" @click="syncOpenLdapUsers">同步原ldap用户信息</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-upload2" type="primary" @click="syncAllUsersToLdap">同步所有MySQL用户到LDAP</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="username" label="用户名" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="nickname" label="中文名" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="givenName" label="花名" min-width="80" />
        <!-- 使用按钮方式展示，以后改成布尔参数比较合适 -->
        <el-table-column label="状态" align="center">
          <template slot-scope="scope">
            <el-switch v-model="scope.row.status" :active-value="1" :inactive-value="2" @change="userStateChanged(scope.row)" />
          </template>
        </el-table-column>
        <!-- <el-table-column show-overflow-tooltip sortable prop="status" label="状态" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success':'danger'" disable-transitions>{{ scope.row.status === 1 ? '正常':'禁用' }}</el-tag>
          </template>
        </el-table-column> -->
        <el-table-column show-overflow-tooltip sortable prop="mail" label="邮箱" min-width="150" />
        <el-table-column show-overflow-tooltip sortable prop="mobile" label="手机号" min-width="110" />
        <el-table-column show-overflow-tooltip sortable prop="jobNumber" label="工号" min-width="80" />
        <el-table-column show-overflow-tooltip sortable prop="departments" label="部门" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="position" label="职位" min-width="80" />
        <el-table-column show-overflow-tooltip sortable prop="uidNumber" label="UID" align="center" min-width="70" />
        <el-table-column show-overflow-tooltip sortable prop="gidNumber" label="GID" align="center" min-width="70" />
        <el-table-column show-overflow-tooltip sortable prop="homeDirectory" label="主目录" min-width="140" />
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" min-width="90" />
        <el-table-column show-overflow-tooltip sortable prop="introduction" label="说明" min-width="100" />
        <el-table-column show-overflow-tooltip sortable prop="userDn" label="DN" min-width="200" />
        <el-table-column show-overflow-tooltip sortable prop="CreatedAt" label="创建时间" min-width="150" />
        <el-table-column show-overflow-tooltip sortable prop="UpdatedAt" label="更新时间" min-width="150" />
        <el-table-column fixed="right" label="操作" align="center" width="190">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="重置密码" effect="dark" placement="top">
              <el-popconfirm title="确定重置该用户密码吗？" @onConfirm="resetUserPassword(scope.row.username)">
                <el-button slot="reference" size="mini" icon="el-icon-key" circle type="warning" />
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip v-if="scope.row.syncState == 2" class="delete-popover" content="同步" effect="dark" placement="top">
              <el-popconfirm title="确定同步吗？" @onConfirm="singleSync(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-upload2" circle type="success" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[1, 5, 10, 30]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
          </div>
        </el-tab-pane>

        <!-- 待审核用户标签页 -->
        <el-tab-pane label="待审核用户" name="pending">
          <div class="tab-content">
            <el-form size="mini" :inline="true" :model="pendingParams" class="demo-form-inline">
              <el-form-item label="用户名">
                <el-input v-model.trim="pendingParams.username" style="width: 100px;" clearable placeholder="用户名" @keyup.enter.native="searchPending" @clear="searchPending" />
              </el-form-item>
              <el-form-item label="真实姓名">
                <el-input v-model.trim="pendingParams.nickname" style="width: 100px;" clearable placeholder="真实姓名" @keyup.enter.native="searchPending" @clear="searchPending" />
              </el-form-item>
              <el-form-item label="邮箱">
                <el-input v-model.trim="pendingParams.mail" style="width: 150px;" clearable placeholder="邮箱" @keyup.enter.native="searchPending" @clear="searchPending" />
              </el-form-item>
              <el-form-item label="审核状态">
                <el-select v-model.trim="pendingParams.status" style="width: 120px;" clearable placeholder="审核状态" @change="searchPending" @clear="searchPending">
                  <el-option label="待审核" :value="0" />
                  <el-option label="已通过" :value="1" />
                  <el-option label="已拒绝" :value="2" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button :loading="pendingLoading" icon="el-icon-search" type="primary" @click="searchPending">查询</el-button>
              </el-form-item>
              <el-form-item>
                <el-button :disabled="pendingMultipleSelection.length === 0" :loading="pendingLoading" icon="el-icon-delete" type="danger" @click="batchDeletePending">批量删除</el-button>
              </el-form-item>
            </el-form>

            <el-table v-loading="pendingLoading" :data="pendingTableData" border stripe style="width: 100%" @selection-change="handlePendingSelectionChange">
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
                    <el-button size="mini" icon="el-icon-check" circle type="success" @click="reviewPending(scope.row)" />
                  </el-tooltip>
                  <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
                    <el-popconfirm title="确定删除吗？" @onConfirm="singleDeletePending(scope.row.ID)">
                      <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
                    </el-popconfirm>
                  </el-tooltip>
                </template>
              </el-table-column>
            </el-table>

            <el-pagination
              :current-page="pendingParams.pageNum"
              :page-size="pendingParams.pageSize"
              :total="pendingTotal"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, prev, pager, next, sizes"
              background
              style="margin-top: 10px;float:right;margin-bottom: 10px;"
              @size-change="handlePendingSizeChange"
              @current-change="handlePendingCurrentChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="50%">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="80px" autocomplete="off">
          <el-row>
            <el-col :span="12">
              <el-form-item label="用户名" prop="username">
                <el-input ref="password" v-model.trim="dialogFormData.username" :disabled="disabled" placeholder="用户名（拼音）" autocomplete="off" name="user-username" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="真实姓名" prop="nickname">
                <el-input v-model.trim="dialogFormData.nickname" placeholder="真实姓名" autocomplete="off" name="user-nickname" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="花名" prop="givenName">
                <el-input v-model.trim="dialogFormData.givenName" placeholder="花名" autocomplete="off" name="user-givenname" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="邮箱" prop="mail">
                <el-input v-model.trim="dialogFormData.mail" placeholder="邮箱" autocomplete="off" name="user-email" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <!-- 修改用户时，不显示密码字段 -->
              <el-form-item v-if="dialogType === 'create'" label="密码" prop="password">
                <el-input v-model.trim="dialogFormData.password" autocomplete="new-password" :type="passwordType" placeholder="请输入密码" name="user-password" />
                <span class="show-pwd" @click="showPwd">
                  <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
                </span>
              </el-form-item>
            </el-col>
            <el-col v-if="dialogType === 'create'" :span="12">
              <el-form-item label="确认密码" prop="confirmPassword">
                <el-input v-model.trim="dialogFormData.confirmPassword" autocomplete="new-password" :type="passwordType" placeholder="请再次输入密码" name="user-confirm-password" />
                <span class="show-pwd" @click="showPwd">
                  <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
                </span>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="角色" prop="roleIds">
                <el-select v-model.trim="dialogFormData.roleIds" multiple placeholder="请选择角色" style="width:100%">
                  <el-option
                    v-for="item in roles"
                    :key="item.ID"
                    :label="item.name"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="状态" prop="status">
                <el-select v-model.trim="dialogFormData.status" placeholder="请选择状态" style="width:100%">
                  <el-option label="正常" :value="1" />
                  <el-option label="禁用" :value="2" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="手机号" prop="mobile">
                <el-input v-model.trim="dialogFormData.mobile" placeholder="手机号" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="工号" prop="jobNumber">
                <el-input v-model.trim="dialogFormData.jobNumber" placeholder="工号" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="职位" prop="position">
                <el-input v-model.trim="dialogFormData.position" placeholder="职业" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="所属部门" prop="departmentId">
                <treeselect
                  v-model="dialogFormData.departmentId"
                  :options="departmentsOptions"
                  placeholder="请选择部门"
                  :normalizer="normalizer"
                  value-consists-of="ALL"
                  :multiple="true"
                  :flat="true"
                  no-children-text="没有更多选项"
                  no-results-text="没有匹配的选项"
                  @input="treeselectInput"
                />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="地址" prop="postalAddress">
                <el-input v-model.trim="dialogFormData.postalAddress" type="textarea" placeholder="地址" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="说明" prop="introduction">
                <el-input v-model.trim="dialogFormData.introduction" type="textarea" placeholder="说明" :autosize="{minRows: 3, maxRows: 6}" show-word-limit maxlength="100" />
              </el-form-item>
            </el-col>
            <!-- Unix用户属性 -->
            <el-col :span="24">
              <el-divider content-position="left">Unix用户属性</el-divider>
            </el-col>
            <el-col :span="12">
              <el-form-item label="UID号" prop="uidNumber">
                <el-input-number v-model="dialogFormData.uidNumber" :min="1000" :max="65534" placeholder="Unix用户ID（留空自动分配）" style="width:100%" />
                <div style="font-size: 12px; color: #909399; margin-top: 4px;">留空将自动分配，建议从1000开始</div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="GID号" prop="gidNumber">
                <el-input-number v-model="dialogFormData.gidNumber" :min="0" :max="65534" placeholder="Unix组ID（留空默认使用UID号）" style="width:100%" />
                <div style="font-size: 12px; color: #909399; margin-top: 4px;">留空将默认使用UID号</div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="主目录" prop="homeDirectory">
                <el-input v-model.trim="dialogFormData.homeDirectory" placeholder="用户主目录路径" />
                <div style="font-size: 12px; color: #909399; margin-top: 4px;">默认：/home/用户名</div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="登录Shell" prop="loginShell">
                <el-select v-model.trim="dialogFormData.loginShell" placeholder="登录Shell" style="width:100%">
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
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>

      <!-- 重置密码结果对话框 -->
      <el-dialog
        title="密码重置成功"
        :visible.sync="resetPasswordDialogVisible"
        width="400px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        @close="closeResetPasswordDialog"
      >
        <div style="text-align: center;">
          <el-alert
            title="请保存新密码"
            type="warning"
            :closable="false"
            show-icon
            style="margin-bottom: 20px;"
          />
          <p style="margin-bottom: 10px; font-weight: bold;">用户：{{ resetUsername }}</p>
          <p style="margin-bottom: 20px; color: #606266;">新密码：</p>
          <el-input
            v-model="newPassword"
            readonly
            style="margin-bottom: 20px;"
          >
            <el-button
              slot="append"
              icon="el-icon-document-copy"
              @click="copyPassword"
            >
              复制
            </el-button>
          </el-input>
          <el-alert
            title="请立即保存密码，关闭对话框后将无法再次查看"
            type="info"
            :closable="false"
            show-icon
          />
        </div>
        <div slot="footer" class="dialog-footer">
          <el-button type="primary" @click="closeResetPasswordDialog">我已保存</el-button>
        </div>
      </el-dialog>

      <!-- 待审核用户审核对话框 -->
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
                    :options="reviewDepartmentsOptions"
                    placeholder="请选择分组（必选）"
                    :normalizer="reviewNormalizer"
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

            <el-divider content-position="left">用户信息（来自用户配置）</el-divider>
            <el-row>
              <el-col :span="12">
                <el-form-item label="手机号" prop="mobile">
                  <el-input v-model.trim="reviewFormData.mobile" placeholder="手机号" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="工号" prop="jobNumber">
                  <el-input v-model.trim="reviewFormData.jobNumber" placeholder="工号" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="职位" prop="position">
                  <el-input v-model.trim="reviewFormData.position" placeholder="职位" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="地址" prop="postalAddress">
                  <el-input v-model.trim="reviewFormData.postalAddress" placeholder="地址" />
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
                  <el-input-number v-model="reviewFormData.gidNumber" :min="0" :max="65534" placeholder="Unix组ID（留空默认使用UID号）" style="width:100%" />
                  <div style="font-size: 12px; color: #909399; margin-top: 4px;">留空将默认使用UID号</div>
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
import JSEncrypt from 'jsencrypt'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'
import { getUsers, createUser, updateUserById, batchDeleteUserByIds, changeUserStatus, syncOpenLdapUsersApi, syncSqlUsers } from '@/api/personnel/user'
import { resetPassword, getPendingUsers, reviewPendingUser, deletePendingUsers } from '@/api/system/user'
import { preConfigGetByUsername } from '@/api/personnel/userPreConfig'
import { getRoles } from '@/api/system/role'
import { getGroupTree } from '@/api/personnel/group'
import { Message } from 'element-ui'

export default {
  name: 'User',
  components: {
    Treeselect
  },
  props: {
    disabled: { // username 默认不可编辑，若需要至为可编辑，请（在新增和编辑处）去掉这个值的控制，且配合后端的ldap-user-name-modify配置使用
      type: Boolean,
      default: false
    }
  },
  data() {
    var checkPhone = (rule, value, callback) => {
      if (value) {
        const reg = /1\d{10}/
        if (reg.test(value)) {
          callback()
        } else {
          return callback(new Error('请输入正确的手机号'))
        }
      }
    }
    return {
      // 当前激活的标签页
      activeTab: 'users',
      // 查询参数
      params: {
        username: '',
        nickname: '',
        status: '',
        syncState: '',
        mobile: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,
      isUpdate: false,
      // 待审核用户相关
      pendingParams: {
        username: '',
        nickname: '',
        mail: '',
        status: 0, // 默认只显示待审核的
        pageNum: 1,
        pageSize: 10
      },
      pendingTableData: [],
      pendingTotal: 0,
      pendingLoading: false,
      pendingMultipleSelection: [],
      reviewDialogVisible: false,
      reviewLoading: false,
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
        uidNumber: undefined,
        gidNumber: undefined,
        homeDirectory: '',
        loginShell: '/bin/bash',
        mobile: '',
        jobNumber: '',
        position: '',
        postalAddress: '',
        departments: '',
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
      reviewDepartmentsOptions: [],
      // 部门信息数据
      treeselectValue: 0,
      // 角色
      roles: [],
      // 部门信息
      departmentsOptions: [],

      passwordType: 'password',

      publicKey: process.env.VUE_APP_PUBLIC_KEY,

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        username: '',
        password: '',
        confirmPassword: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: [],
        ID: '',
        mail: '',
        givenName: '',
        jobNumber: '',
        postalAddress: '',
        departments: '',
        position: '',
        departmentId: undefined,
        uidNumber: undefined,
        gidNumber: 0,
        homeDirectory: '',
        loginShell: '/bin/bash'
      },
      dialogFormRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        password: [
          { 
            validator: (rule, value, callback) => {
              if (this.dialogType === 'create') {
                if (!value) {
                  callback(new Error('请输入密码'))
                } else if (value.length < 6 || value.length > 30) {
                  callback(new Error('密码长度在 6 到 30 个字符'))
                } else {
                  callback()
                }
              } else {
                // 编辑时密码可选
                if (value && (value.length < 6 || value.length > 30)) {
                  callback(new Error('密码长度在 6 到 30 个字符'))
                } else {
                  callback()
                }
              }
            }, 
            trigger: 'blur' 
          }
        ],
        confirmPassword: [
          { 
            validator: (rule, value, callback) => {
              if (this.dialogType === 'create') {
                if (!value) {
                  callback(new Error('请再次输入密码'))
                } else if (value !== this.dialogFormData.password) {
                  callback(new Error('两次输入的密码不一致'))
                } else {
                  callback()
                }
              } else {
                // 编辑时不需要确认密码
                callback()
              }
            }, 
            trigger: 'blur' 
          }
        ],
        mail: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        nickname: [
          { required: true, message: '请输入真实姓名', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        jobNumber: [
          { required: false, message: '请输入工号', trigger: 'blur' },
          { min: 0, max: 20, message: '长度在 0 到 20 个字符', trigger: 'blur' }
        ],
        mobile: [
          { required: false, validator: checkPhone, trigger: 'blur' }
        ],
        status: [
          { required: false, message: '请选择状态', trigger: 'change' }
        ],
        departmentId: [
          { required: false, message: '请选择部门', trigger: 'change' }
        ],
        roleIds: [
          { required: false, message: '请选择角色', trigger: 'change' }
        ],
        introduction: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 255, message: '长度在 0 到 255 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      changeUserStatusFormData: {
        id: '',
        status: ''
      },

      // 重置密码结果对话框
      resetPasswordDialogVisible: false,
      newPassword: '',
      resetUsername: ''
    }
  },
  created() {
    this.getTableData()
    this.getRoles()
    this.loadReviewDepartments()
  },
  methods: {
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getUsers(this.params)
        data.users.forEach(item => {
          const dataStrArr = item.departmentId.split(',')
          const dataIntArr = []
          dataStrArr.forEach(item => {
            dataIntArr.push(+item)
          })
          item.departmentId = dataIntArr
        })
        this.tableData = data.users
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    // 获取所有的分组信息，用于弹框选取上级分组
    async getAllGroups() {
      this.loading = true
      try {
        const checkParams = {
          pageNum: 1,
          pageSize: 1000 // 平常百姓人家应该不会有这么多数据吧
        }
        const { data } = await getGroupTree(checkParams)
        this.departmentsOptions = [{ ID: 0, groupName: '请选择部门信息', groupType: 'T', children: data }]
      } finally {
        this.loading = false
      }
    },
    // 获取角色数据
    async getRoles() {
      const res = await getRoles(null)

      this.roles = res.data.roles
    },

      // 新增
    async create() {
      // 重置表单数据，确保没有残留值
      this.dialogFormData = {
        username: '',
        password: '',
        confirmPassword: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: [],
        ID: '',
        mail: '',
        givenName: '',
        jobNumber: '',
        postalAddress: '',
        departments: '',
        position: '',
        departmentId: undefined,
        uidNumber: undefined,
        gidNumber: undefined,
        homeDirectory: '',
        loginShell: '/bin/bash'
      }
      
      // 加载分组列表，并设置默认分组CMPLabHPC
      await this.getAllGroups()
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
      // 重置表单验证状态
      if (this.$refs['dialogForm']) {
        this.$refs['dialogForm'].resetFields()
      }
      
      // 查找CMPLabHPC分组并设置为默认值（在resetFields之后设置）
      const cmplabHpcId = findCMPLabHPC(this.departmentsOptions[0]?.children || [])
      if (cmplabHpcId) {
        this.dialogFormData.departmentId = [cmplabHpcId]
      }
      
      this.dialogFormTitle = '新增用户'
      this.dialogType = 'create'
      this.disabled = false
      this.passwordType = 'password'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.disabled = true
      this.getAllGroups()
      this.dialogFormData.ID = row.ID
      this.dialogFormData.username = row.username
      this.dialogFormData.password = ''
      this.dialogFormData.nickname = row.nickname
      this.dialogFormData.status = row.status
      this.dialogFormData.mobile = row.mobile
      this.dialogFormData.introduction = row.introduction
      // 遍历角色数组，获取角色ID
      this.dialogFormData.roleIds = row.roles.map(item => item.ID)

      this.dialogFormTitle = '修改用户'
      this.dialogType = 'update'
      this.passwordType = 'password'
      this.dialogFormVisible = true

      this.dialogFormData.mail = row.mail
      this.dialogFormData.givenName = row.givenName
      this.dialogFormData.jobNumber = row.jobNumber
      this.dialogFormData.postalAddress = row.postalAddress
      this.dialogFormData.departments = row.departments
      this.dialogFormData.departmentId = row.departmentId
      this.dialogFormData.position = row.position
      this.dialogFormData.uidNumber = row.uidNumber || undefined
      this.dialogFormData.gidNumber = row.gidNumber || undefined
      this.dialogFormData.homeDirectory = row.homeDirectory || ''
      this.dialogFormData.loginShell = row.loginShell || '/bin/bash'
    },

    // 将 部门id 转换为 部门name
    setDepartmentNameByDepartmentId() {
      const ids = this.dialogFormData.departmentId
      if (!ids || !ids.length) return
      const departments = []
      // 深度优先遍函数
      const dfs = (node, cb) => {
        if (!node) return
        cb(node)
        if (node.children && node.children.length) {
          node.children.forEach(item => {
            dfs(item, cb)
          })
        }
      }
      dfs(this.departmentsOptions[0], node => {
        if (ids.includes(node.ID)) {
          departments.push(node.groupName)
        }
      })
      this.dialogFormData.departments = departments.join(',')
    },

    // 判断结果
    judgeResult(res) {
      if (res.code === 0) {
        Message({
          showClose: true,
          message: '操作成功',
          type: 'success'
        })
      }
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          // 如果是创建用户，验证密码和确认密码
          if (this.dialogType === 'create') {
            if (!this.dialogFormData.password) {
              Message({
                showClose: true,
                message: '请输入密码',
                type: 'error'
              })
              return false
            }
            if (this.dialogFormData.password !== this.dialogFormData.confirmPassword) {
              Message({
                showClose: true,
                message: '密码和确认密码不一致',
                type: 'error'
              })
              return false
            }
          }
          
          this.submitLoading = true
          // 在这里自动填充下部门字段
          this.setDepartmentNameByDepartmentId()
          this.dialogFormDataCopy = { ...this.dialogFormData }
          // 移除确认密码字段，不发送到后端
          delete this.dialogFormDataCopy.confirmPassword
          
          if (this.dialogFormData.password !== '') {
            // 密码RSA加密处理
            const encryptor = new JSEncrypt()
            // 设置公钥
            encryptor.setPublicKey(this.publicKey)
            // 加密密码
            const encPassword = encryptor.encrypt(this.dialogFormData.password)
            this.dialogFormDataCopy.password = encPassword
            // 同时加密确认密码（如果存在）
            if (this.dialogFormData.confirmPassword) {
              this.dialogFormDataCopy.confirmPassword = encryptor.encrypt(this.dialogFormData.confirmPassword)
            }
          }
          try {
            if (this.dialogType === 'create') {
              await createUser(this.dialogFormDataCopy).then(res => {
                this.judgeResult(res)
              })
            } else {
              await updateUserById(this.dialogFormDataCopy).then(res => {
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
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        username: '',
        password: '',
        confirmPassword: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: [],
        ID: '',
        mail: '',
        givenName: '',
        jobNumber: '',
        postalAddress: '',
        departments: '',
        position: '',
        departmentId: undefined,
        uidNumber: undefined,
        gidNumber: undefined,
        homeDirectory: '',
        loginShell: '/bin/bash'
      }
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const userIds = []
        this.multipleSelection.forEach(x => {
          userIds.push(x.ID)
        })
        try {
          await batchDeleteUserByIds({ userIds: userIds }).then(res => {
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
      this.$confirm('此操作批量将数据库的用户同步到Ldap, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const userIds = []
        this.multipleSelection.forEach(x => {
          userIds.push(x.ID)
        })
        try {
          await syncSqlUsers({ userIds: userIds }).then(res => {
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

    // 监听 switch 开关 状态改变
    async userStateChanged(userInfo) {
      this.changeUserStatusFormData.id = userInfo.ID
      this.changeUserStatusFormData.status = userInfo.status
      const { code } = await changeUserStatus(this.changeUserStatusFormData)
      if (code !== 0) {
        return Message.error('更新用户状态失败')
      }
      Message.success('更新用户状态成功')
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      try {
        await batchDeleteUserByIds({ userIds: [Id] }).then(res => {
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
        await syncSqlUsers({ userIds: [Id] }).then(res => {
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
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
        label: node.groupType + '=' + node.groupName,
        // 只禁用root分组，允许选择ou类型的分组（如CMPLabHPC）
        isDisabled: node.groupName === 'root' || (node.groupType === 'T' && node.ID === 0),
        children: node.children
      }
    },
    treeselectInput(value) {
      this.treeselectValue = value
    },
    syncOpenLdapUsers() {
      this.loading = true
      syncOpenLdapUsersApi().then(res => {
        this.judgeResult(res)
        this.loading = false
        this.getTableData()
      })
    },
    // 同步所有MySQL用户到LDAP
    async syncAllUsersToLdap() {
      this.$confirm('此操作将同步所有MySQL用户到LDAP, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        try {
          // 传递空数组表示同步所有用户
          await syncSqlUsers({ userIds: [] }).then(res => {
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

    // 重置用户密码
    async resetUserPassword(username) {
      this.loading = true
      try {
        const res = await resetPassword({ username: username })
        if (res.code === 0) {
          this.newPassword = res.data.newPassword
          this.resetUsername = username
          this.resetPasswordDialogVisible = true
          Message({
            showClose: true,
            message: '密码重置成功',
            type: 'success'
          })
        } else {
          Message({
            showClose: true,
            message: res.msg || '密码重置失败',
            type: 'error'
          })
        }
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    // 复制密码到剪贴板
    copyPassword() {
      const textArea = document.createElement('textarea')
      textArea.value = this.newPassword
      document.body.appendChild(textArea)
      textArea.select()
      try {
        document.execCommand('copy')
        Message({
          showClose: true,
          message: '密码已复制到剪贴板',
          type: 'success'
        })
      } catch (err) {
        Message({
          showClose: true,
          message: '复制失败，请手动复制',
          type: 'error'
        })
      }
      document.body.removeChild(textArea)
    },

    // 关闭重置密码对话框
    closeResetPasswordDialog() {
      this.resetPasswordDialogVisible = false
      this.newPassword = ''
      this.resetUsername = ''
    },
    // 标签页切换
    handleTabClick(tab) {
      if (tab.name === 'pending' && this.pendingTableData.length === 0) {
        this.searchPending()
      }
    },
    // 待审核用户相关方法
    async searchPending() {
      this.pendingParams.pageNum = 1
      await this.getPendingTableData()
    },
    async getPendingTableData() {
      this.pendingLoading = true
      try {
        const res = await getPendingUsers(this.pendingParams)
        if (res.code === 0) {
          this.pendingTableData = res.data.pendingUsers || []
          this.pendingTotal = res.data.total || 0
        } else {
          Message({
            showClose: true,
            message: res.msg || '获取待审核用户列表失败',
            type: 'error'
          })
        }
      } catch (error) {
        Message({
          showClose: true,
          message: error.message || '获取待审核用户列表失败',
          type: 'error'
        })
      } finally {
        this.pendingLoading = false
      }
    },
    handlePendingSelectionChange(val) {
      this.pendingMultipleSelection = val
    },
    handlePendingSizeChange(val) {
      this.pendingParams.pageSize = val
      this.getPendingTableData()
    },
    handlePendingCurrentChange(val) {
      this.pendingParams.pageNum = val
      this.getPendingTableData()
    },
    reviewNormalizer(node) {
      if (node.children && !node.children.length) {
        delete node.children
      }
      return {
        id: node.ID,
        label: node.groupName,
        isDisabled: node.groupName === 'root' || (node.groupType === 'T' && node.ID === 0),
        children: node.children
      }
    },
    async loadReviewDepartments() {
      try {
        const res = await getGroupTree({ pageNum: 1, pageSize: 1000 })
        if (res.code === 0) {
          this.reviewDepartmentsOptions = res.data || []
        }
      } catch (error) {
        Message({
          showClose: true,
          message: '加载分组列表失败',
          type: 'error'
        })
      }
    },
    async reviewPending(row) {
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
        uidNumber: undefined,
        gidNumber: undefined,
        homeDirectory: '',
        loginShell: '/bin/bash',
        mobile: '',
        jobNumber: '',
        position: '',
        postalAddress: '',
        departments: '',
        allowSudo: false,
        allowSSHKey: true,
        sudoRules: ''
      }
      
      // 获取用户预配置信息
      try {
        const preConfigRes = await preConfigGetByUsername({ username: row.username })
        if (preConfigRes.code === 0 && preConfigRes.data) {
          const preConfig = preConfigRes.data
          // 预填充用户配置信息
          if (preConfig.uidNumber > 0) {
            this.reviewFormData.uidNumber = preConfig.uidNumber
          }
          if (preConfig.gidNumber > 0) {
            this.reviewFormData.gidNumber = preConfig.gidNumber
          }
          if (preConfig.homeDirectory) {
            this.reviewFormData.homeDirectory = preConfig.homeDirectory
          }
          if (preConfig.loginShell) {
            this.reviewFormData.loginShell = preConfig.loginShell
          }
          if (preConfig.mobile) {
            this.reviewFormData.mobile = preConfig.mobile
          }
          if (preConfig.jobNumber) {
            this.reviewFormData.jobNumber = preConfig.jobNumber
          }
          if (preConfig.position) {
            this.reviewFormData.position = preConfig.position
          }
          if (preConfig.postalAddress) {
            this.reviewFormData.postalAddress = preConfig.postalAddress
          }
          if (preConfig.departments) {
            this.reviewFormData.departments = preConfig.departments
          }
          // 如果预配置有分组信息
          if (preConfig.departmentId) {
            const deptIds = preConfig.departmentId.split(',').filter(id => id).map(id => parseInt(id))
            if (deptIds.length > 0) {
              this.reviewFormData.departmentId = deptIds
            }
          }
        }
      } catch (error) {
        console.log('获取用户预配置失败:', error)
      }
      
      // 如果没有预配置分组，查找CMPLabHPC分组并设置为默认值
      if (this.reviewFormData.departmentId.length === 0) {
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
        const cmplabHpcId = findCMPLabHPC(this.reviewDepartmentsOptions || [])
        if (cmplabHpcId) {
          this.reviewFormData.departmentId = [cmplabHpcId]
        }
      }
      
      this.reviewDialogVisible = true
    },
    cancelReview() {
      this.reviewDialogVisible = false
      if (this.$refs.reviewForm) {
        this.$refs.reviewForm.resetFields()
      }
    },
    async submitReview() {
      this.$refs.reviewForm.validate(async (valid) => {
        if (valid) {
          this.reviewLoading = true
          try {
            const res = await reviewPendingUser(this.reviewFormData)
            if (res.code === 0) {
              Message({
                showClose: true,
                message: '审核成功',
                type: 'success'
              })
              this.reviewDialogVisible = false
              this.getPendingTableData()
              // 如果审核通过，刷新用户列表
              if (this.reviewFormData.status === 1) {
                this.getTableData()
              }
            } else {
              Message({
                showClose: true,
                message: res.msg || '审核失败',
                type: 'error'
              })
            }
          } catch (error) {
            Message({
              showClose: true,
              message: error.message || '审核失败',
              type: 'error'
            })
          } finally {
            this.reviewLoading = false
          }
        }
      })
    },
    async singleDeletePending(id) {
      try {
        const res = await deletePendingUsers({ ids: [id] })
        if (res.code === 0) {
          Message({
            showClose: true,
            message: '删除成功',
            type: 'success'
          })
          this.getPendingTableData()
        } else {
          Message({
            showClose: true,
            message: res.msg || '删除失败',
            type: 'error'
          })
        }
      } catch (error) {
        Message({
          showClose: true,
          message: error.message || '删除失败',
          type: 'error'
        })
      }
    },
    async batchDeletePending() {
      if (this.pendingMultipleSelection.length === 0) {
        Message({
          showClose: true,
          message: '请选择要删除的记录',
          type: 'warning'
        })
        return
      }
      try {
        const ids = this.pendingMultipleSelection.map(item => item.ID)
        const res = await deletePendingUsers({ ids })
        if (res.code === 0) {
          Message({
            showClose: true,
            message: '批量删除成功',
            type: 'success'
          })
          this.getPendingTableData()
        } else {
          Message({
            showClose: true,
            message: res.msg || '批量删除失败',
            type: 'error'
          })
        }
      } catch (error) {
        Message({
          showClose: true,
          message: error.message || '批量删除失败',
          type: 'error'
        })
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

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 3px;
    font-size: 16px;
    color: #889aa4;
    cursor: pointer;
    user-select: none;
  }
</style>
