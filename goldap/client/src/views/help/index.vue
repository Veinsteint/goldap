<template>
  <div class="help-container">
    <div class="help-content">
      <!-- 页面标题 -->
      <div class="help-header">
        <h1 class="help-title">
          <svg-icon icon-class="question" style="margin-right: 12px;" />
          LDAP统一登录管理系统 - 使用帮助
        </h1>
        <p class="help-subtitle">实验室服务器统一身份认证与管理平台</p>
      </div>

      <!-- 目录导航 -->
      <div class="help-nav">
        <el-anchor :offset="80" :bounds="20">
          <el-anchor-link href="#ldap-intro" title="1. LDAP简介" />
          <el-anchor-link href="#unified-login" title="2. 统一登录管理" />
          <el-anchor-link href="#shared-storage" title="3. 共享存储空间" />
          <el-anchor-link href="#user-registration" title="4. 用户注册流程" />
          <el-anchor-link href="#admin-guide" title="5. 管理员指南" />
        </el-anchor>
      </div>

      <!-- 内容区域 -->
      <div class="help-body">
        <!-- 1. LDAP简介 -->
        <section id="ldap-intro" class="help-section">
          <h2 class="section-title">
            <el-icon class="section-icon"><InfoFilled /></el-icon>
            1. LDAP的作用和功能
          </h2>
          <div class="section-content">
            <el-card shadow="hover" class="info-card">
              <h3>什么是LDAP？</h3>
              <p>
                <strong>LDAP（Lightweight Directory Access Protocol，轻量级目录访问协议）</strong>
                是一种开放的网络协议标准，用于访问和管理分布式目录信息服务。
              </p>
              
              <h3 style="margin-top: 24px;">LDAP的核心功能</h3>
              <el-row :gutter="20" style="margin-top: 16px;">
                <el-col :span="8">
                  <el-card shadow="never" class="feature-card">
                    <div class="feature-icon">
                      <svg-icon icon-class="user" style="font-size: 48px; color: #409EFF;" />
                    </div>
                    <h4>统一身份认证</h4>
                    <p>为多台服务器提供统一的用户账号和密码管理，实现统一认证</p>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card shadow="never" class="feature-card">
                    <div class="feature-icon">
                      <svg-icon icon-class="lock" style="font-size: 48px; color: #67C23A;" />
                    </div>
                    <h4>集中权限管理</h4>
                    <p>集中管理用户权限、组权限以及sudo规则和docker权限，确保安全访问控制</p>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card shadow="never" class="feature-card">
                    <div class="feature-icon">
                      <svg-icon icon-class="key" style="font-size: 48px; color: #E6A23C;" />
                    </div>
                    <h4>SSH密钥管理</h4>
                    <p>统一管理SSH公钥，支持多设备免密登录，提高使用便利性</p>
                  </el-card>
                </el-col>
              </el-row>

              <h3 style="margin-top: 32px;">LDAP的优势</h3>
              <el-timeline style="margin-top: 16px;">
                <el-timeline-item timestamp="集中管理" placement="top">
                  <el-card>
                    <p>所有用户信息集中存储，管理员只需在一个地方进行维护，大大减少管理成本</p>
                  </el-card>
                </el-timeline-item>
                <el-timeline-item timestamp="统一认证" placement="top">
                  <el-card>
                    <p>用户只需记住一套账号密码，即可登录所有配置了LDAP的服务器</p>
                  </el-card>
                </el-timeline-item>
                <el-timeline-item timestamp="安全可靠" placement="top">
                  <el-card>
                    <p>支持TLS加密传输，密码加密存储，提供企业级安全保障</p>
                  </el-card>
                </el-timeline-item>
                <el-timeline-item timestamp="易于扩展" placement="top">
                  <el-card>
                    <p>可以轻松添加新的服务器或服务，只需配置LDAP客户端即可</p>
                  </el-card>
                </el-timeline-item>
              </el-timeline>
            </el-card>
          </div>
        </section>

        <!-- 2. 统一登录管理 -->
        <section id="unified-login" class="help-section">
          <h2 class="section-title">
            <el-icon class="section-icon"><Connection /></el-icon>
            2. 实验室服务器统一登录管理
          </h2>
          <div class="section-content">
            <el-card shadow="hover" class="info-card">
              <h3>服务器架构</h3>
              <p>我们的实验室拥有三台高性能计算服务器，通过LDAP实现统一身份认证：</p>
              
              <div class="server-architecture">
                <div class="architecture-diagram">
                  <div class="ldap-server">
                    <div class="server-box ldap-box">
                      <svg-icon icon-class="server" style="font-size: 32px; margin-bottom: 8px;" />
                      <h4>LDAP服务器</h4>
                      <p>统一身份认证中心</p>
                    </div>
                  </div>
                  
                  <div class="connection-lines">
                    <div class="line"></div>
                    <div class="line"></div>
                    <div class="line"></div>
                  </div>
                  
                  <div class="compute-servers">
                    <div class="server-box hpc-box">
                      <svg-icon icon-class="monitor" style="font-size: 28px; margin-bottom: 8px;" />
                      <h4>CMPLab-HPC1 (192.168.11.6)</h4>
                      <p>高性能计算服务器1</p>
                    </div>
                    <div class="server-box hpc-box">
                      <svg-icon icon-class="monitor" style="font-size: 28px; margin-bottom: 8px;" />
                      <h4>CMPLab-HPC2 (192.168.11.7)</h4>
                      <p>高性能计算服务器2</p>
                    </div>
                    <div class="server-box hpc-box">
                      <svg-icon icon-class="monitor" style="font-size: 28px; margin-bottom: 8px;" />
                      <h4>CMPLab-HPC3 (192.168.11.8)</h4>
                      <p>高性能计算服务器3</p>
                    </div>
                  </div>
                </div>
              </div>

              <el-alert
                title="统一登录说明"
                type="success"
                :closable="false"
                show-icon
                style="margin-top: 24px;"
              >
                <template slot="default">
                  <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                    <li>所有三台服务器（CMPLab-HPC1、CMPLab-HPC2、CMPLab-HPC3）使用相同的LDAP认证系统</li>
                    <li>用户只需注册一次账号，即可在三台服务器上使用相同的用户名和密码登录</li>
                    <li>管理员可以在Web界面统一管理所有用户，无需在每台服务器上单独配置</li>
                    <li>用户权限、(sudo、docker)组权限和sudo规则在所有服务器上保持一致</li>
                  </ul>
                </template>
              </el-alert>
            </el-card>
          </div>
        </section>

        <!-- 3. 共享存储空间 -->
        <section id="shared-storage" class="help-section">
          <h2 class="section-title">
            <el-icon class="section-icon"><FolderOpened /></el-icon>
            3. 共享存储空间
          </h2>
          <div class="section-content">
            <el-card shadow="hover" class="info-card">
              <h3>统一Home目录</h3>
              <p>三台服务器的用户home目录位于同一个共享存储空间上，这意味着：</p>
              
              <el-row :gutter="20" style="margin-top: 20px;">
                <el-col :span="12">
                  <el-card shadow="never" class="benefit-card">
                    <div class="benefit-icon">
                      <el-icon style="font-size: 40px; color: #409EFF;"><Folder /></el-icon>
                    </div>
                    <h4>数据同步</h4>
                    <p>在任意一台服务器上创建或修改的文件，在其他服务器上立即可见</p>
                  </el-card>
                </el-col>
                <el-col :span="12">
                  <el-card shadow="never" class="benefit-card">
                    <div class="benefit-icon">
                      <el-icon style="font-size: 40px; color: #67C23A;"><DocumentCopy /></el-icon>
                    </div>
                    <h4>无缝切换</h4>
                    <p>可以在不同服务器间自由切换，无需担心数据迁移或同步问题</p>
                  </el-card>
                </el-col>
              </el-row>

              <div class="storage-diagram" style="margin-top: 30px;">
                <div class="storage-visual">
                  <div class="storage-center">
                    <el-icon style="font-size: 48px; color: #409EFF;"><FolderOpened /></el-icon>
                    <h4>共享存储</h4>
                    <p>/home/username</p>
                  </div>
                  <div class="storage-connections">
                    <div class="storage-link">
                      <span>CMPLab-HPC1</span>
                    </div>
                    <div class="storage-link">
                      <span>CMPLab-HPC2</span>
                    </div>
                    <div class="storage-link">
                      <span>CMPLab-HPC3</span>
                    </div>
                  </div>
                </div>
              </div>

              <el-alert
                title="注意事项"
                type="warning"
                :closable="false"
                show-icon
                style="margin-top: 24px;"
              >
                <template slot="default">
                  <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                    <li>所有服务器上的用户home目录路径相同，例如：<code>/home/wwtian</code></li>
                    <li>请确保在多台服务器上同时操作同一文件时注意文件锁定，避免数据冲突</li>
                    <li>建议使用版本控制系统（如Git）管理代码，避免并发编辑冲突</li>
                    <li>大文件传输时请注意网络带宽，避免影响其他用户</li>
                  </ul>
                </template>
              </el-alert>
            </el-card>
          </div>
        </section>

        <!-- 4. 用户注册流程 -->
        <section id="user-registration" class="help-section">
          <h2 class="section-title">
            <el-icon class="section-icon"><UserFilled /></el-icon>
            4. 普通用户注册流程详解
          </h2>
          <div class="section-content">
            <el-card shadow="hover" class="info-card">
              <h3>注册流程概览</h3>
              <el-steps :active="4" finish-status="success" align-center style="margin: 30px 0;">
                <el-step title="填写注册信息" description="填写用户名、邮箱等信息" />
                <el-step title="提交注册申请" description="等待管理员审核" />
                <el-step title="管理员审核" description="管理员审核通过" />
                <el-step title="登录使用" description="使用账号密码登录" />
              </el-steps>

              <el-divider content-position="left">详细步骤说明（以用户 wwtian 为例）</el-divider>

              <!-- 步骤1: 访问注册页面 -->
              <div class="step-detail">
                <h4>
                  <el-tag type="primary" size="small">步骤 1</el-tag>
                  访问注册页面
                </h4>
                <p>打开浏览器，访问自助密码服务页面，点击"账号注册"标签页。</p>
                <el-card shadow="never" style="margin: 16px 0; background: #f5f7fa; border-left: 4px solid #409EFF;">
                  <div style="padding: 16px;">
                    <p style="margin: 0; color: #606266;">
                      <el-icon style="color: #409EFF; margin-right: 8px;"><InfoFilled /></el-icon>
                      在浏览器中访问 <code>/password-service</code> 页面，选择"账号注册"标签页即可开始注册流程。
                    </p>
                  </div>
                </el-card>
              </div>

              <!-- 步骤2: 查看有效用户名 -->
              <div class="step-detail">
                <h4>
                  <el-tag type="warning" size="small">步骤 2</el-tag>
                  查看有效用户名列表（预配置模式）
                </h4>
                <p>
                  如果系统开启了<strong>预配置模式</strong>，只有预配置列表中的用户名才能注册。
                  点击"查看有效用户名"按钮查看可用用户名列表。
                </p>
                
                <el-card shadow="never" style="margin: 16px 0; background: #f5f7fa;">
                  <div class="form-example">
                    <h5 style="margin-bottom: 16px; color: #606266;">有效用户名列表示例：</h5>
                    <el-table :data="validUsernamesExample" border style="width: 100%">
                      <el-table-column prop="username" label="用户名" width="150" />
                      <el-table-column prop="nickname" label="姓名" width="120" />
                      <el-table-column prop="status" label="状态" width="120">
                        <template slot-scope="scope">
                          <el-tag
                            :type="scope.row.status === '未注册' ? 'info' : 'success'"
                            size="small"
                          >
                            {{ scope.row.status }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="remark" label="备注" />
                    </el-table>
                    <p style="margin-top: 12px; font-size: 12px; color: #909399;">
                      <el-icon><InfoFilled /></el-icon>
                      只有状态为"未注册"的用户名才能进行注册
                    </p>
                  </div>
                </el-card>
              </div>

              <!-- 步骤3: 填写注册表单 -->
              <div class="step-detail">
                <h4>
                  <el-tag type="success" size="small">步骤 3</el-tag>
                  填写注册表单
                </h4>
                <p>按照以下示例填写注册信息（以用户 <strong>wwtian</strong> 为例）：</p>

                <el-card shadow="never" style="margin: 16px 0; background: #f5f7fa;">
                  <div class="form-example">
                    <h5 style="margin-bottom: 16px; color: #606266;">注册表单示例：</h5>
                    
                    <!-- 模拟注册表单 -->
                    <el-form
                      :model="exampleForm"
                      label-width="100px"
                      size="medium"
                      class="example-register-form"
                      disabled
                    >
                      <el-form-item label="用户名">
                        <el-input
                          v-model="exampleForm.username"
                          placeholder="请输入用户名（3-50个字符）"
                          prefix-icon="el-icon-user"
                        >
                          <template slot="append">
                            <el-button icon="el-icon-view" size="small">查看有效用户名</el-button>
                          </template>
                        </el-input>
                        <div class="form-tip">
                          <el-tag type="warning" size="mini">注意</el-tag>
                          只有预配置列表中的用户名才能注册
                        </div>
                      </el-form-item>

                      <el-form-item label="邮箱">
                        <el-input
                          v-model="exampleForm.email"
                          placeholder="请输入邮箱地址"
                          prefix-icon="el-icon-message"
                        />
                        <div class="form-tip">用于接收验证码和重要通知</div>
                      </el-form-item>

                      <el-form-item label="真实姓名">
                        <el-input
                          v-model="exampleForm.realName"
                          placeholder="请输入真实姓名（仅支持英文字母）"
                          prefix-icon="el-icon-user-solid"
                        />
                        <div class="form-tip">只能包含英文字母（a-z, A-Z）</div>
                      </el-form-item>

                      <el-form-item label="密码">
                        <el-input
                          v-model="exampleForm.password"
                          type="password"
                          placeholder="请输入密码（至少6位）"
                          prefix-icon="el-icon-lock"
                          show-password
                        />
                        <div class="password-strength" style="margin-top: 10px;">
                          <div class="strength-label">密码强度：</div>
                          <div class="strength-bar">
                            <div class="strength-item strong" style="width: 100%;"></div>
                          </div>
                          <span class="strength-text strong">强</span>
                        </div>
                      </el-form-item>

                      <el-form-item label="确认密码">
                        <el-input
                          v-model="exampleForm.confirmPassword"
                          type="password"
                          placeholder="请再次输入密码"
                          prefix-icon="el-icon-lock"
                          show-password
                        />
                      </el-form-item>

                      <el-form-item label="备注">
                        <el-input
                          v-model="exampleForm.remark"
                          type="textarea"
                          :rows="3"
                          placeholder="选填：注册原因或备注信息"
                          maxlength="200"
                          show-word-limit
                        />
                      </el-form-item>

                      <el-form-item>
                        <el-button type="primary">注册</el-button>
                        <el-button>重置</el-button>
                      </el-form-item>
                    </el-form>
                  </div>
                </el-card>

                <el-alert
                  title="填写注意事项"
                  type="info"
                  :closable="false"
                  show-icon
                  style="margin-top: 16px;"
                >
                  <template slot="default">
                    <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                      <li><strong>用户名</strong>：3-50个字符，只能包含字母、数字和下划线，例如：<code>wwtian</code></li>
                      <li><strong>邮箱</strong>：必须填写有效的邮箱地址，用于接收审核通知，例如：<code>tianweiwei@cimrbj.ac.cn</code></li>
                      <li><strong>真实姓名</strong>：只能包含英文字母，例如：<code>weiweitian</code> 或 <code>WeiweiTian</code></li>
                      <li><strong>密码</strong>：至少6位，建议使用包含大小写字母、数字和特殊字符的强密码</li>
                      <li><strong>备注</strong>：可选，可以填写注册原因或需要说明的信息</li>
                    </ul>
                  </template>
                </el-alert>
              </div>

              <!-- 步骤4: 提交注册 -->
              <div class="step-detail">
                <h4>
                  <el-tag type="info" size="small">步骤 4</el-tag>
                  提交注册申请
                </h4>
                <p>填写完所有信息后，点击"注册"按钮提交申请。系统会显示注册成功提示：</p>
                
                <el-alert
                  title="注册成功"
                  type="success"
                  :closable="false"
                  show-icon
                  style="margin: 16px 0;"
                >
                  <template slot="default">
                    <p>您的注册申请已提交，等待管理员审核后即可使用。</p>
                    <p style="margin-top: 10px;">审核结果将通过邮箱通知您。</p>
                  </template>
                </el-alert>
              </div>

              <!-- 步骤5: 等待审核 -->
              <div class="step-detail" style="margin-bottom: 60px;">
                <h4>
                  <el-tag type="warning" size="small">步骤 5</el-tag>
                  等待管理员审核
                </h4>
                <p>提交注册申请后，需要等待管理员审核。审核通过后，您将收到邮件通知，即可使用账号登录系统。</p>
              </div>
              
              <div style="margin-top: 200px;"><p/><p/></div>
              <!-- 忘记密码 -->
              <el-divider content-position="left" style="margin-top: 40px;">忘记密码</el-divider>
              <div class="step-detail" style="margin-top: 0;">
                <h4>如何重置密码？</h4>
                <p>如果您忘记了密码，可以通过以下步骤重置：</p>
                
                <el-steps :active="2" finish-status="success" align-center style="margin: 20px 0;">
                  <el-step title="验证身份" description="输入注册邮箱" />
                  <el-step title="获取验证码" description="系统发送验证码到邮箱" />
                  <el-step title="重置密码" description="输入验证码，系统生成新密码" />
                  <el-step title="完成" description="使用新密码登录" />
                </el-steps>

                <el-card shadow="never" style="margin: 16px 0; background: #f5f7fa;">
                  <div class="form-example">
                    <h5 style="margin-bottom: 16px; color: #606266;">忘记密码表单示例：</h5>
                    <el-form
                      :model="forgotPasswordExample"
                      label-width="120px"
                      size="medium"
                      class="example-register-form"
                      disabled
                    >
                      <el-form-item label="邮箱地址">
                        <el-input
                          v-model="forgotPasswordExample.mail"
                          placeholder="请输入注册时使用的邮箱地址"
                          prefix-icon="el-icon-message"
                        >
                          <template slot="append">
                            <el-button type="primary" size="small">发送验证码</el-button>
                          </template>
                        </el-input>
                      </el-form-item>
                      <el-form-item label="验证码">
                        <el-input
                          v-model="forgotPasswordExample.code"
                          placeholder="请输入邮箱收到的验证码"
                          prefix-icon="el-icon-key"
                          maxlength="6"
                        />
                      </el-form-item>
                      <el-form-item>
                        <el-button type="primary">下一步</el-button>
                      </el-form-item>
                    </el-form>
                  </div>
                </el-card>

                <el-alert
                  title="重置密码说明"
                  type="info"
                  :closable="false"
                  show-icon
                  style="margin-top: 16px;"
                >
                  <template slot="default">
                    <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                      <li>输入注册时使用的邮箱地址</li>
                      <li>点击"发送验证码"按钮，系统会发送6位数字验证码到您的邮箱</li>
                      <li>输入收到的验证码，点击"下一步"</li>
                      <li>系统会自动生成新密码并发送到您的邮箱</li>
                      <li>使用新密码登录后，建议立即修改为个人密码</li>
                    </ul>
                  </template>
                </el-alert>
              </div>

              <!-- 修改密码 -->
              <el-divider content-position="left">修改密码</el-divider>
              <div class="step-detail">
                <h4>如何修改密码？</h4>
                <p>登录后，您可以通过以下方式修改密码：</p>
                
                <el-card shadow="never" style="margin: 16px 0; background: #f5f7fa;">
                  <div class="form-example">
                    <h5 style="margin-bottom: 16px; color: #606266;">修改密码表单示例：</h5>
                    <el-form
                      :model="changePasswordExample"
                      label-width="100px"
                      size="medium"
                      class="example-register-form"
                      disabled
                    >
                      <el-form-item label="原密码">
                        <el-input
                          v-model="changePasswordExample.oldPassword"
                          type="password"
                          placeholder="请输入当前密码"
                          prefix-icon="el-icon-lock"
                          show-password
                        />
                      </el-form-item>
                      <el-form-item label="新密码">
                        <el-input
                          v-model="changePasswordExample.newPassword"
                          type="password"
                          placeholder="请输入新密码（至少6位）"
                          prefix-icon="el-icon-lock"
                          show-password
                        />
                        <div class="password-strength" style="margin-top: 10px;">
                          <div class="strength-label">密码强度：</div>
                          <div class="strength-bar">
                            <div class="strength-item strong" style="width: 100%;"></div>
                          </div>
                          <span class="strength-text strong">强</span>
                        </div>
                      </el-form-item>
                      <el-form-item label="确认密码">
                        <el-input
                          v-model="changePasswordExample.confirmPassword"
                          type="password"
                          placeholder="请再次输入新密码"
                          prefix-icon="el-icon-lock"
                          show-password
                        />
                      </el-form-item>
                      <el-form-item>
                        <el-button type="primary">确认修改</el-button>
                        <el-button>重置</el-button>
                      </el-form-item>
                    </el-form>
                  </div>
                </el-card>

                <el-alert
                  title="修改密码说明"
                  type="info"
                  :closable="false"
                  show-icon
                  style="margin-top: 16px;"
                >
                  <template slot="default">
                    <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                      <li>登录系统后，访问"自助密码服务"页面，选择"修改密码"标签页</li>
                      <li>或访问"个人中心"，在个人信息页面修改密码</li>
                      <li>输入当前密码和新密码（至少6位）</li>
                      <li>密码修改成功后，系统会自动退出登录，需要使用新密码重新登录</li>
                      <li>建议定期更换密码，提高账户安全性</li>
                    </ul>
                  </template>
                </el-alert>
              </div>
            </el-card>
          </div>
        </section>

        <!-- 5. 管理员指南 -->
        <section id="admin-guide" class="help-section">
          <h2 class="section-title">
            <el-icon class="section-icon"><Setting /></el-icon>
            5. 管理员用户介绍
          </h2>
          <div class="section-content">
            <el-card shadow="hover" class="info-card">
              <h3>管理员权限概述</h3>
              <p>管理员拥有系统的最高权限，可以管理用户、组、权限等所有功能。主要功能包括：</p>

              <el-row :gutter="20" style="margin-top: 20px;">
                <el-col :span="12">
                  <el-card shadow="never" class="admin-feature-card">
                    <div class="admin-feature-icon">
                      <svg-icon icon-class="user" style="font-size: 40px; color: #409EFF;" />
                    </div>
                    <h4>用户管理</h4>
                    <ul>
                      <li>查看所有用户列表</li>
                      <li>创建、编辑、删除用户</li>
                      <li>审核用户注册申请</li>
                      <li>重置用户密码</li>
                      <li>启用/禁用用户账号</li>
                    </ul>
                  </el-card>
                </el-col>
                <el-col :span="12">
                  <el-card shadow="never" class="admin-feature-card">
                    <div class="admin-feature-icon">
                      <svg-icon icon-class="peoples" style="font-size: 40px; color: #67C23A;" />
                    </div>
                    <h4>组管理</h4>
                    <ul>
                      <li>创建和管理用户组</li>
                      <li>分配用户到组</li>
                      <li>设置组权限</li>
                      <li>管理组成员关系</li>
                      <li>...</li>
                    </ul>
                  </el-card>
                </el-col>
                <el-col :span="12" style="margin-top: 20px;">
                  <el-card shadow="never" class="admin-feature-card">
                    <div class="admin-feature-icon">
                      <svg-icon icon-class="lock" style="font-size: 40px; color: #E6A23C;" />
                    </div>
                    <h4>权限管理</h4>
                    <ul>
                      <li>配置sudo规则(免密、需要密码)</li>
                      <li>管理角色权限（sudo权限、docker组）</li>
                      <li>设置API访问权限</li>
                      <li>配置系统参数</li>
                    </ul>
                  </el-card>
                </el-col>
                <el-col :span="12" style="margin-top: 20px;">
                  <el-card shadow="never" class="admin-feature-card">
                    <div class="admin-feature-icon">
                      <svg-icon icon-class="monitor" style="font-size: 40px; color: #F56C6C;" />
                    </div>
                    <h4>系统监控</h4>
                    <ul>
                      <li>查看系统统计信息</li>
                      <li>查看操作日志</li>
                      <li>监控系统状态</li>
                      <li>管理SSH密钥</li>
                    </ul>
                  </el-card>
                </el-col>
              </el-row>

              <el-divider content-position="left" style="margin: 40px 0 24px 0;">审核用户注册申请</el-divider>
              
              <div class="admin-workflow">
                <h4>审核流程</h4>
                <el-steps :active="3" finish-status="success" align-center style="margin: 20px 0;">
                  <el-step title="查看待审核列表" description="进入用户管理-待审核用户" />
                  <el-step title="审核用户信息" description="查看用户注册信息" />
                  <el-step title="通过/拒绝" description="设置用户属性和权限" />
                  <el-step title="完成" description="用户收到通知邮件" />
                </el-steps>

                <el-card shadow="never" style="margin: 20px 0; background: #f5f7fa;">
                  <div class="form-example">
                    <h5 style="margin-bottom: 16px; color: #606266;">审核界面示例：</h5>
                    
                    <!-- 待审核用户列表 -->
                    <div style="margin-bottom: 24px;">
                      <h6 style="color: #909399; margin-bottom: 12px;">待审核用户列表</h6>
                      <el-table :data="pendingUsersExample" border style="width: 100%">
                        <el-table-column prop="username" label="用户名" width="120" />
                        <el-table-column prop="nickname" label="真实姓名" width="120" />
                        <el-table-column prop="mail" label="邮箱" />
                        <el-table-column prop="status" label="状态" width="100">
                          <template slot-scope="scope">
                            <el-tag type="warning" size="small">待审核</el-tag>
                          </template>
                        </el-table-column>
                        <el-table-column label="操作" width="100" align="center">
                          <template slot-scope="scope">
                            <el-button size="mini" type="success" icon="el-icon-check" circle />
                          </template>
                        </el-table-column>
                      </el-table>
                    </div>

                    <!-- 审核对话框 -->
                    <div>
                      <h6 style="color: #909399; margin-bottom: 12px;">审核对话框</h6>
                      <el-form
                        :model="reviewFormExample"
                        label-width="120px"
                        size="small"
                        class="example-register-form"
                        disabled
                      >
                        <el-alert
                          title="用户注册信息"
                          type="info"
                          :closable="false"
                          show-icon
                          style="margin-bottom: 16px;"
                        />
                        <el-row>
                          <el-col :span="12">
                            <el-form-item label="用户名">
                              <el-input v-model="reviewFormExample.username" disabled />
                            </el-form-item>
                          </el-col>
                          <el-col :span="12">
                            <el-form-item label="真实姓名">
                              <el-input v-model="reviewFormExample.nickname" disabled />
                            </el-form-item>
                          </el-col>
                          <el-col :span="12">
                            <el-form-item label="邮箱">
                              <el-input v-model="reviewFormExample.mail" disabled />
                            </el-form-item>
                          </el-col>
                        </el-row>

                        <el-divider content-position="left">审核设置</el-divider>
                        <el-row>
                          <el-col :span="12">
                            <el-form-item label="审核结果">
                              <el-radio-group v-model="reviewFormExample.status">
                                <el-radio :label="1">通过</el-radio>
                                <el-radio :label="2">拒绝</el-radio>
                              </el-radio-group>
                            </el-form-item>
                          </el-col>
                          <el-col :span="24">
                            <el-form-item label="审核备注">
                              <el-input
                                v-model="reviewFormExample.reviewRemark"
                                type="textarea"
                                :rows="3"
                                placeholder="请输入审核备注"
                              />
                            </el-form-item>
                          </el-col>
                        </el-row>

                        <el-divider content-position="left" v-if="reviewFormExample.status === 1">用户配置</el-divider>
                        <el-row v-if="reviewFormExample.status === 1">
                          <el-col :span="24">
                            <el-form-item label="分配分组">
                              <el-input placeholder="请选择分组（必选）" disabled />
                            </el-form-item>
                          </el-col>
                          <el-col :span="12">
                            <el-form-item label="分配角色">
                              <el-input placeholder="请选择角色" disabled />
                            </el-form-item>
                          </el-col>
                        </el-row>

                        <el-form-item>
                          <el-button type="primary">确认审核</el-button>
                          <el-button>取消</el-button>
                        </el-form-item>
                      </el-form>
                    </div>
                  </div>
                </el-card>

                <el-alert
                  title="审核注意事项"
                  type="warning"
                  :closable="false"
                  show-icon
                  style="margin-top: 16px;"
                >
                  <template slot="default">
                    <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                      <li>仔细核对用户注册信息的真实性和准确性</li>
                      <li>确保用户名符合命名规范，不与现有用户冲突</li>
                      <li>为通过审核的用户分配合适的分组和角色</li>
                      <li>填写审核备注，便于后续追溯</li>
                      <li>拒绝申请时，务必填写拒绝原因</li>
                      <li>审核通过后，系统会自动发送通知邮件给用户</li>
                    </ul>
                  </template>
                </el-alert>
              </div>

              <el-divider content-position="left" style="margin: 40px 0 24px 0;">管理员操作指南</el-divider>
              
              <el-collapse v-model="activeAdminCollapse" style="margin-top: 24px;">
                <el-collapse-item title="用户管理操作" name="user-management">
                  <div>
                    <h5>创建用户</h5>
                    <ol>
                      <li>进入"人员管理" → "用户管理"</li>
                      <li>点击"新增"按钮</li>
                      <li>填写用户信息（用户名、真实姓名、邮箱等）</li>
                      <li>设置用户密码</li>
                      <li>分配分组和角色</li>
                      <li>配置用户属性（UID、GID、Home目录等）</li>
                      <li>点击"确定"保存</li>
                    </ol>

                    <h5 style="margin-top: 20px;">编辑用户</h5>
                    <ol>
                      <li>在用户列表中点击"编辑"按钮</li>
                      <li>修改用户信息</li>
                      <li>更新分组、角色或权限</li>
                      <li>保存更改</li>
                    </ol>

                    <h5 style="margin-top: 20px;">重置用户密码</h5>
                    <ol>
                      <li>在用户列表中找到目标用户</li>
                      <li>点击"重置密码"按钮</li>
                      <li>系统会生成新密码并发送到用户邮箱</li>
                    </ol>
                  </div>
                </el-collapse-item>

                <el-collapse-item title="组管理操作" name="group-management">
                  <div>
                    <h5>创建组</h5>
                    <ol>
                      <li>进入"人员管理" → "组管理"</li>
                      <li>点击"新增"按钮</li>
                      <li>填写组信息（组名、组类型、GID等）</li>
                      <li>添加组成员</li>
                      <li>保存组配置</li>
                    </ol>

                    <h5 style="margin-top: 20px;">管理组成员</h5>
                    <ol>
                      <li>选择要管理的组</li>
                      <li>点击"成员管理"</li>
                      <li>添加或移除组成员</li>
                      <li>保存更改</li>
                    </ol>

                    <h5 style="margin-top: 20px;">Docker组管理</h5>
                    <p style="margin-bottom: 12px; color: #606266;">
                      Docker组是系统预置的posixGroup（GID=984），用于授予用户docker命令执行权限。
                    </p>
                    <ol>
                      <li>在组管理列表中找到"docker"组</li>
                      <li>点击"成员管理"按钮</li>
                      <li>在成员管理界面中添加需要docker权限的用户</li>
                      <li>保存后，用户将自动加入docker posixGroup</li>
                      <li>用户需要重新登录或执行<code>newgrp docker</code>使权限生效</li>
                    </ol>
                    <el-card shadow="never" style="margin-top: 16px; background: #f5f7fa;">
                      <div style="padding: 12px;">
                        <h6 style="margin: 0 0 8px 0; color: #606266; font-size: 14px;">Docker组配置示例：</h6>
                        <ul style="margin: 0; padding-left: 20px; color: #606266; font-size: 13px;">
                          <li><strong>组名</strong>：docker</li>
                          <li><strong>组类型</strong>：posix（Linux系统组）</li>
                          <li><strong>GID号</strong>：984</li>
                          <li><strong>作用</strong>：成员可直接使用docker命令，无需sudo权限</li>
                        </ul>
                      </div>
                    </el-card>
                  </div>
                </el-collapse-item>

                <el-collapse-item title="权限管理操作" name="permission-management">
                  <div>
                    <h5>配置Sudo规则</h5>
                    <ol>
                      <li>进入"Sudo管理"</li>
                      <li>创建新的sudo规则</li>
                      <li>设置允许的用户、主机和命令</li>
                      <li>配置是否需要密码验证</li>
                      <li>保存规则</li>
                    </ol>

                    <h5 style="margin-top: 20px;">管理角色权限</h5>
                    <ol>
                      <li>进入"系统管理" → "角色管理"</li>
                      <li>选择要配置的角色</li>
                      <li>分配API权限和菜单权限</li>
                      <li>保存权限配置</li>
                    </ol>

                    <h5 style="margin-top: 20px;">Docker组权限管理</h5>
                    <p style="margin-bottom: 12px; color: #606266;">
                      Docker组是一个特殊的posixGroup（GID=984），成员可以使用docker命令而无需sudo权限。
                    </p>
                    <ol>
                      <li>进入"人员管理" → "组管理"</li>
                      <li>找到或创建"docker"组（类型为posix，GID=984）</li>
                      <li>点击"成员管理"按钮</li>
                      <li>在成员管理界面中，添加需要docker权限的用户</li>
                      <li>保存后，用户将自动加入docker组</li>
                    </ol>
                    <el-alert
                      title="Docker组说明"
                      type="info"
                      :closable="false"
                      show-icon
                      style="margin-top: 16px;"
                    >
                      <template slot="default">
                        <ul style="margin: 10px 0 0 0; padding-left: 20px;">
                          <li>Docker组是posixGroup类型，GID固定为984</li>
                          <li>用户加入docker组后，可以直接使用<code>docker</code>命令，无需sudo</li>
                          <li>例如：<code>docker ps</code>、<code>docker run</code>等命令可以直接执行</li>
                          <li>组成员信息会同步到LDAP，在所有服务器上生效</li>
                          <li>用户需要重新登录或执行<code>newgrp docker</code>才能使组权限生效</li>
                        </ul>
                      </template>
                    </el-alert>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </el-card>
          </div>
        </section>

        <!-- 底部联系信息 -->
        <div class="help-footer">
          <el-card shadow="never">
            <h3>需要帮助？</h3>
            <p>如果您在使用过程中遇到任何问题，请联系系统管理员或查看系统日志。</p>
            <p style="margin-top: 12px; color: #909399; font-size: 14px;">
              <el-icon><InfoFilled /></el-icon>
              本文档会持续更新，请关注最新版本
            </p>
          </el-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Help',
  data() {
    return {
      // 示例数据
      validUsernamesExample: [
        { username: 'wwtian', nickname: 'tianweiwei', status: '未注册', remark: '田维伟' },
        { username: 'zhsan', nickname: 'zhangsan', status: '已注册-有效', remark: '张三' },
        { username: 'lisi', nickname: 'lisi', status: '已注册-失效', remark: '李四' },
        { username: 'wangwu', nickname: 'wangwu', status: '待审核', remark: '王五' }
      ],
      exampleForm: {
        username: 'wwtian',
        email: 'tianweiwei@cimrbj.ac.cn',
        realName: 'WeiweiTian',
        password: '••••••••',
        confirmPassword: '••••••••',
        remark: '田维伟'
      },
      forgotPasswordExample: {
        mail: 'tianweiwei@cimrbj.ac.cn',
        code: '123456'
      },
      changePasswordExample: {
        oldPassword: '••••••••',
        newPassword: '••••••••',
        confirmPassword: '••••••••'
      },
      pendingUsersExample: [
        { username: 'wwtian', nickname: 'tianweiwei', mail: 'tianweiwei@cimrbj.ac.cn', status: 0 }
      ],
      reviewFormExample: {
        username: 'wwtian',
        nickname: 'tianweiwei',
        mail: 'tianweiwei@cimrbj.ac.cn',
        status: 1,
        reviewRemark: '审核通过，已分配实验室成员组'
      },
      activeAdminCollapse: []
    }
  }
}
</script>

