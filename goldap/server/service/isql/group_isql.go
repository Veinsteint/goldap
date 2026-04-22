package isql

import (
	"errors"
	"fmt"
	"strings"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"

	ldap "github.com/go-ldap/ldap/v3"
	"gorm.io/gorm"
)

type GroupService struct{}

// List returns paginated group list
func (s GroupService) List(req *request.GroupListReq) ([]*model.Group, error) {
	var list []*model.Group
	db := common.DB.Model(&model.Group{}).Order("created_at DESC")

	if groupName := strings.TrimSpace(req.GroupName); groupName != "" {
		db = db.Where("group_name LIKE ?", fmt.Sprintf("%%%s%%", groupName))
	}
	if groupRemark := strings.TrimSpace(req.Remark); groupRemark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", groupRemark))
	}
	if req.SyncState != 0 {
		db = db.Where("sync_state = ?", req.SyncState)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("Users").Find(&list).Error
	return list, err
}

// ListTree returns groups for tree display
func (s GroupService) ListTree(req *request.GroupListReq) ([]*model.Group, error) {
	var list []*model.Group
	db := common.DB.Model(&model.Group{}).Order("created_at DESC")

	if groupName := strings.TrimSpace(req.GroupName); groupName != "" {
		db = db.Where("group_name LIKE ?", fmt.Sprintf("%%%s%%", groupName))
	}
	if groupRemark := strings.TrimSpace(req.Remark); groupRemark != "" {
		db = db.Where("remark LIKE ?", fmt.Sprintf("%%%s%%", groupRemark))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// ListAll returns all groups
func (s GroupService) ListAll() (list []*model.Group, err error) {
	err = common.DB.Model(&model.Group{}).Order("created_at DESC").Find(&list).Error
	return list, err
}

// GenGroupTree generates group tree structure
func GenGroupTree(parentId uint, groups []*model.Group) []*model.Group {
	tree := make([]*model.Group, 0)
	for _, g := range groups {
		if g.ParentId == parentId {
			g.Children = GenGroupTree(g.ID, groups)
			tree = append(tree, g)
		}
	}
	return tree
}

// Count returns total group count
func (s GroupService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Group{}).Count(&count).Error
	return count, err
}

// Add creates group and syncs to LDAP
func (s GroupService) Add(data *model.Group) error {
	err := common.DB.Create(data).Error
	if err != nil {
		return err
	}

	if data.GroupDN != "" {
		err = ildap.Group.Add(data)
		if err != nil {
			// For groupOfUniqueNames without users
			if strings.Contains(data.GroupDN, "ou=sudoers,") == false && data.GroupType == "cn" {
				if strings.Contains(err.Error(), "uniqueMember") || strings.Contains(err.Error(), "Object Class Violation") {
					common.Log.Infof("Group [%s] created in MySQL. LDAP group will be created when first user is added.", data.GroupDN)
					_ = s.ChangeSyncState(int(data.ID), 2)
					return nil // Don't treat as error
				}
			}
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				// Group already exists
				_ = s.ChangeSyncState(int(data.ID), 1)
			} else {
				common.Log.Warnf("Failed to sync group [%s] to LDAP: %v", data.GroupDN, err)
				_ = s.ChangeSyncState(int(data.ID), 2)
			}
		} else {
			_ = s.ChangeSyncState(int(data.ID), 1)
		}
	}

	return nil
}

// Update modifies group and syncs to LDAP
func (s GroupService) Update(dataObj *model.Group) error {
	var oldGroup model.Group
	err := common.DB.Where("id = ?", dataObj.ID).First(&oldGroup).Error
	if err != nil {
		return err
	}

	err = common.DB.Model(dataObj).Where("id = ?", dataObj.ID).
		Select("group_name", "remark", "creator", "group_dn", "ip_ranges", "updated_at").
		Updates(dataObj).Error
	if err != nil {
		return err
	}

	if oldGroup.GroupDN != "" && dataObj.GroupDN != "" {
		err = ildap.Group.Update(&oldGroup, dataObj)
		if err != nil {
			common.Log.Warnf("Failed to sync group [%s] to LDAP: %v", dataObj.GroupDN, err)
			_ = s.ChangeSyncState(int(dataObj.ID), 2)
		} else {
			_ = s.ChangeSyncState(int(dataObj.ID), 1)
		}
	}

	return nil
}

// ChangeSyncState updates group sync state
func (s GroupService) ChangeSyncState(id, status int) error {
	return common.DB.Model(&model.Group{}).Where("id = ?", id).Update("sync_state", status).Error
}

// Find gets single group by filter
func (s GroupService) Find(filter map[string]any, data *model.Group, args ...any) error {
	return common.DB.Where(filter, args).Preload("Users").First(&data).Error
}

// Exist checks if group exists
func (s GroupService) Exist(filter map[string]any) bool {
	var dataObj model.Group
	err := common.DB.Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete removes groups from both MySQL and LDAP
func (s GroupService) Delete(groups []*model.Group) error {
	for _, group := range groups {
		if group.GroupDN != "" {
			if err := ildap.Group.Delete(group.GroupDN); err != nil {
				common.Log.Warnf("Failed to delete LDAP group [%s]: %v", group.GroupDN, err)
			}
		}
	}
	return common.DB.Select("Users").Unscoped().Delete(&groups).Error
}

// GetGroupByIds gets groups by IDs
func (s GroupService) GetGroupByIds(ids []uint) (datas []*model.Group, err error) {
	err = common.DB.Where("id IN (?)", ids).Preload("Users").Find(&datas).Error
	return datas, err
}

// AddUserToGroup adds users to group
func (s GroupService) AddUserToGroup(group *model.Group, users []model.User) error {
	return common.DB.Model(&group).Association("Users").Append(users)
}

// RemoveUserFromGroup removes users from group
func (s GroupService) RemoveUserFromGroup(group *model.Group, users []model.User) error {
	return common.DB.Model(&group).Association("Users").Delete(users)
}

// DeptIdsToGroupIds converts department IDs to group IDs
func (s GroupService) DeptIdsToGroupIds(ids []string) (groupIds []uint, err error) {
	var tempGroups []model.Group
	err = common.DB.Model(&model.Group{}).Where("source_dept_id IN (?)", ids).Find(&tempGroups).Error
	if err != nil {
		return nil, err
	}
	var tempGroupIds []uint
	for _, g := range tempGroups {
		tempGroupIds = append(tempGroupIds, g.ID)
	}
	return tempGroupIds, nil
}

// UserInGroup checks if user is in group
func (s GroupService) UserInGroup(groupID, userID uint) bool {
	var count int64
	err := common.DB.Table("group_users").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return err == nil && count > 0
}
