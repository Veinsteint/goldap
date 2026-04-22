package common

import (
	"errors"
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/tools"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

// InitData initializes default data in MySQL
func InitData() {
	if !config.Conf.System.InitData {
		return
	}

	// Sync existing users to pre-config
	syncExistingUsersToPreConfig()

	// Initialize roles
	newRoles := make([]*model.Role, 0)
	roles := []*model.Role{
		{Model: gorm.Model{ID: 1}, Name: "Administrator", Keyword: "admin", Sort: 1, Status: 1, Creator: "system"},
		{Model: gorm.Model{ID: 2}, Name: "User", Keyword: "user", Sort: 3, Status: 1, Creator: "system"},
		{Model: gorm.Model{ID: 3}, Name: "Guest", Keyword: "guest", Sort: 5, Status: 1, Creator: "system"},
	}

	for _, role := range roles {
		err := DB.First(&role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) > 0 {
		if err := DB.Create(&newRoles).Error; err != nil {
			Log.Errorf("Failed to create roles: %v", err)
		}
	}

	// Initialize menus
	newMenus := make([]model.Menu, 0)
	var uint0 uint = 0
	var uint1 uint = 1
	var uint4 uint = 5
	var uint8 uint = 9
	var uint12 uint = 11

	menus := []model.Menu{
		{Model: gorm.Model{ID: 1}, Name: "UserManage", Title: "用户管理", Icon: "user", Path: "/personnel", Component: "Layout", Redirect: "/personnel/user", Sort: 5, ParentId: uint0, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 2}, Name: "User", Title: "用户列表", Icon: "people", Path: "user", Component: "/personnel/user/index", Sort: 6, ParentId: uint1, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 13}, Name: "UserPreConfig", Title: "用户配置", Icon: "el-icon-setting", Path: "userPreConfig", Component: "/personnel/userPreConfig/index", Sort: 7, ParentId: uint1, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 3}, Name: "Group", Title: "分组管理", Icon: "peoples", Path: "group", Component: "/personnel/group/index", Sort: 8, ParentId: uint1, NoCache: 1, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 4}, Name: "FieldRelation", Title: "字段映射", Icon: "el-icon-s-tools", Path: "fieldRelation", Component: "/personnel/fieldRelation/index", Sort: 9, ParentId: uint1, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 5}, Name: "System", Title: "系统管理", Icon: "component", Path: "/system", Component: "Layout", Redirect: "/system/role", Sort: 10, ParentId: uint0, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 6}, Name: "Role", Title: "角色管理", Icon: "eye-open", Path: "role", Component: "/system/role/index", Sort: 11, ParentId: uint4, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 7}, Name: "Menu", Title: "菜单管理", Icon: "tree-table", Path: "menu", Component: "/system/menu/index", Sort: 13, ParentId: uint4, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 8}, Name: "Api", Title: "接口管理", Icon: "tree", Path: "api", Component: "/system/api/index", Sort: 14, ParentId: uint4, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 14}, Name: "SystemConfig", Title: "系统维护", Icon: "el-icon-setting", Path: "systemConfig", Component: "/system/systemConfig/index", Sort: 15, ParentId: uint4, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 9}, Name: "Log", Title: "日志管理", Icon: "documentation", Path: "/log", Component: "Layout", Redirect: "/log/operation-log", Sort: 20, ParentId: uint0, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 10}, Name: "OperationLog", Title: "操作日志", Icon: "form", Path: "operation-log", Component: "/log/operation-log/index", Sort: 21, ParentId: uint8, Roles: roles[:1], Creator: "system"},
		{Model: gorm.Model{ID: 11}, Name: "Service", Title: "自助服务", Icon: "example", Path: "/service", Component: "Layout", Redirect: "/service/passwordService", Sort: 30, ParentId: uint0, Roles: roles[:2], Creator: "system"},
		{Model: gorm.Model{ID: 12}, Name: "Password", Title: "密码服务", Icon: "lock", Path: "passwordService", Component: "/service/passwordService/index", Sort: 31, ParentId: uint12, Roles: roles[:2], Creator: "system"},
	}

	for _, menu := range menus {
		var existingMenu model.Menu
		err := DB.First(&existingMenu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		} else {
			menu.Model = existingMenu.Model
			_ = DB.Model(&menu).Updates(menu).Error
			_ = DB.Model(&menu).Association("Roles").Replace(menu.Roles)
		}
	}
	if len(newMenus) > 0 {
		if err := DB.Create(&newMenus).Error; err != nil {
			Log.Errorf("Failed to create menus: %v", err)
		}
	}

	// Initialize admin user
	newUsers := make([]*model.User, 0)
	users := []*model.User{
		{
			Model:         gorm.Model{ID: 1},
			Username:      "admin",
			Password:      tools.NewGenPasswd(config.Conf.Ldap.AdminPass),
			Nickname:      "Administrator",
			GivenName:     "Admin",
			Mail:          "admin@" + config.Conf.Ldap.DefaultEmailSuffix,
			JobNumber:     "ADMIN001",
			Mobile:        "",
			PostalAddress: "",
			Departments:   "System",
			Position:      "Administrator",
			Introduction:  "System Administrator",
			Status:        1,
			Creator:       "system",
			Roles:         roles[:1],
			UserDN:        config.Conf.Ldap.AdminDN,
		},
	}

	for _, user := range users {
		err := DB.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		if err := DB.Create(&newUsers).Error; err != nil {
			Log.Errorf("Failed to create users: %v", err)
		}
	}

	// Initialize APIs
	apis := []model.Api{
		{Method: "POST", Path: "/base/login", Category: "base", Remark: "User login", Creator: "system"},
		{Method: "POST", Path: "/base/logout", Category: "base", Remark: "User logout", Creator: "system"},
		{Method: "POST", Path: "/base/refreshToken", Category: "base", Remark: "Refresh JWT token", Creator: "system"},
		{Method: "POST", Path: "/base/sendcode", Category: "base", Remark: "Send verification code", Creator: "system"},
		{Method: "POST", Path: "/base/changePwd", Category: "base", Remark: "Change password via email", Creator: "system"},
		{Method: "GET", Path: "/base/dashboard", Category: "base", Remark: "Dashboard data", Creator: "system"},
		{Method: "GET", Path: "/user/info", Category: "user", Remark: "Get current user info", Creator: "system"},
		{Method: "GET", Path: "/user/list", Category: "user", Remark: "Get user list", Creator: "system"},
		{Method: "POST", Path: "/user/changePwd", Category: "user", Remark: "Change password", Creator: "system"},
		{Method: "POST", Path: "/user/resetPassword", Category: "user", Remark: "Reset password", Creator: "system"},
		{Method: "POST", Path: "/user/add", Category: "user", Remark: "Create user", Creator: "system"},
		{Method: "POST", Path: "/user/update", Category: "user", Remark: "Update user", Creator: "system"},
		{Method: "POST", Path: "/user/delete", Category: "user", Remark: "Delete users", Creator: "system"},
		{Method: "POST", Path: "/user/changeUserStatus", Category: "user", Remark: "Change user status", Creator: "system"},
		{Method: "POST", Path: "/user/syncOpenLdapUsers", Category: "user", Remark: "Sync users from LDAP", Creator: "system"},
		{Method: "POST", Path: "/user/syncSqlUsers", Category: "user", Remark: "Sync users to LDAP", Creator: "system"},
		{Method: "GET", Path: "/user/ssh-keys", Category: "user", Remark: "Get SSH keys", Creator: "system"},
		{Method: "POST", Path: "/user/ssh-keys", Category: "user", Remark: "Add SSH key", Creator: "system"},
		{Method: "DELETE", Path: "/user/ssh-keys/:id", Category: "user", Remark: "Delete SSH key", Creator: "system"},
		{Method: "GET", Path: "/user/pending/list", Category: "user", Remark: "Get pending users", Creator: "system"},
		{Method: "POST", Path: "/user/pending/review", Category: "user", Remark: "Review pending user", Creator: "system"},
		{Method: "POST", Path: "/user/pending/delete", Category: "user", Remark: "Delete pending user", Creator: "system"},
		{Method: "GET", Path: "/user/preconfig/list", Category: "user", Remark: "Get user pre-configs", Creator: "system"},
		{Method: "POST", Path: "/user/preconfig/add", Category: "user", Remark: "Add user pre-config", Creator: "system"},
		{Method: "POST", Path: "/user/preconfig/update", Category: "user", Remark: "Update user pre-config", Creator: "system"},
		{Method: "POST", Path: "/user/preconfig/delete", Category: "user", Remark: "Delete user pre-config", Creator: "system"},
		{Method: "GET", Path: "/user/preconfig/getByUsername", Category: "user", Remark: "Get pre-config by username", Creator: "system"},
		{Method: "POST", Path: "/user/preconfig/syncUsers", Category: "user", Remark: "Sync existing users to pre-config", Creator: "system"},
		{Method: "GET", Path: "/group/list", Category: "group", Remark: "Get group list", Creator: "system"},
		{Method: "GET", Path: "/group/tree", Category: "group", Remark: "Get group tree", Creator: "system"},
		{Method: "POST", Path: "/group/add", Category: "group", Remark: "Create group", Creator: "system"},
		{Method: "POST", Path: "/group/update", Category: "group", Remark: "Update group", Creator: "system"},
		{Method: "POST", Path: "/group/delete", Category: "group", Remark: "Delete groups", Creator: "system"},
		{Method: "POST", Path: "/group/adduser", Category: "group", Remark: "Add user to group", Creator: "system"},
		{Method: "POST", Path: "/group/removeuser", Category: "group", Remark: "Remove user from group", Creator: "system"},
		{Method: "GET", Path: "/group/useringroup", Category: "group", Remark: "Get users in group", Creator: "system"},
		{Method: "GET", Path: "/group/usernoingroup", Category: "group", Remark: "Get users not in group", Creator: "system"},
		{Method: "GET", Path: "/group-user-permission", Category: "group", Remark: "Get group user permissions", Creator: "system"},
		{Method: "POST", Path: "/group-user-permission", Category: "group", Remark: "Add group user permission", Creator: "system"},
		{Method: "PUT", Path: "/group-user-permission", Category: "group", Remark: "Update group user permission", Creator: "system"},
		{Method: "DELETE", Path: "/group-user-permission", Category: "group", Remark: "Delete group user permission", Creator: "system"},
		{Method: "POST", Path: "/group/syncOpenLdapDepts", Category: "group", Remark: "Sync groups from LDAP", Creator: "system"},
		{Method: "POST", Path: "/group/syncSqlGroups", Category: "group", Remark: "Sync groups to LDAP", Creator: "system"},
		{Method: "GET", Path: "/role/list", Category: "role", Remark: "Get role list", Creator: "system"},
		{Method: "POST", Path: "/role/add", Category: "role", Remark: "Create role", Creator: "system"},
		{Method: "POST", Path: "/role/update", Category: "role", Remark: "Update role", Creator: "system"},
		{Method: "GET", Path: "/role/getmenulist", Category: "role", Remark: "Get role menus", Creator: "system"},
		{Method: "POST", Path: "/role/updatemenus", Category: "role", Remark: "Update role menus", Creator: "system"},
		{Method: "GET", Path: "/role/getapilist", Category: "role", Remark: "Get role APIs", Creator: "system"},
		{Method: "POST", Path: "/role/updateapis", Category: "role", Remark: "Update role APIs", Creator: "system"},
		{Method: "POST", Path: "/role/delete", Category: "role", Remark: "Delete roles", Creator: "system"},
		{Method: "GET", Path: "/menu/tree", Category: "menu", Remark: "Get menu tree", Creator: "system"},
		{Method: "GET", Path: "/menu/access/tree", Category: "menu", Remark: "Get user menu tree", Creator: "system"},
		{Method: "POST", Path: "/menu/add", Category: "menu", Remark: "Create menu", Creator: "system"},
		{Method: "POST", Path: "/menu/update", Category: "menu", Remark: "Update menu", Creator: "system"},
		{Method: "POST", Path: "/menu/delete", Category: "menu", Remark: "Delete menus", Creator: "system"},
		{Method: "GET", Path: "/api/list", Category: "api", Remark: "Get API list", Creator: "system"},
		{Method: "GET", Path: "/api/tree", Category: "api", Remark: "Get API tree", Creator: "system"},
		{Method: "POST", Path: "/api/add", Category: "api", Remark: "Create API", Creator: "system"},
		{Method: "POST", Path: "/api/update", Category: "api", Remark: "Update API", Creator: "system"},
		{Method: "POST", Path: "/api/delete", Category: "api", Remark: "Delete APIs", Creator: "system"},
		{Method: "GET", Path: "/system/config/get", Category: "system", Remark: "Get system configuration", Creator: "system"},
		{Method: "POST", Path: "/system/config/update", Category: "system", Remark: "Update system configuration", Creator: "system"},
		{Method: "GET", Path: "/fieldrelation/list", Category: "fieldrelation", Remark: "Get field relations", Creator: "system"},
		{Method: "POST", Path: "/fieldrelation/add", Category: "fieldrelation", Remark: "Create field relation", Creator: "system"},
		{Method: "POST", Path: "/fieldrelation/update", Category: "fieldrelation", Remark: "Update field relation", Creator: "system"},
		{Method: "POST", Path: "/fieldrelation/delete", Category: "fieldrelation", Remark: "Delete field relations", Creator: "system"},
		{Method: "GET", Path: "/log/operation/list", Category: "log", Remark: "Get operation logs", Creator: "system"},
		{Method: "POST", Path: "/log/operation/delete", Category: "log", Remark: "Delete operation logs", Creator: "system"},
		{Method: "DELETE", Path: "/log/operation/clean", Category: "log", Remark: "Clean operation logs", Creator: "system"},
	}

	newApi := make([]model.Api, 0)
	newRoleCasbin := make([]model.RoleCasbin, 0)
	basePaths := []string{
		"/base/login", "/base/logout", "/base/refreshToken", "/base/sendcode",
		"/base/changePwd", "/base/dashboard", "/user/info", "/user/changePwd",
		"/user/ssh-keys", "/user/ssh-keys/:id", "/menu/access/tree", "/log/operation/list",
	}

	for _, api := range apis {
		var existingApi model.Api
		err := DB.Where("method = ? AND path = ?", api.Method, api.Path).First(&existingApi).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ID = 0
			newApi = append(newApi, api)
		}

		// Admin permissions
		adminPolicies := CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword, api.Path, api.Method)
		if len(adminPolicies) == 0 {
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{Keyword: roles[0].Keyword, Path: api.Path, Method: api.Method})
		}

		// User base permissions
		isBasePath := funk.ContainsString(basePaths, api.Path) || strings.HasPrefix(api.Path, "/user/ssh-keys")
		if isBasePath {
			userPolicies := CasbinEnforcer.GetFilteredPolicy(0, roles[1].Keyword, api.Path, api.Method)
			if len(userPolicies) == 0 {
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{Keyword: roles[1].Keyword, Path: api.Path, Method: api.Method})
			}
		}
	}

	if len(newApi) > 0 {
		if err := DB.Create(&newApi).Error; err != nil {
			Log.Errorf("Failed to create APIs: %v", err)
		}
	}

	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{c.Keyword, c.Path, c.Method})
		}
		isAdd, err := CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			Log.Errorf("Failed to create Casbin policies: %v", err)
		}
	}

	// Initialize groups
	newGroups := make([]model.Group, 0)
	groups := []model.Group{
		{Model: gorm.Model{ID: 6}, GroupName: "CMPLabHPC", Remark: "Default group", Creator: "system", GroupType: "ou", ParentId: 0, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_0", GroupDN: fmt.Sprintf("ou=%s,%s", "CMPLabHPC", config.Conf.Ldap.BaseDN), IPRanges: `["192.168.11.6-192.168.11.8"]`},
		{Model: gorm.Model{ID: 7}, GroupName: "sudoers", Remark: "Sudo rules", Creator: "system", GroupType: "ou", ParentId: 0, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_0", GroupDN: fmt.Sprintf("ou=%s,%s", "sudoers", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 8}, GroupName: "sudouser-nopasswd", Remark: "Passwordless sudo", Creator: "system", GroupType: "cn", ParentId: 7, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_0", GroupDN: fmt.Sprintf("cn=%s,ou=%s,%s", "sudouser-nopasswd", "sudoers", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 9}, GroupName: "sudouser-other", Remark: "Sudo with password", Creator: "system", GroupType: "cn", ParentId: 7, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_0", GroupDN: fmt.Sprintf("cn=%s,ou=%s,%s", "sudouser-other", "sudoers", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 10}, GroupName: "docker", Remark: "Docker用户组(posixGroup GID=984)，成员可使用docker命令", Creator: "system", GroupType: "posix", GIDNumber: 984, ParentId: 0, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_0", GroupDN: fmt.Sprintf("cn=%s,%s", "docker", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 11}, GroupName: "ClusterServer", Remark: "机房服务器集群登录域", Creator: "system", GroupType: "cn", ParentId: 6, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_6", GroupDN: fmt.Sprintf("cn=%s,ou=%s,%s", "ClusterServer", "CMPLabHPC", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 12}, GroupName: "ClusterHPC", Remark: "实验室HPC集群登录域", Creator: "system", GroupType: "cn", ParentId: 6, SourceDeptId: "platform_0", Source: "platform", SourceDeptParentId: "platform_6", GroupDN: fmt.Sprintf("cn=%s,ou=%s,%s", "ClusterHPC", "CMPLabHPC", config.Conf.Ldap.BaseDN)},
	}

	for _, group := range groups {
		err := DB.First(&group, group.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newGroups = append(newGroups, group)
		}
	}
	if len(newGroups) > 0 {
		if err := DB.Create(&newGroups).Error; err != nil {
			Log.Errorf("Failed to create groups: %v", err)
		} else {
			for _, group := range newGroups {
				ensureGroupInLDAP(group)
			}
		}
	}

	// Ensure default sudo groups exist in LDAP
	defaultSudoGroups := []model.Group{
		{Model: gorm.Model{ID: 7}, GroupName: "sudoers", Remark: "Sudo rules", Creator: "system", GroupType: "ou", GroupDN: fmt.Sprintf("ou=%s,%s", "sudoers", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 8}, GroupName: "sudouser-nopasswd", Remark: "Passwordless sudo", Creator: "system", GroupType: "cn", ParentId: 7, GroupDN: fmt.Sprintf("cn=%s,ou=%s,%s", "sudouser-nopasswd", "sudoers", config.Conf.Ldap.BaseDN)},
		{Model: gorm.Model{ID: 9}, GroupName: "sudouser-other", Remark: "Sudo with password", Creator: "system", GroupType: "cn", ParentId: 7, GroupDN: fmt.Sprintf("cn=%s,ou=%s,%s", "sudouser-other", "sudoers", config.Conf.Ldap.BaseDN)},
	}

	for _, group := range defaultSudoGroups {
		ensureGroupInLDAP(group)
	}
}

// ensureGroupInLDAP ensures group exists in LDAP
func ensureGroupInLDAP(group model.Group) {
	searchConn, err := GetLDAPConn()
	if err != nil {
		return
	}
	defer PutLADPConn(searchConn)

	searchRequest := ldap.NewSearchRequest(group.GroupDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
	_, err = searchConn.Search(searchRequest)
	if err == nil {
		return
	}

	isSudoRole := strings.Contains(group.GroupDN, "ou=sudoers,") && group.GroupType == "cn"
	isPosixGroup := group.GroupType == "posix"

	modifyConn, err := GetLDAPConnForModify()
	if err != nil {
		return
	}
	defer modifyConn.Close()

	if isSudoRole {
		// Ensure ou=sudoers exists
		sudoersDN := fmt.Sprintf("ou=sudoers,%s", config.Conf.Ldap.BaseDN)
		sudoersSearchRequest := ldap.NewSearchRequest(sudoersDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
		_, sudoersErr := searchConn.Search(sudoersSearchRequest)
		if sudoersErr != nil {
			sudoersAdd := ldap.NewAddRequest(sudoersDN, nil)
			sudoersAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
			sudoersAdd.Attribute("ou", []string{"sudoers"})
			sudoersAdd.Attribute("description", []string{"Sudo rules container"})
			if addErr := modifyConn.Add(sudoersAdd); addErr != nil {
				if ldapErr, ok := addErr.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
					Log.Warnf("Failed to create ou=sudoers: %v", addErr)
				}
			}
		}

		// Create sudoRole - use "ALL" as sudoUser for default sudo groups if no users
		var groupWithUsers model.Group
		err = DB.Preload("Users").First(&groupWithUsers, group.ID).Error
		sudoUser := "ALL"
		if err == nil && len(groupWithUsers.Users) > 0 {
			sudoUser = groupWithUsers.Users[0].Username
		}
		
		sudoAdd := ldap.NewAddRequest(group.GroupDN, nil)
		sudoAdd.Attribute("objectClass", []string{"top", "sudoRole"})
		sudoAdd.Attribute("cn", []string{group.GroupName})
		if group.Remark != "" {
			sudoAdd.Attribute("description", []string{group.Remark})
		}
		sudoAdd.Attribute("sudoUser", []string{sudoUser})
		sudoAdd.Attribute("sudoHost", []string{"ALL"})
		sudoAdd.Attribute("sudoCommand", []string{"ALL"})
		sudoAdd.Attribute("sudoRunAsUser", []string{"ALL"})
		if group.GroupName == "sudouser-nopasswd" {
			sudoAdd.Attribute("sudoOption", []string{"!authenticate"})
		}
		if addErr := modifyConn.Add(sudoAdd); addErr != nil {
			if ldapErr, ok := addErr.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				// Group already exists in LDAP
				Log.Infof("SudoRole %s already exists in LDAP", group.GroupDN)
			} else {
				Log.Errorf("Failed to create sudoRole %s: %v", group.GroupDN, addErr)
			}
		}
	} else if isPosixGroup {
		// Ensure parent OU exists first
		parentDN := extractParentDN(group.GroupDN)
		if parentDN != "" && parentDN != config.Conf.Ldap.BaseDN {
			parentSearchRequest := ldap.NewSearchRequest(parentDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
			_, parentErr := searchConn.Search(parentSearchRequest)
			if parentErr != nil {
				ouName := extractOUName(parentDN)
				parentAdd := ldap.NewAddRequest(parentDN, nil)
				parentAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
				parentAdd.Attribute("ou", []string{ouName})
				_ = modifyConn.Add(parentAdd)
			}
		}

		// Create posixGroup
		posixAdd := ldap.NewAddRequest(group.GroupDN, nil)
		posixAdd.Attribute("objectClass", []string{"posixGroup", "top"})
		posixAdd.Attribute("cn", []string{group.GroupName})
		posixAdd.Attribute("gidNumber", []string{fmt.Sprintf("%d", group.GIDNumber)})
		if group.Remark != "" {
			posixAdd.Attribute("description", []string{group.Remark})
		}
		// Add members if any
		var groupWithUsers model.Group
		err = DB.Preload("Users").First(&groupWithUsers, group.ID).Error
		if err == nil && len(groupWithUsers.Users) > 0 {
			memberUids := make([]string, 0)
			for _, u := range groupWithUsers.Users {
				memberUids = append(memberUids, u.Username)
			}
			posixAdd.Attribute("memberUid", memberUids)
		}
		if addErr := modifyConn.Add(posixAdd); addErr != nil {
			if ldapErr, ok := addErr.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				// Group already exists in LDAP
				Log.Infof("PosixGroup %s already exists in LDAP", group.GroupDN)
			} else {
				Log.Errorf("Failed to create posixGroup %s: %v", group.GroupDN, addErr)
			}
		}
	} else {
		// Ensure parent OU exists first for cn groups
		if group.GroupType == "cn" {
			parentDN := extractParentDN(group.GroupDN)
			if parentDN != "" && parentDN != config.Conf.Ldap.BaseDN {
				parentSearchRequest := ldap.NewSearchRequest(parentDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
				_, parentErr := searchConn.Search(parentSearchRequest)
				if parentErr != nil {
					ouName := extractOUName(parentDN)
					if ouName != "" {
						parentAdd := ldap.NewAddRequest(parentDN, nil)
						parentAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
						parentAdd.Attribute("ou", []string{ouName})
						if addErr := modifyConn.Add(parentAdd); addErr != nil {
							if ldapErr, ok := addErr.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
								Log.Warnf("Failed to create parent DN %s: %v", parentDN, addErr)
							}
						}
					}
				}
			}
		}

		ouAdd := ldap.NewAddRequest(group.GroupDN, nil)
		if group.GroupType == "ou" {
			ouAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
			ouAdd.Attribute("ou", []string{group.GroupName})
		} else if group.GroupType == "cn" {
			ouAdd.Attribute("objectClass", []string{"groupOfUniqueNames", "top"})
			ouAdd.Attribute("cn", []string{group.GroupName})
			var groupWithUsers model.Group
			err = DB.Preload("Users").First(&groupWithUsers, group.ID).Error
			if err == nil && len(groupWithUsers.Users) > 0 {
				uniqueMembers := make([]string, 0)
				for _, u := range groupWithUsers.Users {
					if u.UserDN != "" {
						uniqueMembers = append(uniqueMembers, u.UserDN)
					}
				}
				if len(uniqueMembers) > 0 {
					ouAdd.Attribute("uniqueMember", uniqueMembers)
				}
			}
		}
		if group.Remark != "" {
			ouAdd.Attribute("description", []string{group.Remark})
		}
		if addErr := modifyConn.Add(ouAdd); addErr != nil {
			if ldapErr, ok := addErr.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				// Group already exists in LDAP
				Log.Infof("Group %s already exists in LDAP", group.GroupDN)
			} else {
				Log.Errorf("Failed to create group %s: %v", group.GroupDN, addErr)
			}
		}
	}
}

