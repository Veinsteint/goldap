package controller

import (
	"fmt"

	"goldap-server/logic"
	"goldap-server/model/request"
	"goldap-server/public/tools"

	"github.com/gin-gonic/gin"
)

type SSHKeyController struct{}

// GetSSHKeys returns user's SSH public keys
// @Summary Get SSH Public Keys
// @Description Get current user's SSH public key list
// @Tags SSH Key Management
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /user/ssh-keys [get]
// @Security ApiKeyAuth
func (sc *SSHKeyController) GetSSHKeys(c *gin.Context) {
	req := new(request.SSHKeyListReq)
	Run(c, req, func() (any, any) {
		return logic.SSHKey.List(c, req)
	})
}

// AddSSHKey adds SSH public key
// @Summary Add SSH Public Key
// @Description Add SSH public key for current user
// @Tags SSH Key Management
// @Accept application/json
// @Produce application/json
// @Param data body request.SSHKeyAddReq true "SSH key data"
// @Success 200 {object} response.ResponseBody
// @Router /user/ssh-keys [post]
// @Security ApiKeyAuth
func (sc *SSHKeyController) AddSSHKey(c *gin.Context) {
	req := new(request.SSHKeyAddReq)
	Run(c, req, func() (any, any) {
		return logic.SSHKey.Add(c, req)
	})
}

// DeleteSSHKey removes SSH public key
// @Summary Delete SSH Public Key
// @Description Delete SSH public key by ID
// @Tags SSH Key Management
// @Accept application/json
// @Produce application/json
// @Param id path int true "SSH Key ID"
// @Success 200 {object} response.ResponseBody
// @Router /user/ssh-keys/{id} [delete]
// @Security ApiKeyAuth
func (sc *SSHKeyController) DeleteSSHKey(c *gin.Context) {
	req := new(request.SSHKeyDeleteReq)

	if err := c.ShouldBindUri(req); err != nil {
		_ = c.ShouldBindJSON(req)
	}

	if req.ID == 0 {
		tools.Err(c, tools.NewValidatorError(fmt.Errorf("SSH key ID is required")), nil)
		return
	}

	Run(c, req, func() (any, any) {
		return logic.SSHKey.Delete(c, req)
	})
}
