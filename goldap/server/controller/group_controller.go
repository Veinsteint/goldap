package controller

import (
	"goldap-server/logic"
	"goldap-server/model/request"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

// List returns group list
// @Summary Get Group List
// @Description Get group list
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/list [get]
// @Security ApiKeyAuth
func (m *GroupController) List(c *gin.Context) {
	req := new(request.GroupListReq)
	Run(c, req, func() (any, any) {
		return logic.Group.List(c, req)
	})
}

// UserInGroup returns users in group
// @Summary Get Users in Group
// @Description Get users that belong to the group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param groupId query int true "Group ID"
// @Param nickname query string false "Nickname filter"
// @Success 200 {object} response.ResponseBody
// @Router /group/useringroup [get]
// @Security ApiKeyAuth
func (m *GroupController) UserInGroup(c *gin.Context) {
	req := new(request.UserInGroupReq)
	Run(c, req, func() (any, any) {
		return logic.Group.UserInGroup(c, req)
	})
}

// UserNoInGroup returns users not in group
// @Summary Get Users Not in Group
// @Description Get users that do not belong to the group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param groupId query int true "Group ID"
// @Param nickname query string false "Nickname filter"
// @Success 200 {object} response.ResponseBody
// @Router /group/usernoingroup [get]
// @Security ApiKeyAuth
func (m *GroupController) UserNoInGroup(c *gin.Context) {
	req := new(request.UserNoInGroupReq)
	Run(c, req, func() (any, any) {
		return logic.Group.UserNoInGroup(c, req)
	})
}

// GetTree returns group tree
// @Summary Get Group Tree
// @Description Get group tree structure
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/tree [get]
// @Security ApiKeyAuth
func (m *GroupController) GetTree(c *gin.Context) {
	req := new(request.GroupListReq)
	Run(c, req, func() (any, any) {
		return logic.Group.GetTree(c, req)
	})
}

// Add creates new group
// @Summary Create Group
// @Description Create new group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param data body request.GroupAddReq true "Group data"
// @Success 200 {object} response.ResponseBody
// @Router /group/add [post]
// @Security ApiKeyAuth
func (m *GroupController) Add(c *gin.Context) {
	req := new(request.GroupAddReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Add(c, req)
	})
}

// Update updates group
// @Summary Update Group
// @Description Update existing group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param data body request.GroupUpdateReq true "Group data"
// @Success 200 {object} response.ResponseBody
// @Router /group/update [post]
// @Security ApiKeyAuth
func (m *GroupController) Update(c *gin.Context) {
	req := new(request.GroupUpdateReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Update(c, req)
	})
}

// Delete removes group
// @Summary Delete Group
// @Description Delete group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param data body request.GroupDeleteReq true "Group ID"
// @Success 200 {object} response.ResponseBody
// @Router /group/delete [post]
// @Security ApiKeyAuth
func (m *GroupController) Delete(c *gin.Context) {
	req := new(request.GroupDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.Group.Delete(c, req)
	})
}

// AddUser adds user to group
// @Summary Add User to Group
// @Description Add user to group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param data body request.GroupAddUserReq true "User and group data"
// @Success 200 {object} response.ResponseBody
// @Router /group/adduser [post]
// @Security ApiKeyAuth
func (m *GroupController) AddUser(c *gin.Context) {
	req := new(request.GroupAddUserReq)
	Run(c, req, func() (any, any) {
		return logic.Group.AddUser(c, req)
	})
}

// RemoveUser removes user from group
// @Summary Remove User from Group
// @Description Remove user from group
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Param data body request.GroupRemoveUserReq true "User and group data"
// @Success 200 {object} response.ResponseBody
// @Router /group/removeuser [post]
// @Security ApiKeyAuth
func (m *GroupController) RemoveUser(c *gin.Context) {
	req := new(request.GroupRemoveUserReq)
	Run(c, req, func() (any, any) {
		return logic.Group.RemoveUser(c, req)
	})
}

// SyncOpenLdapDepts syncs departments from OpenLDAP
// @Summary Sync OpenLDAP Departments
// @Description Sync department data from OpenLDAP
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncOpenLdapDepts [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncOpenLdapDepts(c *gin.Context) {
	req := new(request.SyncOpenLdapDeptsReq)
	Run(c, req, func() (any, any) {
		return logic.OpenLdap.SyncOpenLdapDepts(c, req)
	})
}

// SyncSqlGroups syncs groups from MySQL to LDAP
// @Summary Sync SQL Groups to LDAP
// @Description Sync group data from MySQL to LDAP
// @Tags Group Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /group/syncSqlGroups [post]
// @Security ApiKeyAuth
func (m *GroupController) SyncSqlGroups(c *gin.Context) {
	req := new(request.SyncSqlGrooupsReq)
	Run(c, req, func() (any, any) {
		return logic.Sql.SyncSqlGroups(c, req)
	})
}
