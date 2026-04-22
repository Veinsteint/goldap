package middleware

import (
	"fmt"
	"strings"
	"time"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

var OperationLogChan = make(chan *model.OperationLog, 30)

// OperationLogMiddleware logs API operations
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		timeCost := endTime.Sub(startTime).Milliseconds()

		var username string
		ctxUser, _ := c.Get("user")
		if user, ok := ctxUser.(model.User); ok {
			username = user.Username
		} else {
			username = "anonymous"
		}

		path := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		method := c.Request.Method

		api := new(model.Api)
		if path != "" && isql.Api.Exist(tools.H{"path": path, "method": method}) {
			_ = isql.Api.Find(tools.H{"path": path, "method": method}, api)
		}

		OperationLogChan <- &model.OperationLog{
			Username:   username,
			Ip:         c.ClientIP(),
			IpLocation: "",
			Method:     method,
			Path:       path,
			Remark:     api.Remark,
			Status:     c.Writer.Status(),
			StartTime:  fmt.Sprintf("%v", startTime),
			TimeCost:   timeCost,
			UserAgent:  c.Request.UserAgent(),
		}
	}
}
