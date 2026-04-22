package logic

import (
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/model/request"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/gin-gonic/gin"
)

type SqlLogic struct{}

// SyncSqlUsers syncs MySQL users to LDAP
func (d *SqlLogic) SyncSqlUsers(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SyncSqlUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	var users []*model.User
	var err error

	if len(r.UserIds) == 0 {
		// Full sync: sync all users and delete extras from LDAP
		var allUsers []*model.User
		err = common.DB.Unscoped().Model(&model.User{}).Order("created_at DESC").Find(&allUsers).Error
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %v", err))
		}

		users = []*model.User{}
		for _, u := range allUsers {
			if u.DeletedAt.Time.IsZero() {
				users = append(users, u)
			}
		}

		// Get LDAP users
		ldapUsers, err := ildap.User.ListUserDN()
		if err != nil {
			return nil, tools.NewLdapError(fmt.Errorf("failed to get LDAP users: %v", err))
		}

		// Create mappings
		mysqlUserDNMap := make(map[string]*model.User)
		mysqlUsernameMap := make(map[string]*model.User)
		for i := range allUsers {
			if allUsers[i].UserDN != "" {
				mysqlUserDNMap[allUsers[i].UserDN] = allUsers[i]
			}
			if allUsers[i].Username != "" {
				mysqlUsernameMap[allUsers[i].Username] = allUsers[i]
			}
		}

		// Delete LDAP users not in MySQL
		deletedCount := 0
		for _, ldapUser := range ldapUsers {
			if ldapUser.UserDN == config.Conf.Ldap.AdminDN {
				continue
			}

			username := extractUsernameFromDN(ldapUser.UserDN)
			mysqlUser, existsByDN := mysqlUserDNMap[ldapUser.UserDN]
			mysqlUserByUsername, existsByUsername := mysqlUsernameMap[username]

			shouldDelete := false
			if existsByDN {
				isDeleted := mysqlUser.DeletedAt.Valid && !mysqlUser.DeletedAt.Time.IsZero()
				if isDeleted || mysqlUser.Status != 1 {
					shouldDelete = true
				}
			} else if existsByUsername {
				isDeleted := mysqlUserByUsername.DeletedAt.Valid && !mysqlUserByUsername.DeletedAt.Time.IsZero()
				if isDeleted || mysqlUserByUsername.Status != 1 {
					shouldDelete = true
				} else {
					shouldDelete = true // Old format DN
				}
			} else {
				shouldDelete = true
			}

			if shouldDelete {
				if err := ildap.User.Delete(ldapUser.UserDN); err == nil {
					deletedCount++
				}
			}
		}

		// Sync MySQL users to LDAP (smart sync)
		syncedCount := 0
		updatedCount := 0
		createdCount := 0
		skippedCount := 0
		for _, user := range users {
			if user.UserDN == config.Conf.Ldap.AdminDN || user.Status != 1 {
				continue
			}

			existsInLDAP, _ := ildap.User.Exist(tools.H{"dn": user.UserDN})
			if existsInLDAP {
				// Check if update is needed
				ldapAttrs, err := ildap.User.GetUserDetails(user.UserDN)
				if err != nil {
					_ = isql.User.ChangeSyncState(int(user.ID), 2)
					continue
				}

				if ildap.User.NeedsUpdate(user, ldapAttrs) {
					if err := ildap.User.Update(user.Username, user); err == nil {
						updatedCount++
						_ = isql.User.ChangeSyncState(int(user.ID), 1)
					} else {
						_ = isql.User.ChangeSyncState(int(user.ID), 2)
					}
				} else {
					skippedCount++ // No changes
					_ = isql.User.ChangeSyncState(int(user.ID), 1)
				}
			} else {
				if err := ildap.User.Add(user); err == nil {
					createdCount++
					_ = isql.User.ChangeSyncState(int(user.ID), 1)
				} else {
					_ = isql.User.ChangeSyncState(int(user.ID), 2)
				}
			}
			syncedCount++
		}

		return map[string]int{"synced": syncedCount, "created": createdCount, "updated": updatedCount, "skipped": skippedCount, "deleted": deletedCount}, nil
	}

	// Sync specific users
	for _, id := range r.UserIds {
		if !isql.User.Exist(tools.H{"id": int(id)}) {
			return nil, tools.NewMySqlError(fmt.Errorf("user ID %d not found", id))
		}
	}
	userList, err := isql.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get users: %v", err))
	}
	users = make([]*model.User, len(userList))
	for i := range userList {
		users[i] = &userList[i]
	}

	syncedCount := 0
	updatedCount := 0
	createdCount := 0
	skippedCount := 0
	for _, user := range users {
		if user.UserDN == config.Conf.Ldap.AdminDN || user.Status != 1 {
			continue
		}

		exists, _ := ildap.User.Exist(tools.H{"dn": user.UserDN})
		if exists {
			// Check if update is needed
			ldapAttrs, err := ildap.User.GetUserDetails(user.UserDN)
			if err != nil {
				_ = isql.User.ChangeSyncState(int(user.ID), 2)
				continue
			}

			if ildap.User.NeedsUpdate(user, ldapAttrs) {
				if err := ildap.User.Update(user.Username, user); err == nil {
					updatedCount++
				} else {
					_ = isql.User.ChangeSyncState(int(user.ID), 2)
					continue
				}
			} else {
				skippedCount++ // No changes
			}
		} else {
			if err := ildap.User.Add(user); err == nil {
				createdCount++
			} else {
				_ = isql.User.ChangeSyncState(int(user.ID), 2)
				continue
			}
		}

		groups, err := isql.Group.GetGroupByIds(tools.StringToSlice(user.DepartmentId, ","))
		if err == nil {
			for _, group := range groups {
				_ = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
			}
		}

		_ = isql.User.ChangeSyncState(int(user.ID), 1)
		syncedCount++
	}

	return map[string]interface{}{"synced": syncedCount, "created": createdCount, "updated": updatedCount, "skipped": skippedCount}, nil
}

