package isql

import (
	"goldap-server/model"
	"goldap-server/public/common"

	"gorm.io/gorm"
)

type SSHKeyService struct{}

// Add creates SSH key
func (s SSHKeyService) Add(sshKey *model.SSHKey) error {
	return common.DB.Create(sshKey).Error
}

// GetByUserID gets all SSH keys for user
func (s SSHKeyService) GetByUserID(userID uint) ([]*model.SSHKey, error) {
	var sshKeys []*model.SSHKey
	err := common.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&sshKeys).Error
	return sshKeys, err
}

// GetByID gets SSH key by ID
func (s SSHKeyService) GetByID(id uint) (*model.SSHKey, error) {
	var sshKey model.SSHKey
	err := common.DB.First(&sshKey, id).Error
	if err != nil {
		return nil, err
	}
	return &sshKey, nil
}

// Delete removes SSH key (user can only delete own keys)
func (s SSHKeyService) Delete(id uint, userID uint) error {
	result := common.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.SSHKey{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllByUserID gets all SSH keys for authorized_keys
func (s SSHKeyService) GetAllByUserID(userID uint) ([]*model.SSHKey, error) {
	var sshKeys []*model.SSHKey
	err := common.DB.Where("user_id = ?", userID).Find(&sshKeys).Error
	return sshKeys, err
}
