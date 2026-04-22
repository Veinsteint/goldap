package controller

import (
	"fmt"
	"net/http"
	"regexp"

	"goldap-server/public/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zht "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Api                   = &ApiController{}
	Group                 = &GroupController{}
	Menu                  = &MenuController{}
	Role                  = &RoleController{}
	User                  = &UserController{}
	OperationLog          = &OperationLogController{}
	Base                  = &BaseController{}
	FieldRelation         = &FieldRelationController{}
	SSHKey                = &SSHKeyController{}
	IPGroup               = &IPGroupController{}
	IPGroupUserPermission = &IPGroupUserPermissionController{}
	GroupUserPermission   = &GroupUserPermissionController{}
	SudoRule              = &SudoRuleController{}
	PendingUser           = &PendingUserController{}
	UserPreConfig         = &UserPreConfigController{}
	SystemConfig          = &SystemConfigController{}

	validate = validator.New()
	trans    ut.Translator
)

func init() {
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	_ = zht.RegisterDefaultTranslations(validate, trans)
	_ = validate.RegisterValidation("checkMobile", checkMobile)
}

// checkMobile validates mobile phone number format
func checkMobile(fl validator.FieldLevel) bool {
	reg := `1\d{10}`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}

// Run binds request, validates, and executes handler
func Run(c *gin.Context, req any, fn func() (any, any)) {
	if req != nil {
		if err := c.Bind(req); err != nil {
			tools.Err(c, tools.NewValidatorError(err), nil)
			return
		}

		if err := validate.Struct(req); err != nil {
			// Check if it's ValidationErrors before type assertion
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrors {
					tools.Err(c, tools.NewValidatorError(fmt.Errorf("%s", e.Translate(trans))), nil)
					return
				}
			} else {
				// Handle other validation errors (e.g., InvalidValidationError)
				tools.Err(c, tools.NewValidatorError(err), nil)
				return
			}
		}
	}

	data, err := fn()
	if err != nil {
		tools.Err(c, tools.ReloadErr(err), data)
		return
	}
	tools.Success(c, data)
}

// Demo Health check endpoint
// @Summary Health Check
// @Tags Base
// @Produce json
// @Description Health check endpoint
// @Success 200 {object} response.ResponseBody
// @Router /base/ping [get]
func Demo(c *gin.Context) {
	CodeDebug()
	c.JSON(http.StatusOK, tools.H{"code": 200, "msg": "ok", "data": "pong"})
}

// CodeDebug placeholder for debugging
func CodeDebug() {}
