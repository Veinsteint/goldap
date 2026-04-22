package isql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type UserService struct{}

var userInfoCache = cache.New(24*time.Hour, 48*time.Hour)

// Add creates a new user with auto-assigned Unix attributes
func (s UserService) Add(user *model.User) error {
	user.Password = tools.NewGenPasswd(user.Password)
	
	if user.UIDNumber == 0 {
		uid, err := s.GetNextUIDNumber()
		if err != nil {
			return fmt.Errorf("failed to get next UID: %v", err)
		}
		user.UIDNumber = uid
	}
	
	if user.GIDNumber == 0 && user.UIDNumber > 0 {
		user.GIDNumber = user.UIDNumber
	}
	
	if user.HomeDirectory == "" {
		user.HomeDirectory = fmt.Sprintf("/home/%s", user.Username)
	}
	
	if user.LoginShell == "" {
		user.LoginShell = "/bin/bash"
	}
	
	if user.Gecos == "" {
		if user.Nickname != "" {
			user.Gecos = user.Nickname
		} else {
			user.Gecos = user.Username
		}
	}
	
	if user.SyncState == 0 {
		user.SyncState = 1
	}
	
	err := common.DB.Create(user).Error
	if err != nil {
		return err
	}
	
	// Sync to LDAP
	if user.SyncState == 1 {
		var syncUser model.User
		if err := common.DB.Preload("Roles").First(&syncUser, user.ID).Error; err != nil {
			common.Log.Warnf("Failed to sync user to LDAP: query error: %v", err)
			common.DB.Model(user).Update("sync_state", 2)
			return nil
		}
		if err := ildap.User.Add(&syncUser); err != nil {
			common.Log.Warnf("Failed to sync user to LDAP: %v", err)
			common.DB.Model(&syncUser).Update("sync_state", 2)
		} else {
			common.DB.Model(&syncUser).Update("sync_state", 1)
		}
	}
	
	return nil
}

// GetNextUIDNumber returns the next available UID (starting from 1000)
func (s UserService) GetNextUIDNumber() (uint, error) {
	var maxUID uint
	err := common.DB.Model(&model.User{}).
		Select("COALESCE(MAX(uid_number), 0)").
		Scan(&maxUID).Error
	if err != nil {
		return 0, err
	}
	
	if maxUID < 1000 {
		return 1000, nil
	}
	return maxUID + 1, nil
}