// extractParentDN extracts the parent DN from a full DN
func extractParentDN(dn string) string {
	parts := strings.SplitN(dn, ",", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// extractOUName extracts the OU name from a DN like "ou=docker,dc=cmplab,dc=com"
func extractOUName(dn string) string {
	parts := strings.Split(dn, ",")
	if len(parts) == 0 {
		return ""
	}
	first := parts[0]
	if strings.HasPrefix(first, "ou=") {
		return strings.TrimPrefix(first, "ou=")
	}
	return ""
}

// syncExistingUsersToPreConfig syncs existing users (except admin) to user pre-config table
func syncExistingUsersToPreConfig() {
	var users []model.User
	// Get all users except admin
	if err := DB.Where("username != ?", "admin").Find(&users).Error; err != nil {
		Log.Errorf("Failed to get users for pre-config sync: %v", err)
		return
	}

	for _, user := range users {
		// Check if pre-config already exists for this username
		var existingConfig model.UserPreConfig
		err := DB.Where("username = ?", user.Username).First(&existingConfig).Error
		if err == nil {
			// Pre-config already exists, skip
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			Log.Errorf("Failed to check pre-config for user %s: %v", user.Username, err)
			continue
		}

		// Create new pre-config from user data
		preConfig := model.UserPreConfig{
			Username:      user.Username,
			Mail:          user.Mail,
			UIDNumber:     user.UIDNumber,
			GIDNumber:     user.GIDNumber,
			DepartmentId:  user.DepartmentId,
			Departments:   user.Departments,
			Mobile:        user.Mobile,
			JobNumber:     user.JobNumber,
			Position:      user.Position,
			PostalAddress: user.PostalAddress,
			Introduction:  user.Introduction,
			HomeDirectory: user.HomeDirectory,
			LoginShell:    user.LoginShell,
			Nickname:      user.Nickname,
			GivenName:     user.GivenName,
			Creator:       "system",
			IsUsed:        true, // Mark as used since user already exists
		}

		if err := DB.Create(&preConfig).Error; err != nil {
			Log.Errorf("Failed to create pre-config for user %s: %v", user.Username, err)
		}
	}
}
