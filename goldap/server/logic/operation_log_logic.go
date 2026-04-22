package logic

import (
	"fmt"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/model/response"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type OperationLogLogic struct{}

// List returns operation log list
func (l OperationLogLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.OperationLogListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	logs, err := isql.OperationLog.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get log list: %s", err.Error()))
	}

	rets := make([]model.OperationLog, 0)
	for _, log := range logs {
		rets = append(rets, *log)
	}

	count, err := isql.OperationLog.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get log count"))
	}

	return response.LogListRsp{
		Total: count,
		Logs:  rets,
	}, nil
}

// Delete removes log records
func (l OperationLogLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.OperationLogDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.OperationLogIds {
		filter := tools.H{"id": int(id)}
		if !isql.OperationLog.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("log not found"))
		}
	}

	if err := isql.OperationLog.Delete(r.OperationLogIds); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete log: %s", err.Error()))
	}

	return nil, nil
}

// Clean clears all logs
func (l OperationLogLogic) Clean(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.OperationLogListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if err := isql.OperationLog.Clean(); err != nil {
		return err, nil
	}

	return "Operation logs cleared", nil
}
