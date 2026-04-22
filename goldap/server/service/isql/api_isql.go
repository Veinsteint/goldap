package isql

import (
	"errors"
	"fmt"
	"strings"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"

	"gorm.io/gorm"
)

type ApiService struct{}

// List returns paginated API list
func (s ApiService) List(req *request.ApiListReq) ([]*model.Api, error) {
	var list []*model.Api
	db := common.DB.Model(&model.Api{}).Order("created_at DESC")

	if method := strings.TrimSpace(req.Method); method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	if path := strings.TrimSpace(req.Path); path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	if category := strings.TrimSpace(req.Category); category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	if creator := strings.TrimSpace(req.Creator); creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// ListAll returns all APIs
func (s ApiService) ListAll() (list []*model.Api, err error) {
	err = common.DB.Model(&model.Api{}).Order("created_at DESC").Find(&list).Error
	return list, err
}

// Count returns total API count
func (s ApiService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Api{}).Count(&count).Error
	return count, err
}

// Add creates a new API
func (s ApiService) Add(api *model.Api) error {
	return common.DB.Create(api).Error
}

// Update modifies API and updates Casbin policies
func (s ApiService) Update(api *model.Api) error {
	var oldApi model.Api
	err := common.DB.First(&oldApi, api.ID).Error
	if err != nil {
		return errors.New("API not found")
	}
	
	err = common.DB.Model(api).Where("id = ?", api.ID).Updates(api).Error
	if err != nil {
		return err
	}
	
	// Update Casbin policy if path or method changed
	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies := common.CasbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		if len(policies) > 0 {
			isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("failed to update API permission")
			}
			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}
			isAdded, _ := common.CasbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("failed to update API permission")
			}
			err := common.CasbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("API updated but policy reload failed")
			}
			return err
		}
	}
	return err
}

// Find gets single API by filter
func (s ApiService) Find(filter map[string]any, data *model.Api) error {
	err := common.DB.Where(filter).First(data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return err
}

// Exist checks if API exists
func (s ApiService) Exist(filter map[string]any) bool {
	var count int64
	common.DB.Model(&model.Api{}).Where(filter).Count(&count)
	return count > 0
}

// Delete removes APIs and Casbin policies
func (s ApiService) Delete(ids []uint) error {
	var apis []model.Api
	for _, id := range ids {
		api := new(model.Api)
		err := s.Find(tools.H{"id": id}, api)
		if err != nil {
			return fmt.Errorf("failed to get API: %v", err)
		}
		apis = append(apis, *api)
	}

	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.Api{}).Error
	if err == nil {
		for _, api := range apis {
			policies := common.CasbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("failed to delete API permission")
				}
			}
		}
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("API deleted but policy reload failed")
		}
		return err
	}
	return err
}

// GetApisById gets APIs by IDs
func (s ApiService) GetApisById(apiIds []uint) ([]*model.Api, error) {
	var apis []*model.Api
	err := common.DB.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}