// SyncSqlGroups syncs MySQL groups to LDAP
func (d *SqlLogic) SyncSqlGroups(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.SyncSqlGrooupsReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.GroupIds {
		if !isql.Group.Exist(tools.H{"id": int(id)}) {
			return nil, tools.NewMySqlError(fmt.Errorf("group ID %d not found", id))
		}
	}

	groups, err := isql.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("failed to get groups: %v", err))
	}

	for _, group := range groups {
		// Ensure Users are loaded
		if err := common.DB.Preload("Users").First(group, group.ID).Error; err != nil {
			common.Log.Warnf("Failed to load users for group %s: %v", group.GroupName, err)
		}
		
		// Check if group already exists in LDAP
		searchConn, err := common.GetLDAPConn()
		groupExists := false
		if err == nil {
			searchRequest := ldap.NewSearchRequest(group.GroupDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
			_, searchErr := searchConn.Search(searchRequest)
			common.PutLADPConn(searchConn)
			
			if searchErr == nil {
				// Group already exists, just sync members
				groupExists = true
				common.Log.Infof("Group %s already exists in LDAP, syncing members only", group.GroupDN)
			}
		}
		
		if !groupExists {
			// Group doesn't exist, try to create it
			if group.GroupType == "cn" && len(group.Users) == 0 {
				common.Log.Warnf("Group %s (groupOfUniqueNames) has no users, skipping LDAP creation. Add users first.", group.GroupDN)
				// Mark as pending sync, will be created when first user is added
				_ = isql.Group.ChangeSyncState(int(group.ID), 2)
				continue
			}
			
			if err := ildap.Group.Add(group); err != nil {
				// If group already exists, continue with member sync
				if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
					common.Log.Infof("Group %s already exists in LDAP, syncing members only", group.GroupDN)
					groupExists = true
				} else {
					if strings.Contains(err.Error(), "uniqueMember") || strings.Contains(err.Error(), "Object Class Violation") {
						common.Log.Warnf("Group %s (groupOfUniqueNames) cannot be created without users: %v. Will be created when first user is added.", group.GroupDN, err)
						_ = isql.Group.ChangeSyncState(int(group.ID), 2)
						continue
					}
					return nil, tools.NewLdapError(fmt.Errorf("failed to sync group %s: %v", group.GroupName, err))
				}
			}
		}
		
		// Sync group members - MySQL as source of truth
		// Get current LDAP members
		conn, err := common.GetLDAPConnForModify()
		if err == nil {
			defer conn.Close()
			
			searchRequest := ldap.NewSearchRequest(
				group.GroupDN,
				ldap.ScopeBaseObject,
				ldap.NeverDerefAliases,
				0, 0, false,
				"(objectClass=*)",
				[]string{"objectClass", "uniqueMember", "memberUid", "sudoUser"},
				nil,
			)
			sr, searchErr := conn.Search(searchRequest)
			
			if searchErr == nil && len(sr.Entries) > 0 {
				entry := sr.Entries[0]
				objectClasses := entry.GetAttributeValues("objectClass")
				isPosixGroup := false
				isSudoRole := false
				isGroupOfUniqueNames := false
				for _, oc := range objectClasses {
					if oc == "posixGroup" {
						isPosixGroup = true
					}
					if oc == "sudoRole" {
						isSudoRole = true
					}
					if oc == "groupOfUniqueNames" {
						isGroupOfUniqueNames = true
					}
				}
				
				// Build MySQL user set
				mysqlUserDNs := make(map[string]bool)
				mysqlUsernames := make(map[string]bool)
				for _, user := range group.Users {
					if user.UserDN != "" && user.UserDN != config.Conf.Ldap.AdminDN {
						mysqlUserDNs[user.UserDN] = true
					}
					if user.Username != "" && user.Username != "admin" {
						mysqlUsernames[user.Username] = true
					}
				}
				
				modify := ldap.NewModifyRequest(group.GroupDN, nil)
				needUpdate := false
				
				if isSudoRole {
					// Handle sudoRole
					currentSudoUsers := entry.GetAttributeValues("sudoUser")
					currentSet := make(map[string]bool)
					for _, u := range currentSudoUsers {
						if u != "ALL" {
							currentSet[u] = true
						}
					}
					
					// Remove users not in MySQL
					toRemove := make([]string, 0)
					for u := range currentSet {
						if !mysqlUsernames[u] {
							toRemove = append(toRemove, u)
							needUpdate = true
						}
					}
					if len(toRemove) > 0 {
						modify.Delete("sudoUser", toRemove)
					}
					
					// Add users in MySQL but not in LDAP
					toAdd := make([]string, 0)
					for u := range mysqlUsernames {
						if !currentSet[u] {
							toAdd = append(toAdd, u)
							needUpdate = true
						}
					}
					if len(toAdd) > 0 {
						// If current has only "ALL", replace it
						if len(currentSudoUsers) == 1 && currentSudoUsers[0] == "ALL" {
							modify.Replace("sudoUser", toAdd)
						} else {
							modify.Add("sudoUser", toAdd)
						}
					}
					
					if needUpdate {
						if err := conn.Modify(modify); err != nil {
							common.Log.Warnf("Failed to sync sudoRole members for %s: %v", group.GroupDN, err)
						}
					}
				} else if isPosixGroup {
					// Handle posixGroup
					currentMembers := entry.GetAttributeValues("memberUid")
					currentSet := make(map[string]bool)
					for _, m := range currentMembers {
						currentSet[m] = true
					}
					
					// Remove users not in MySQL
					toRemove := make([]string, 0)
					for m := range currentSet {
						if m != "admin" && !mysqlUsernames[m] {
							toRemove = append(toRemove, m)
							needUpdate = true
						}
					}
					if len(toRemove) > 0 {
						modify.Delete("memberUid", toRemove)
					}
					
					// Add users in MySQL but not in LDAP
					toAdd := make([]string, 0)
					for u := range mysqlUsernames {
						if !currentSet[u] {
							toAdd = append(toAdd, u)
							needUpdate = true
						}
					}
					if len(toAdd) > 0 {
						modify.Add("memberUid", toAdd)
					}
					
					if needUpdate {
						if err := conn.Modify(modify); err != nil {
							common.Log.Warnf("Failed to sync posixGroup members for %s: %v", group.GroupDN, err)
						}
					}
				} else if isGroupOfUniqueNames {
					// Handle groupOfUniqueNames
					currentMembers := entry.GetAttributeValues("uniqueMember")
					currentSet := make(map[string]bool)
					for _, m := range currentMembers {
						if m != config.Conf.Ldap.AdminDN {
							currentSet[m] = true
						}
					}
					
					// Remove users not in MySQL
					toRemove := make([]string, 0)
					for m := range currentSet {
						if !mysqlUserDNs[m] {
							toRemove = append(toRemove, m)
							needUpdate = true
						}
					}
					if len(toRemove) > 0 {
						// Don't remove if it's the last member (groupOfUniqueNames requires at least one)
						if len(currentMembers)-len(toRemove) > 0 {
							modify.Delete("uniqueMember", toRemove)
						} else {
							common.Log.Warnf("Cannot remove all members from groupOfUniqueNames %s, keeping existing members", group.GroupDN)
							needUpdate = false
						}
					}
					
					// Add users in MySQL but not in LDAP
					toAdd := make([]string, 0)
					for u := range mysqlUserDNs {
						if !currentSet[u] {
							toAdd = append(toAdd, u)
							needUpdate = true
						}
					}
					if len(toAdd) > 0 {
						modify.Add("uniqueMember", toAdd)
					}
					
					if needUpdate {
						if err := conn.Modify(modify); err != nil {
							common.Log.Warnf("Failed to sync groupOfUniqueNames members for %s: %v", group.GroupDN, err)
						}
					}
				}
			} else {
				// Group doesn't exist in LDAP, add all MySQL users
				for _, user := range group.Users {
					if user.UserDN != "" && user.UserDN != config.Conf.Ldap.AdminDN {
						_ = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
					}
				}
			}
		}
		_ = isql.Group.ChangeSyncState(int(group.ID), 1)
	}

	return nil, nil
}

