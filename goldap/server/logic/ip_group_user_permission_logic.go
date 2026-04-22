package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IPGroupUserPermissionLogic struct{}

// Add creates an IP group user permission
func (l IPGroupUserPermissionLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupUserPermissionAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if !isql.IPGroup.Exist(tools.H{"id": r.IPGroupID}) {
		return nil, tools.NewValidatorError(fmt.Errorf("IP group not found"))
	}

	if !isql.User.Exist(tools.H{"id": r.UserID}) {
		return nil, tools.NewValidatorError(fmt.Errorf("user not found"))
	}

	if isql.IPGroupUserPermission.Exist(tools.H{"ip_group_id": r.IPGroupID, "user_id": r.UserID}) {
		return nil, tools.NewValidatorError(fmt.Errorf("permission already exists"))
	}

	permission := &model.IPGroupUserPermission{
		IPGroupID:   r.IPGroupID,
		UserID:      r.UserID,
		AllowLogin:  r.AllowLogin,
		AllowSudo:   r.AllowSudo,
		AllowSSHKey: r.AllowSSHKey,
		SudoRules:   r.SudoRules,
	}

	if err := isql.IPGroupUserPermission.Add(permission); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to add permission: %v", err))
	}

	return permission, nil
}

// Update modifies an IP group user permission
func (l IPGroupUserPermissionLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupUserPermissionUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}

	permission := &model.IPGroupUserPermission{
		Model:       gorm.Model{ID: r.ID},
		AllowLogin:  r.AllowLogin,
		AllowSudo:   r.AllowSudo,
		AllowSSHKey: r.AllowSSHKey,
		SudoRules:   r.SudoRules,
	}

	if err := isql.IPGroupUserPermission.Update(permission); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update permission: %v", err))
	}

	return permission, nil
}

// Delete removes an IP group user permission
func (l IPGroupUserPermissionLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupUserPermissionDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if err := isql.IPGroupUserPermission.Delete(r.ID); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete permission: %v", err))
	}

	return nil, nil
}

// List returns IP group user permission list
func (l IPGroupUserPermissionLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupUserPermissionListReq)
	if !ok {
		return nil, ReqAssertErr
	}

	var permissions []*model.IPGroupUserPermission
	var err error

	if r.IPGroupID > 0 {
		permissions, err = isql.IPGroupUserPermission.GetByIPGroupID(r.IPGroupID)
	} else if r.UserID > 0 {
		permissions, err = isql.IPGroupUserPermission.GetByUserID(r.UserID)
	} else {
		return nil, tools.NewValidatorError(fmt.Errorf("IP group ID or user ID required"))
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
			"ipGroupId":   perm.IPGroupID,
			"ipGroupName": perm.IPGroup.Name,
			"userId":      perm.UserID,
			"username":    perm.User.Username,
			"allowLogin":  perm.AllowLogin,
			"allowSudo":   perm.AllowSudo,
			"allowSSHKey": perm.AllowSSHKey,
			"sudoRules":   sudoRules,
			"createdAt":   perm.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// CheckPermission checks user permission for given IP
func (l IPGroupUserPermissionLogic) CheckPermission(userID uint, clientIP string) (*model.IPGroupUserPermission, error) {
	permissions, err := isql.IPGroupUserPermission.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, perm := range permissions {
		matched, err := tools.IPInRanges(clientIP, perm.IPGroup.IPRanges)
		if err != nil {
			common.Log.Warnf("Failed to check IP range: %v", err)
			continue
		}
		if matched {
			return perm, nil
		}
	}

	return nil, fmt.Errorf("no matching IP group permission found")
}
