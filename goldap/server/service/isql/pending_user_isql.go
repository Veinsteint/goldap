package isql

import (
	"fmt"
	"strings"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
)

type PendingUserService struct{}

// Add creates pending user
func (s PendingUserService) Add(pendingUser *model.PendingUser) error {
	return common.DB.Create(pendingUser).Error
}

// List returns paginated pending user list
func (s PendingUserService) List(req *request.PendingUserListReq) ([]*model.PendingUser, error) {
	var list []*model.PendingUser
	db := common.DB.Model(&model.PendingUser{}).Order("id DESC")

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if nickname := strings.TrimSpace(req.Nickname); nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if mail := strings.TrimSpace(req.Mail); mail != "" {
		db = db.Where("mail LIKE ?", fmt.Sprintf("%%%s%%", mail))
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	} else {
		db = db.Where("status = ?", 0) // Default: pending only
	}

	if req.PageNum > 0 && req.PageSize > 0 {
		offset := (req.PageNum - 1) * req.PageSize
		db = db.Offset(offset).Limit(req.PageSize)
	}

	err := db.Find(&list).Error
	return list, err
}

// ListCount returns total pending user count
func (s PendingUserService) ListCount(req *request.PendingUserListReq) (int64, error) {
	var count int64
	db := common.DB.Model(&model.PendingUser{})

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if nickname := strings.TrimSpace(req.Nickname); nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if mail := strings.TrimSpace(req.Mail); mail != "" {
		db = db.Where("mail LIKE ?", fmt.Sprintf("%%%s%%", mail))
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	} else {
		db = db.Where("status = ?", 0)
	}

	err := db.Count(&count).Error
	return count, err
}

// Find gets pending user by ID
func (s PendingUserService) Find(id uint) (*model.PendingUser, error) {
	var pendingUser model.PendingUser
	err := common.DB.Where("id = ?", id).First(&pendingUser).Error
	return &pendingUser, err
}

// Update modifies pending user
func (s PendingUserService) Update(pendingUser *model.PendingUser) error {
	return common.DB.Model(pendingUser).Updates(pendingUser).Error
}

// Delete removes pending user
func (s PendingUserService) Delete(id uint) error {
	return common.DB.Delete(&model.PendingUser{}, id).Error
}

// Exist checks if pending user exists
func (s PendingUserService) Exist(filter map[string]any) bool {
	var dataObj model.PendingUser
	err := common.DB.Where(filter).First(&dataObj).Error
	return err == nil
}
