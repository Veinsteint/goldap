package logic

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	"github.com/gin-gonic/gin"
)

type SSHKeyLogic struct{}

// Add creates an SSH public key
func (l SSHKeyLogic) Add(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SSHKeyAddReq)
	if !ok {
		return nil, ReqAssertErr
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user: %v", err))
	}

	// Check IP-based permissions
	clientIP := tools.GetClientIP(c.ClientIP(), c.GetHeader("X-Forwarded-For"), c.GetHeader("X-Real-IP"))
	if clientIP != "" {
		permission, err := IPGroupUserPermission.CheckPermission(ctxUser.ID, clientIP)
		if err == nil && !permission.AllowSSHKey {
			return nil, tools.NewValidatorError(fmt.Errorf("SSH key not allowed from IP %s", clientIP))
		}
	}

	key := strings.TrimSpace(r.Key)
	if !l.validateSSHKey(key) {
		return nil, tools.NewValidatorError(fmt.Errorf("invalid SSH key format"))
	}

	keyType := l.extractKeyType(key)

	// Check for duplicates
	existingKeys, err := isql.SSHKey.GetByUserID(ctxUser.ID)
	if err == nil {
		for _, existingKey := range existingKeys {
			if strings.TrimSpace(existingKey.Key) == key {
				return nil, tools.NewValidatorError(fmt.Errorf("SSH key already exists"))
			}
		}
	}

	sshKey := &model.SSHKey{
		UserID:  ctxUser.ID,
		Title:   strings.TrimSpace(r.Title),
		Key:     key,
		KeyType: keyType,
	}

	if err := isql.SSHKey.Add(sshKey); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to add SSH key: %v", err))
	}

	if err := l.updateAuthorizedKeys(ctxUser.Username, ctxUser.ID); err != nil {
		common.Log.Warnf("Failed to update authorized_keys: %v", err)
	}

	if err := l.syncSSHKeysToLDAP(ctxUser); err != nil {
		common.Log.Warnf("Failed to sync SSH keys to LDAP: %v", err)
	}

	return sshKey, nil
}

// List returns user's SSH keys
func (l SSHKeyLogic) List(c *gin.Context, req any) (data any, rspError any) {
	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user: %v", err))
	}

	sshKeys, err := isql.SSHKey.GetByUserID(ctxUser.ID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get SSH keys: %v", err))
	}

	var result []map[string]interface{}
	for _, key := range sshKeys {
		result = append(result, map[string]interface{}{
			"id":        key.ID,
			"title":     key.Title,
			"key":       key.Key,
			"keyType":   key.KeyType,
			"createdAt": key.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// Delete removes an SSH key
func (l SSHKeyLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SSHKeyDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}

	ctxUser, err := isql.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get current user: %v", err))
	}

	sshKey, err := isql.SSHKey.GetByID(r.ID)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("SSH key not found"))
	}

	if sshKey.UserID != ctxUser.ID {
		return nil, tools.NewValidatorError(fmt.Errorf("permission denied"))
	}

	if err := isql.SSHKey.Delete(r.ID, ctxUser.ID); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to delete SSH key: %v", err))
	}

	if err := l.updateAuthorizedKeys(ctxUser.Username, ctxUser.ID); err != nil {
		common.Log.Warnf("Failed to update authorized_keys: %v", err)
	}

	if err := l.syncSSHKeysToLDAP(ctxUser); err != nil {
		common.Log.Warnf("Failed to sync SSH keys to LDAP: %v", err)
	}

	return nil, nil
}

// syncSSHKeysToLDAP syncs SSH keys to LDAP
func (l SSHKeyLogic) syncSSHKeysToLDAP(user model.User) error {
	if user.UserDN == "" {
		return fmt.Errorf("user DN is empty")
	}

	sshKeys, err := isql.SSHKey.GetAllByUserID(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get SSH keys: %v", err)
	}

	var sshKeyStrings []string
	for _, sshKey := range sshKeys {
		key := strings.TrimSpace(sshKey.Key)
		if key != "" {
			sshKeyStrings = append(sshKeyStrings, key)
		}
	}

	return ildap.User.UpdateSSHKeys(user.UserDN, sshKeyStrings)
}

// validateSSHKey validates SSH key format
func (l SSHKeyLogic) validateSSHKey(key string) bool {
	key = strings.TrimSpace(key)
	validPrefixes := []string{
		"ssh-rsa", "ssh-ed25519", "ecdsa-sha2-nistp256",
		"ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521", "ssh-dss",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(key, prefix) && len(key) >= 50 {
			return true
		}
	}
	return false
}

// extractKeyType extracts SSH key type
func (l SSHKeyLogic) extractKeyType(key string) string {
	parts := strings.Fields(strings.TrimSpace(key))
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// updateAuthorizedKeys updates user's authorized_keys file
func (l SSHKeyLogic) updateAuthorizedKeys(username string, userID uint) error {
	sshKeys, err := isql.SSHKey.GetAllByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get SSH keys: %v", err)
	}

	sysUser, err := user.Lookup(username)
	if err != nil {
		common.Log.Warnf("System user %s not found, skipping authorized_keys update", username)
		return nil
	}

	sshDir := filepath.Join(sysUser.HomeDir, ".ssh")
	authorizedKeysPath := filepath.Join(sshDir, "authorized_keys")

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}

	var content strings.Builder
	for _, sshKey := range sshKeys {
		key := strings.TrimSpace(sshKey.Key)
		if key != "" {
			content.WriteString(fmt.Sprintf("# %s (ID: %d)\n", sshKey.Title, sshKey.ID))
			content.WriteString(key)
			content.WriteString("\n\n")
		}
	}

	if err := os.WriteFile(authorizedKeysPath, []byte(content.String()), 0600); err != nil {
		return fmt.Errorf("failed to write authorized_keys: %v", err)
	}

	_ = os.Chmod(authorizedKeysPath, 0600)
	_ = os.Chmod(sshDir, 0700)

	return nil
}

