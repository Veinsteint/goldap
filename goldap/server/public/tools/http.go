package tools

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SystemErr    = 500
	MySqlErr     = 501
	LdapErr      = 505
	OperationErr = 506
	ValidatorErr = 412
)

type RspError struct {
	code int
	err  error
}

func (re *RspError) Error() string {
	return re.err.Error()
}

func (re *RspError) Code() int {
	return re.code
}

// NewRspError creates a response error
func NewRspError(code int, err error) *RspError {
	return &RspError{code: code, err: err}
}

// NewMySqlError creates MySQL error
func NewMySqlError(err error) *RspError {
	return NewRspError(MySqlErr, err)
}

// NewValidatorError creates validation error
func NewValidatorError(err error) *RspError {
	return NewRspError(ValidatorErr, err)
}

// NewLdapError creates LDAP error
func NewLdapError(err error) *RspError {
	return NewRspError(LdapErr, err)
}

// NewOperationError creates operation error
func NewOperationError(err error) *RspError {
	return NewRspError(OperationErr, err)
}

// ReloadErr converts any error to RspError
func ReloadErr(err any) *RspError {
	rspErr, ok := err.(*RspError)
	if !ok {
		rspError, ok := err.(error)
		if !ok {
			return &RspError{code: SystemErr, err: fmt.Errorf("unknown error")}
		}
		return &RspError{code: SystemErr, err: rspError}
	}
	return rspErr
}

// Success sends success response
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// Err sends error response
func Err(c *gin.Context, err *RspError, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": err.Code(),
		"msg":  err.Error(),
		"data": data,
	})
}

// Response sends custom response
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}