<style lang="scss" scoped>
.help-container {
  min-height: calc(100vh - 84px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
}

.help-content {
  max-width: 1200px;
  margin: 0 auto;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.help-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  padding: 40px;
  text-align: center;

  .help-title {
    font-size: 32px;
    margin: 0 0 12px 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .help-subtitle {
    font-size: 16px;
    margin: 0;
    opacity: 0.9;
  }
}

.help-nav {
  background: #f5f7fa;
  padding: 20px 40px;
  border-bottom: 1px solid #e4e7ed;

  ::v-deep .el-anchor {
    display: flex;
    justify-content: center;
    gap: 30px;
  }
}

.help-body {
  padding: 40px;
}

.help-section {
  margin-bottom: 60px;

  &:last-child {
    margin-bottom: 0;
  }
}

.section-title {
  font-size: 24px;
  color: #303133;
  margin: 0 0 24px 0;
  display: flex;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 2px solid #e4e7ed;

  .section-icon {
    margin-right: 12px;
    font-size: 28px;
    color: #409EFF;
  }
}

.section-content {
  .info-card {
    ::v-deep .el-card__body {
      padding: 30px;
    }

    h3 {
      font-size: 20px;
      color: #303133;
      margin: 0 0 16px 0;
    }

    h4 {
      font-size: 18px;
      color: #606266;
      margin: 24px 0 12px 0;
    }

    p {
      color: #606266;
      line-height: 1.8;
      margin: 12px 0;
    }

    ul, ol {
      color: #606266;
      line-height: 1.8;
      margin: 12px 0;
      padding-left: 24px;
    }

    code {
      background: #f5f7fa;
      padding: 2px 6px;
      border-radius: 3px;
      font-family: 'Courier New', monospace;
      color: #e6a23c;
    }
  }
}

.feature-card {
  text-align: center;
  height: 100%;

  .feature-icon {
    margin-bottom: 16px;
  }

  h4 {
    font-size: 18px;
    color: #303133;
    margin: 0 0 12px 0;
  }

  p {
    color: #606266;
    line-height: 1.6;
    margin: 0;
  }
}

.benefit-card {
  text-align: center;
  height: 100%;

  .benefit-icon {
    margin-bottom: 16px;
  }

  h4 {
    font-size: 18px;
    color: #303133;
    margin: 0 0 12px 0;
  }

  p {
    color: #606266;
    line-height: 1.6;
    margin: 0;
  }
}

.server-architecture {
  margin: 30px 0;

  .architecture-diagram {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 30px;

    .ldap-server {
      .ldap-box {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: #fff;
        padding: 30px;
        border-radius: 12px;
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);

        h4 {
          margin: 12px 0 8px 0;
          color: #fff;
          font-weight: 600;
          text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }

        p {
          margin: 0;
          color: rgba(255, 255, 255, 0.95);
          font-size: 14px;
          text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }
      }
    }

    .connection-lines {
      display: flex;
      gap: 20px;
      width: 100%;
      justify-content: center;

      .line {
        width: 2px;
        height: 40px;
        background: linear-gradient(to bottom, #409EFF, transparent);
      }
    }

    .compute-servers {
      display: flex;
      gap: 30px;
      width: 100%;
      justify-content: center;

      .hpc-box {
        background: #f5f7fa;
        border: 2px solid #409EFF;
        padding: 24px;
        border-radius: 8px;
        flex: 1;
        max-width: 200px;
        text-align: center;
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
          transform: translateY(-4px);
        }

        h4 {
          margin: 12px 0 8px 0;
          color: #303133;
        }

        p {
          margin: 0;
          color: #909399;
          font-size: 14px;
        }
      }
    }

    .server-box {
      text-align: center;

      h4 {
        margin: 12px 0 8px 0;
      }

      p {
        margin: 0;
        font-size: 14px;
      }
    }
  }
}

.storage-diagram {
  .storage-visual {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;

    .storage-center {
      background: linear-gradient(135deg, #409EFF 0%, #67C23A 100%);
      color: #fff;
      padding: 40px;
      border-radius: 12px;
      text-align: center;
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);

      h4 {
        margin: 16px 0 8px 0;
        color: #fff;
      }

      p {
        margin: 0;
        font-size: 16px;
        opacity: 0.9;
      }
    }

    .storage-connections {
      display: flex;
      gap: 40px;
      margin-top: 20px;

      .storage-link {
        padding: 12px 24px;
        background: #f5f7fa;
        border: 2px solid #409EFF;
        border-radius: 6px;
        font-weight: 500;
        color: #409EFF;
        transition: all 0.3s;

        &:hover {
          background: #409EFF;
          color: #fff;
        }
      }
    }
  }
}

.step-detail {
  margin-top: 40px;

  &:first-of-type {
    margin-top: 0;
  }

  h4 {
    font-size: 18px;
    color: #303133;
    margin: 0 0 12px 0;
    display: flex;
    align-items: center;
    gap: 12px;
  }

  p {
    color: #606266;
    line-height: 1.8;
    margin: 12px 0;
  }
}

.form-example {
  padding: 20px 30px;
  background: #fff;
  border-radius: 4px;

  .example-register-form {
    ::v-deep .el-form-item__label {
      font-weight: 500;
      color: #606266;
    }

    .form-tip {
      font-size: 12px;
      color: #909399;
      margin-top: 5px;
      line-height: 1.5;
    }

    .password-strength {
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
  }
}

.admin-feature-card {
  min-height: 240px;
  display: flex;
  flex-direction: column;

  .admin-feature-icon {
    margin-bottom: 16px;
    text-align: center;
  }

  h4 {
    font-size: 18px;
    color: #303133;
    margin: 0 0 16px 0;
    text-align: center;
  }

  ul {
    margin: 0;
    padding-left: 20px;
    color: #606266;
    line-height: 1.8;
    flex: 1;

    li {
      margin: 8px 0;
    }
  }
}

.admin-workflow {
  margin-top: 40px;

  h4 {
    font-size: 18px;
    color: #303133;
    margin: 0 0 24px 0;
  }

  > * {
    margin-bottom: 32px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

.help-footer {
  margin-top: 60px;
  padding-top: 40px;
  border-top: 2px solid #e4e7ed;

  h3 {
    font-size: 20px;
    color: #303133;
    margin: 0 0 12px 0;
  }

  p {
    color: #606266;
    line-height: 1.8;
    margin: 8px 0;
  }
}

@media (max-width: 768px) {
  .help-body {
    padding: 20px;
  }

  .section-title {
    font-size: 20px;
  }

  .server-architecture .architecture-diagram .compute-servers {
    flex-direction: column;
  }
}
</style>