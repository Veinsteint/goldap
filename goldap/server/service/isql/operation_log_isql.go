package isql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"

	"gorm.io/gorm"
)

type OperationLogService struct{}

// SaveOperationLogChannel batches operation logs for efficiency
func (s OperationLogService) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	Logs := make([]model.OperationLog, 0)
	duration := 5 * time.Second
	timer := time.NewTimer(duration)
	defer timer.Stop()
	
	for {
		select {
		case log := <-olc:
			Logs = append(Logs, *log)
			if len(Logs) > 5 {
				common.DB.Create(&Logs)
				Logs = make([]model.OperationLog, 0)
				timer.Reset(duration)
			}
		case <-timer.C:
			if len(Logs) > 0 {
				common.DB.Create(&Logs)
				Logs = make([]model.OperationLog, 0)
			}
			timer.Reset(duration)
		}
	}
}

// List returns paginated operation logs
func (s OperationLogService) List(req *request.OperationLogListReq) ([]*model.OperationLog, error) {
	var list []*model.OperationLog
	db := common.DB.Model(&model.OperationLog{}).Order("id DESC")

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if ip := strings.TrimSpace(req.Ip); ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	if path := strings.TrimSpace(req.Path); path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	if method := strings.TrimSpace(req.Method); method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// Count returns total log count
func (s OperationLogService) Count() (count int64, err error) {
	err = common.DB.Model(&model.OperationLog{}).Count(&count).Error
	return count, err
}

// Find gets single log by filter
func (s OperationLogService) Find(filter map[string]any, data *model.OperationLog) error {
	return common.DB.Where(filter).First(&data).Error
}

// Exist checks if log exists
func (s OperationLogService) Exist(filter map[string]any) bool {
	var dataObj model.OperationLog
	err := common.DB.Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete removes logs by IDs
func (s OperationLogService) Delete(operationLogIds []uint) error {
	return common.DB.Where("id IN (?)", operationLogIds).Unscoped().Delete(&model.OperationLog{}).Error
}

// Clean removes all logs
func (s OperationLogService) Clean() error {
	return common.DB.Exec("DELETE FROM operation_logs").Error
}