// List returns paginated user list
func (s UserService) List(req *request.UserListReq) ([]*model.User, error) {
	var list []*model.User
	db := common.DB.Model(&model.User{}).Order("id DESC")

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if nickname := strings.TrimSpace(req.Nickname); nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if mobile := strings.TrimSpace(req.Mobile); mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	if len(req.DepartmentId) > 0 {
		db = db.Where("department_id = ?", req.DepartmentId)
	}
	if givenName := strings.TrimSpace(req.GivenName); givenName != "" {
		db = db.Where("given_name LIKE ?", fmt.Sprintf("%%%s%%", givenName))
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.SyncState != 0 {
		db = db.Where("sync_state = ?", req.SyncState)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("Roles").Find(&list).Error
	return list, err
}

// ListCount returns total count matching filter
func (s UserService) ListCount(req *request.UserListReq) (int64, error) {
	var count int64
	db := common.DB.Model(&model.User{}).Order("id DESC")

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if nickname := strings.TrimSpace(req.Nickname); nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if mobile := strings.TrimSpace(req.Mobile); mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	if len(req.DepartmentId) > 0 {
		db = db.Where("department_id = ?", req.DepartmentId)
	}
	if givenName := strings.TrimSpace(req.GivenName); givenName != "" {
		db = db.Where("given_name LIKE ?", fmt.Sprintf("%%%s%%", givenName))
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.SyncState != 0 {
		db = db.Where("sync_state = ?", req.SyncState)
	}

	err := db.Count(&count).Error
	return count, err
}

// ListAll returns all users
func (s UserService) ListAll() (list []*model.User, err error) {
	err = common.DB.Model(&model.User{}).Order("created_at DESC").Find(&list).Error
	return list, err
}

// Count returns total user count
func (s UserService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.User{}).Count(&count).Error
	return count, err
}

// Exist checks if user exists
func (s UserService) Exist(filter map[string]any) bool {
	var dataObj model.User
	err := common.DB.Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// FindTheSameUserName finds user with similar username pattern
func (s UserService) FindTheSameUserName(username string, data *model.User) error {
	return common.DB.Where("username REGEXP ? ", fmt.Sprintf("^%s[0-9]{0,3}$", username)).Order("username desc").First(&data).Error
}

// Find gets single user by filter
func (s UserService) Find(filter map[string]any, data *model.User) error {
	return common.DB.Where(filter).Preload("Roles").First(&data).Error
}

// Update updates user and syncs to LDAP
func (s UserService) Update(user *model.User, skipLdapSync ...bool) error {
	shouldSkipLdapSync := len(skipLdapSync) > 0 && skipLdapSync[0]

	var oldUser model.User
	err := common.DB.Where("id = ?", user.ID).First(&oldUser).Error
	if err != nil {
		return err
	}

	// Auto-assign Unix attributes if not set
	if user.UIDNumber == 0 {
		uid, err := s.GetNextUIDNumber()
		if err != nil {
			common.Log.Warnf("Failed to assign UID for [%s]: %v", user.Username, err)
		} else {
			user.UIDNumber = uid
			if user.GIDNumber == 0 {
				user.GIDNumber = user.UIDNumber
			}
			if user.HomeDirectory == "" {
				user.HomeDirectory = fmt.Sprintf("/home/%s", user.Username)
			}
			if user.LoginShell == "" {
				user.LoginShell = "/bin/bash"
			}
			if user.Gecos == "" {
				if user.Nickname != "" {
					user.Gecos = user.Nickname
				} else {
					user.Gecos = user.Username
				}
			}
		}
	}

	// Update MySQL
	err = common.DB.Model(user).
		Select("username", "nickname", "given_name", "mail", "job_number", "mobile", "avatar",
			"postal_address", "departments", "position", "introduction", "creator", "department_id",
			"source", "user_dn", "uid_number", "gid_number", "home_directory", "login_shell",
			"gecos", "sync_state", "updated_at").
		Updates(user).Error
	if err != nil {
		return err
	}

	// Explicit GIDNumber update
	if user.GIDNumber != oldUser.GIDNumber {
		err = common.DB.Model(user).Where("id = ?", user.ID).Update("gid_number", user.GIDNumber).Error
		if err != nil {
			return err
		}
	}
	
	err = common.DB.Model(user).Association("Roles").Replace(user.Roles)
	if err != nil {
		return err
	}

	// Sync to LDAP
	if !shouldSkipLdapSync && oldUser.UserDN != "" {
		oldUsername := oldUser.Username
		newUsername := user.Username
		
		if !config.Conf.Ldap.UserNameModify {
			newUsername = oldUsername
		}

		syncUser := *user
		syncUser.Username = newUsername

		err = ildap.User.Update(oldUsername, &syncUser)
		if err != nil {
			common.Log.Warnf("LDAP sync failed for [%s]: %s", user.Username, err.Error())
			_ = s.ChangeSyncState(int(user.ID), 2)
		} else {
			_ = s.ChangeSyncState(int(user.ID), 1)
		}
	}

	// Update cache
	if err == nil {
		userDb := &model.User{}
		common.DB.Where("id = ?", user.ID).Preload("Roles").First(&userDb)
		userInfoCache.Set(userDb.Username, *userDb, cache.DefaultExpiration)
		if oldUser.Username != userDb.Username {
			userInfoCache.Delete(oldUser.Username)
		}
	}
	return err
}

// GetUserMinRoleSortsByIds gets minimum role sort for users
func (s UserService) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("no users found")
	}
	var roleMinSortList []int
	for _, user := range userList {
		var roleSortList []int
		for _, role := range user.Roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList).(int)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// GetCurrentUserMinRoleSort gets current user's highest role level
func (s UserService) GetCurrentUserMinRoleSort(c *gin.Context) (uint, model.User, error) {
	ctxUser, err := s.GetCurrentLoginUser(c)
	if err != nil {
		return 999, ctxUser, err
	}
	var currentRoleSorts []int
	for _, role := range ctxUser.Roles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts).(int))
	return currentRoleSortMin, ctxUser, nil
}

