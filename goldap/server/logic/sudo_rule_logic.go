package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type SudoRuleLogic struct{}

// Add creates a sudo rule
func (l SudoRuleLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SudoRuleAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user"))
	}

	rule := &model.SudoRule{
		Name:        r.Name,
		Description: r.Description,
		User:        r.User,
		Host:        r.Host,
		Command:     r.Command,
		RunAsUser:   r.RunAsUser,
		RunAsGroup:  r.RunAsGroup,
		Options:     r.Options,
		Creator:     ctxUser.Username,
	}

	if err := isql.SudoRule.Add(rule); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to add sudo rule: %v", err))
	}

	if err := ildap.Sudo.Add(rule); err != nil {
		common.Log.Errorf("Failed to sync sudo rule %s to LDAP: %v", rule.Name, err)
	}

	return rule, nil
}

// Update modifies a sudo rule
func (l SudoRuleLogic) Update(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SudoRuleUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}

	oldRule, err := isql.SudoRule.GetByID(r.ID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("sudo rule not found"))
	}

	rule := &model.SudoRule{
		Model:       oldRule.Model,
		Name:        r.Name,
		Description: r.Description,
		User:        r.User,
		Host:        r.Host,
		Command:     r.Command,
		RunAsUser:   r.RunAsUser,
		RunAsGroup:  r.RunAsGroup,
		Options:     r.Options,
		Creator:     oldRule.Creator,
	}

	if err := isql.SudoRule.Update(rule); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to update sudo rule: %v", err))
	}

	if err := ildap.Sudo.Update(rule); err != nil {
		common.Log.Errorf("Failed to sync sudo rule %s to LDAP: %v", rule.Name, err)
	}

	return rule, nil
}

// Delete removes a sudo rule
func (l SudoRuleLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SudoRuleDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}

	rule, err := isql.SudoRule.GetByID(r.ID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("sudo rule not found"))
	}

	if err := ildap.Sudo.Delete(rule.Name); err != nil {
		common.Log.Errorf("Failed to delete sudo rule %s from LDAP: %v", rule.Name, err)
	}

	if err := isql.SudoRule.Delete(r.ID); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete sudo rule: %v", err))
	}

	return nil, nil
}

// List returns sudo rule list
func (l SudoRuleLogic) List(c *gin.Context, req any) (data any, rspError any) {
	rules, err := isql.SudoRule.List()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get sudo rule list: %v", err))
	}

	var result []map[string]interface{}
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":          rule.ID,
			"name":        rule.Name,
			"description": rule.Description,
			"user":        rule.User,
			"host":        rule.Host,
			"command":     rule.Command,
			"runAsUser":   rule.RunAsUser,
			"runAsGroup":  rule.RunAsGroup,
			"options":     rule.Options,
			"creator":     rule.Creator,
			"createdAt":   rule.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}
