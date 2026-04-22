package middleware

import (
	"fmt"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/isql"

	"time"

	"goldap-server/model/request"
	"goldap-server/model/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitAuth initializes JWT authentication middleware
func InitAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.Jwt.Realm,
		Key:             []byte(config.Conf.Jwt.Key),
		Timeout:         time.Hour * time.Duration(config.Conf.Jwt.Timeout),
		MaxRefresh:      time.Hour * time.Duration(config.Conf.Jwt.MaxRefresh),
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   login,
		Authorizator:    authorizator,
		Unauthorized:    unauthorized,
		LoginResponse:   loginResponse,
		LogoutResponse:  logoutResponse,
		RefreshResponse: refreshResponse,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
	return authMiddleware, err
}

func payloadFunc(data any) jwt.MapClaims {
	if v, ok := data.(tools.H); ok {
		var user model.User
		tools.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) any {
	claims := jwt.ExtractClaims(c)
	return tools.H{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

func login(c *gin.Context) (any, error) {
	var req request.RegisterAndLoginReq
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}

	decodeData, err := tools.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username: req.Username,
		Password: string(decodeData),
	}

	user, err := isql.User.Login(u)
	if err != nil {
		return nil, err
	}
	return tools.H{
		"user": tools.Struct2Json(user),
	}, nil
}

func authorizator(data any, c *gin.Context) bool {
	if v, ok := data.(tools.H); ok {
		userStr := v["user"].(string)
		var user model.User
		tools.Json2Struct(userStr, &user)
		c.Set("user", user)
		return true
	}
	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	common.Log.Debugf("JWT auth failed, code: %d, message: %s", code, message)
	response.Response(c, code, code, nil, fmt.Sprintf("JWT auth failed: %d, %s", code, message))
}

func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.Response(c, code, code,
		gin.H{
			"token":   token,
			"expires": expires.Format("2006-01-02 15:04:05"),
		},
		"Login success")
}

func logoutResponse(c *gin.Context, code int) {
	response.Success(c, nil, "Logout success")
}

func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	response.Response(c, code, code,
		gin.H{
			"token":   token,
			"expires": expires,
		},
		"Token refreshed")
}
