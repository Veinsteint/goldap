package middleware

import (
	"strings"
	"sync"

	"goldap-server/config"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

var checkLock sync.Mutex

// CasbinMiddleware RBAC authorization middleware
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := isql.User.GetCurrentLoginUser(c)
		if err != nil {
			tools.Response(c, 401, 401, nil, "User not logged in")
			c.Abort()
			return
		}
		if user.Status != 1 {
			tools.Response(c, 401, 401, nil, "User disabled")
			c.Abort()
			return
		}

		var subs []string
		for _, role := range user.Roles {
			if role.Status == 1 {
				subs = append(subs, role.Keyword)
			}
		}

		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		act := c.Request.Method
		
		if !check(subs, obj, act) {
			tools.Response(c, 401, 401, nil, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

func check(subs []string, obj string, act string) bool {
	checkLock.Lock()
	defer checkLock.Unlock()
	for _, sub := range subs {
		if pass, _ := common.CasbinEnforcer.Enforce(sub, obj, act); pass {
			return true
		}
	}
	return false
}
