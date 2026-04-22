package isql

var (
	User                    = &UserService{}
	Group                   = &GroupService{}
	Api                     = &ApiService{}
	Menu                    = &MenuService{}
	Role                    = &RoleService{}
	OperationLog            = &OperationLogService{}
	FieldRelation           = &FieldRelationService{}
	SSHKey                  = &SSHKeyService{}
	IPGroup                 = &IPGroupService{}
	IPGroupUserPermission   = &IPGroupUserPermissionService{}
	GroupUserPermission     = &GroupUserPermissionService{}
	SudoRule                = &SudoRuleService{}
	PendingUser             = &PendingUserService{}
	UserPreConfig           = &UserPreConfigService{}
	SystemConfig            = &SystemConfigService{}
)
