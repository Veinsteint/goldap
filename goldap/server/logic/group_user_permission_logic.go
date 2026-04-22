package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupUserPermissionLogic struct{}

// Add creates a group user permission
func (l GroupUserPermissionLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUserPermissionAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if !isql.Group.Exist(tools.H{"id": r.GroupID}) {
		return nil, tools.NewValidatorError(fmt.Errorf("group not found"))
	}

	if !isql.User.Exist(tools.H{"id": r.UserID}) {
		return nil, tools.NewValidatorError(fmt.Errorf("user not found"))
	}

	if !isql.Group.UserInGroup(r.GroupID, r.UserID) {
		return nil, tools.NewValidatorError(fmt.Errorf("user not in group"))
	}

	if isql.GroupUserPermission.Exist(tools.H{"group_id": r.GroupID, "user_id": r.UserID}) {
		return nil, tools.NewValidatorError(fmt.Errorf("permission already exists"))
	}

	permission := &model.GroupUserPermission{
		GroupID:     r.GroupID,
		UserID:      r.UserID,
		AllowSudo:   r.AllowSudo,
		AllowSSHKey: r.AllowSSHKey,
		SudoRules:   r.SudoRules,
	}

	if err := isql.GroupUserPermission.Add(permission); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to add permission: %v", err))
	}

	return permission, nil
}

// Update modifies a group user permission
func (l GroupUserPermissionLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUserPermissionUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}

	permission := &model.GroupUserPermission{
		Model:       gorm.Model{ID: r.ID},
		AllowSudo:   r.AllowSudo,
		AllowSSHKey: r.AllowSSHKey,
		SudoRules:   r.SudoRules,
	}

	if err := isql.GroupUserPermission.Update(permission); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update permission: %v", err))
	}

	return permission, nil
}

// Delete removes a group user permission
func (l GroupUserPermissionLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUserPermissionDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if err := isql.GroupUserPermission.Delete(r.ID); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete permission: %v", err))
	}

	return nil, nil
}

// List returns group user permission list
func (l GroupUserPermissionLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.GroupUserPermissionListReq)
	if !ok {
		return nil, ReqAssertErr
	}

	var permissions []*model.GroupUserPermission
	var err error

	if r.GroupID > 0 {
		permissions, err = isql.GroupUserPermission.GetByGroupID(r.GroupID)
	} else if r.UserID > 0 {
		permissions, err = isql.GroupUserPermission.GetByUserID(r.UserID)
	} else {
		return nil, tools.NewValidatorError(fmt.Errorf("group ID or user ID required"))
	}

	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get permissions: %v", err))
	}

	var result []map[string]interface{}
	for _, perm := range permissions {
		var sudoRules interface{}
		if perm.SudoRules != "" {
			_ = json.Unmarshal([]byte(perm.SudoRules), &sudoRules)
		}

		result = append(result, map[string]interface{}{
			"id":          perm.ID,
			"groupId":     perm.GroupID,
			"groupName":   perm.Group.GroupName,
			"userId":      perm.UserID,
			"username":    perm.User.Username,
			"nickname":    perm.User.Nickname,
			"allowSudo":   perm.AllowSudo,
			"allowSSHKey": perm.AllowSSHKey,
			"sudoRules":   sudoRules,
			"createdAt":   perm.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}
