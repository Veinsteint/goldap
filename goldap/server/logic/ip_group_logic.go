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

type IPGroupLogic struct{}

// Add creates an IP group
func (l IPGroupLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	for _, ipRange := range r.IPRanges {
		if err := tools.ValidateIPRange(ipRange); err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("invalid IP range: %v", err))
		}
	}

	ipRangesJSON, err := json.Marshal(r.IPRanges)
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to serialize IP ranges: %v", err))
	}

	ipGroup := &model.IPGroup{
		Name:        r.Name,
		Description: r.Description,
		IPRanges:    string(ipRangesJSON),
		Creator:     ctxUser.Username,
	}

	if err := isql.IPGroup.Add(ipGroup); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to add IP group: %v", err))
	}

	return ipGroup, nil
}

// Update modifies an IP group
func (l IPGroupLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}

	for _, ipRange := range r.IPRanges {
		if err := tools.ValidateIPRange(ipRange); err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("invalid IP range: %v", err))
		}
	}

	ipRangesJSON, err := json.Marshal(r.IPRanges)
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("failed to serialize IP ranges: %v", err))
	}

	ipGroup := &model.IPGroup{
		Model:       gorm.Model{ID: r.ID},
		Name:        r.Name,
		Description: r.Description,
		IPRanges:    string(ipRangesJSON),
	}

	if err := isql.IPGroup.Update(ipGroup); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update IP group: %v", err))
	}

	return ipGroup, nil
}

// Delete removes an IP group
func (l IPGroupLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.IPGroupDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}

	if err := isql.IPGroup.Delete(r.ID); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete IP group: %v", err))
	}

	return nil, nil
}

// List returns IP group list
func (l IPGroupLogic) List(c *gin.Context, req any) (data any, rspError any) {
	ipGroups, err := isql.IPGroup.List()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get IP group list: %v", err))
	}

	var result []map[string]interface{}
	for _, ipGroup := range ipGroups {
		var ipRanges []string
		if err := json.Unmarshal([]byte(ipGroup.IPRanges), &ipRanges); err != nil {
			common.Log.Warnf("Failed to parse IP ranges for group %s: %v", ipGroup.Name, err)
			ipRanges = []string{}
		}

		result = append(result, map[string]interface{}{
			"id":          ipGroup.ID,
			"name":        ipGroup.Name,
			"description": ipGroup.Description,
			"ipRanges":    ipRanges,
			"creator":     ipGroup.Creator,
			"createdAt":   ipGroup.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}