// SearchGroupDiff finds groups not synced to LDAP
func SearchGroupDiff() (err error) {
	sqlGroupList, err := isql.Group.ListAll()
	if err != nil {
		return err
	}
	ldapGroupList, err := ildap.Group.ListGroupDN()
	if err != nil {
		return err
	}

	groups := diffGroup(sqlGroupList, ldapGroupList)
	for _, group := range groups {
		if group.GroupDN != config.Conf.Ldap.BaseDN {
			_ = isql.Group.ChangeSyncState(int(group.ID), 2)
		}
	}
	return
}

// SearchUserDiff syncs MySQL users to LDAP (smart sync - only updates changed users)
func SearchUserDiff() (err error) {
	var sqlUserList []*model.User
	err = common.DB.Model(&model.User{}).Preload("Roles").Find(&sqlUserList).Error
	if err != nil {
		return err
	}

	ldapUserList, err := ildap.User.ListUserDN()
	if err != nil {
		return err
	}

	// Build maps for quick lookup
	mysqlUserDNMap := make(map[string]*model.User)
	for i := range sqlUserList {
		if sqlUserList[i].UserDN != "" && sqlUserList[i].Status == 1 {
			mysqlUserDNMap[sqlUserList[i].UserDN] = sqlUserList[i]
		}
	}

	ldapUserDNMap := make(map[string]bool)
	for i := range ldapUserList {
		ldapUserDNMap[ldapUserList[i].UserDN] = true
	}

	// Sync MySQL users to LDAP
	for _, user := range sqlUserList {
		if user.UserDN == config.Conf.Ldap.AdminDN || user.UserDN == "" {
			continue
		}

		if user.Status != 1 {
			// Disabled user - delete from LDAP if exists
			if ldapUserDNMap[user.UserDN] {
				_ = ildap.User.Delete(user.UserDN)
			}
			_ = isql.User.ChangeSyncState(int(user.ID), 2)
			continue
		}

		if ldapUserDNMap[user.UserDN] {
			// User exists in LDAP - check if update is needed
			ldapAttrs, err := ildap.User.GetUserDetails(user.UserDN)
			if err != nil {
				// Cannot get LDAP details, skip update
				continue
			}

			if ildap.User.NeedsUpdate(user, ldapAttrs) {
				// Attributes changed, update LDAP (preserves password)
				if err := ildap.User.Update(user.Username, user); err != nil {
					_ = isql.User.ChangeSyncState(int(user.ID), 2)
				} else {
					_ = isql.User.ChangeSyncState(int(user.ID), 1)
				}
			}
			// No changes - do nothing
		} else {
			// New user - Add to LDAP
			if err := ildap.User.Add(user); err != nil {
				_ = isql.User.ChangeSyncState(int(user.ID), 2)
			} else {
				_ = isql.User.ChangeSyncState(int(user.ID), 1)
			}
		}
	}

	// Delete LDAP users not in MySQL (ensure LDAP <= MySQL)
	for _, ldapUser := range ldapUserList {
		if ldapUser.UserDN == config.Conf.Ldap.AdminDN {
			continue
		}
		if _, exists := mysqlUserDNMap[ldapUser.UserDN]; !exists {
			_ = ildap.User.Delete(ldapUser.UserDN)
		}
	}

	return
}

func diffGroup(sqlGroup, ldapGroup []*model.Group) (rst []*model.Group) {
	tmp := make(map[string]struct{})
	for _, v := range ldapGroup {
		tmp[v.GroupDN] = struct{}{}
	}
	for _, v := range sqlGroup {
		if _, ok := tmp[v.GroupDN]; !ok {
			rst = append(rst, v)
		}
	}
	return
}

func extractUsernameFromDN(dn string) string {
	parts := strings.Split(dn, ",")
	if len(parts) > 0 {
		uidPart := parts[0]
		if strings.HasPrefix(uidPart, "uid=") {
			return strings.TrimPrefix(uidPart, "uid=")
		}
	}
	return ""
}