// Delete removes users from both MySQL and LDAP
func (s UserService) Delete(ids []uint) error {
	var users []model.User
	for _, id := range ids {
		user := new(model.User)
		err := s.Find(tools.H{"id": id}, user)
		if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}
		users = append(users, *user)
	}

	// Delete from LDAP first
	for _, user := range users {
		if user.UserDN != "" && user.UserDN != config.Conf.Ldap.AdminDN {
			if err := ildap.User.Delete(user.UserDN); err != nil {
				common.Log.Warnf("Failed to delete LDAP user [%s]: %v", user.UserDN, err)
			}
		}
	}

	// Delete from MySQL
	err := common.DB.Select("Roles").Unscoped().Delete(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		userInfoCache.Delete(user.Username)
	}

	// Remove group associations
	return common.DB.Exec("DELETE FROM group_users WHERE user_id IN (?)", ids).Error
}

// GetUserByIds gets users by IDs
func (s UserService) GetUserByIds(ids []uint) ([]model.User, error) {
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	return userList, err
}

// ChangePwd updates user password
func (s UserService) ChangePwd(username string, hashNewPasswd string) error {
	err := common.DB.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	
	cacheUser, found := userInfoCache.Get(username)
	if err == nil {
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		} else {
			var user model.User
			common.DB.Where("username = ?", username).Preload("Roles").First(&user)
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		}
	}
	return err
}

// ChangeStatus updates user status
func (s UserService) ChangeStatus(id, status int) error {
	return common.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// ChangeSyncState updates user sync state
func (s UserService) ChangeSyncState(id, status int) error {
	return common.DB.Model(&model.User{}).Where("id = ?", id).Update("sync_state", status).Error
}

// GetCurrentLoginUser gets current logged-in user with caching
func (s UserService) GetCurrentLoginUser(c *gin.Context) (model.User, error) {
	var newUser model.User
	ctxUser, exist := c.Get("user")
	if !exist {
		return newUser, errors.New("user not logged in")
	}
	u, _ := ctxUser.(model.User)

	cacheUser, found := userInfoCache.Get(u.Username)
	var user model.User
	var err error
	if found {
		user = cacheUser.(model.User)
	} else {
		user, err = s.GetUserById(u.ID)
		if err != nil {
			userInfoCache.Delete(u.Username)
		} else {
			userInfoCache.Set(u.Username, user, cache.DefaultExpiration)
		}
	}
	return user, err
}

// Login validates user credentials
func (s UserService) Login(user *model.User) (*model.User, error) {
	var firstUser model.User
	err := s.Find(tools.H{"username": user.Username}, &firstUser)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	// Prevent uid=admin format login
	if strings.HasPrefix(firstUser.UserDN, "uid=admin,") {
		return nil, errors.New("uid=admin format login not allowed")
	}
	
	if firstUser.Status != 1 {
		return nil, errors.New("user disabled")
	}

	if tools.NewParPasswd(firstUser.Password) != user.Password {
		return nil, errors.New("incorrect password")
	}

	return &firstUser, nil
}

// ClearUserInfoCache clears all user cache
func (s UserService) ClearUserInfoCache() {
	userInfoCache.Flush()
}

// GetUserById gets user by ID
func (us UserService) GetUserById(id uint) (model.User, error) {
	var user model.User
	err := common.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}
