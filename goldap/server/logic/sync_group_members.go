package logic

import (
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/public/common"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	ldap "github.com/go-ldap/ldap/v3"
)

// SyncGroupMembersFromDB syncs all group members from database to LDAP
func SyncGroupMembersFromDB() error {
	// First sync sudo posixGroup based on sudoRole membership
	if err := SyncSudoPosixGroupFromDB(); err != nil {
		common.Log.Warnf("Failed to sync sudo posixGroup: %v", err)
	}

	// Sync docker posixGroup
	if err := SyncDockerPosixGroupFromDB(); err != nil {
		common.Log.Warnf("Failed to sync docker posixGroup: %v", err)
	}

	groups, err := isql.Group.ListAll()
	if err != nil {
		return fmt.Errorf("failed to get group list: %v", err)
	}

	syncedCount := 0
	failedCount := 0

	for _, group := range groups {
		if err = common.DB.Preload("Users").First(&group, group.ID).Error; err != nil {
			continue
		}

		isSudoRole := strings.Contains(group.GroupDN, "ou=sudoers,") && group.GroupType == "cn"

		if isSudoRole {
			if len(group.Users) == 0 {
				// No members, delete sudoRole if exists
				conn, err := common.GetLDAPConnForModify()
				if err != nil {
					failedCount++
					continue
				}

				searchRequest := ldap.NewSearchRequest(group.GroupDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
				_, err = conn.Search(searchRequest)
				if err == nil {
					del := ldap.NewDelRequest(group.GroupDN, nil)
					_ = conn.Del(del)
					syncedCount++
				}
				conn.Close()
				continue
			}

			// Get current sudoUser values
			conn, err := common.GetLDAPConnForModify()
			if err != nil {
				failedCount++
				continue
			}

			searchRequest := ldap.NewSearchRequest(group.GroupDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"sudoUser"}, nil)
			sr, err := conn.Search(searchRequest)
			if err != nil || len(sr.Entries) == 0 {
				conn.Close()
				failedCount++
				continue
			}

			currentSudoUsers := sr.Entries[0].GetAttributeValues("sudoUser")
			conn.Close()

			// Build username list from database
			dbUsernames := make([]string, 0, len(group.Users))
			for _, user := range group.Users {
				if user.Username != "" {
					dbUsernames = append(dbUsernames, user.Username)
				}
			}

			// Compare and update if different
			needUpdate := len(currentSudoUsers) != len(dbUsernames)
			if !needUpdate {
				for _, dbUsername := range dbUsernames {
					found := false
					for _, ldapUser := range currentSudoUsers {
						if ldapUser == dbUsername {
							found = true
							break
						}
					}
					if !found {
						needUpdate = true
						break
					}
				}
			}

			if needUpdate {
				modifyConn, err := common.GetLDAPConnForModify()
				if err != nil {
					failedCount++
					continue
				}

				modify := ldap.NewModifyRequest(group.GroupDN, nil)
				modify.Replace("sudoUser", dbUsernames)
				err = modifyConn.Modify(modify)
				modifyConn.Close()

				if err != nil {
					failedCount++
				} else {
					syncedCount++
				}
			} else {
				syncedCount++
			}
		} else {
			isOUGroup := len(group.GroupDN) >= 3 && group.GroupDN[:3] == "ou="
			if isOUGroup {
				syncedCount++
				continue
			}
			
			conn, err := common.GetLDAPConnForModify()
			if err != nil {
				failedCount++
				continue
			}
			
			searchRequest := ldap.NewSearchRequest(
				group.GroupDN,
				ldap.ScopeBaseObject,
				ldap.NeverDerefAliases,
				0, 0, false,
				"(objectClass=*)",
				[]string{"objectClass", "uniqueMember", "memberUid"},
				nil,
			)
			sr, err := conn.Search(searchRequest)
			if err != nil || len(sr.Entries) == 0 {
				conn.Close()
				for _, user := range group.Users {
					if user.UserDN != "" && user.UserDN != config.Conf.Ldap.AdminDN {
						_ = ildap.Group.AddUserToGroup(group.GroupDN, user.UserDN)
					}
				}
				syncedCount++
				continue
			}
			
			entry := sr.Entries[0]
			objectClasses := entry.GetAttributeValues("objectClass")
			isPosixGroup := false
			for _, oc := range objectClasses {
				if oc == "posixGroup" {
					isPosixGroup = true
					break
				}
			}
			
			// Build MySQL user sets
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
			
			if isPosixGroup {
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
			} else {
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
			}
			
			if needUpdate {
				if err := conn.Modify(modify); err != nil {
					common.Log.Warnf("Failed to sync group members for %s: %v", group.GroupDN, err)
					failedCount++
				} else {
					syncedCount++
				}
			} else {
				syncedCount++
			}
			conn.Close()
		}
	}

	return nil
}

// SyncSudoPosixGroupFromDB syncs sudo posixGroup membership from MySQL sudoRole groups
func SyncSudoPosixGroupFromDB() error {
	// Get all users in sudouser-nopasswd and sudouser-other groups from MySQL
	sudoUsers := make(map[string]bool)

	groups, err := isql.Group.ListAll()
	if err != nil {
		return fmt.Errorf("failed to get group list: %v", err)
	}

	for _, group := range groups {
		// Check if this is a sudoRole group
		isSudoRole := strings.Contains(group.GroupDN, "ou=sudoers,") && group.GroupType == "cn"
		if !isSudoRole {
			continue
		}

		// Load users for this group
		if err = common.DB.Preload("Users").First(&group, group.ID).Error; err != nil {
			continue
		}

		for _, user := range group.Users {
			if user.Username != "" && user.Username != "admin" {
				sudoUsers[user.Username] = true
			}
		}
	}

	// Get LDAP connection
	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return fmt.Errorf("failed to get LDAP connection: %v", err)
	}
	defer conn.Close()

	sudoGroupDN := fmt.Sprintf("cn=sudo,%s", config.Conf.Ldap.BaseDN)

	// Check if sudo posixGroup exists
	searchRequest := ldap.NewSearchRequest(
		sudoGroupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"memberUid"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	
	// If no sudo users, delete the group if it exists
	if len(sudoUsers) == 0 {
		if err == nil && len(sr.Entries) > 0 {
			// Group exists but no sudo users, delete it
			del := ldap.NewDelRequest(sudoGroupDN, nil)
			if delErr := conn.Del(del); delErr != nil {
				// If group doesn't exist or already deleted, that's fine
				if ldapErr, ok := delErr.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
					return nil
				}
				common.Log.Warnf("Failed to delete empty sudo posixGroup: %v", delErr)
			}
		}
		return nil
	}

	// Convert map to slice
	memberUids := make([]string, 0, len(sudoUsers))
	for username := range sudoUsers {
		memberUids = append(memberUids, username)
	}

	if err != nil {
		// Group doesn't exist, create it
		add := ldap.NewAddRequest(sudoGroupDN, nil)
		add.Attribute("objectClass", []string{"posixGroup", "top"})
		add.Attribute("cn", []string{"sudo"})
		add.Attribute("gidNumber", []string{"27"})
		add.Attribute("description", []string{"Sudo posixGroup (GID=27)"})
		add.Attribute("memberUid", memberUids)
		if err := conn.Add(add); err != nil {
			return fmt.Errorf("failed to create sudo posixGroup: %v", err)
		}
		return nil
	}

	// Group exists, update memberUid
	if len(sr.Entries) > 0 {
		currentMembers := sr.Entries[0].GetAttributeValues("memberUid")

		// Check if update is needed
		needUpdate := len(currentMembers) != len(memberUids)
		if !needUpdate {
			currentSet := make(map[string]bool)
			for _, m := range currentMembers {
				currentSet[m] = true
			}
			for _, m := range memberUids {
				if !currentSet[m] {
					needUpdate = true
					break
				}
			}
		}

		if needUpdate {
			modify := ldap.NewModifyRequest(sudoGroupDN, nil)
			modify.Replace("memberUid", memberUids)
			if err := conn.Modify(modify); err != nil {
				return fmt.Errorf("failed to update sudo posixGroup: %v", err)
			}
		}
	}

	return nil
}

// SyncDockerPosixGroupFromDB syncs docker posixGroup membership from MySQL docker group
func SyncDockerPosixGroupFromDB() error {
	// Find docker group in MySQL (posix type with name "docker")
	groups, err := isql.Group.ListAll()
	if err != nil {
		return fmt.Errorf("failed to get group list: %v", err)
	}

	var dockerGroup *struct {
		ID        uint
		GroupName string
		GroupType string
		GIDNumber uint
		GroupDN   string
	}

	for _, group := range groups {
		if group.GroupName == "docker" && group.GroupType == "posix" {
			dockerGroup = &struct {
				ID        uint
				GroupName string
				GroupType string
				GIDNumber uint
				GroupDN   string
			}{
				ID:        group.ID,
				GroupName: group.GroupName,
				GroupType: group.GroupType,
				GIDNumber: group.GIDNumber,
				GroupDN:   group.GroupDN,
			}
			break
		}
	}

	if dockerGroup == nil {
		// No docker group configured
		return nil
	}

	// Load users for docker group using raw query on join table
	var usernames []string
	err = common.DB.Table("group_users").
		Select("users.username").
		Joins("JOIN users ON users.id = group_users.user_id").
		Where("group_users.group_id = ? AND users.deleted_at IS NULL AND users.username != ?", dockerGroup.ID, "admin").
		Pluck("users.username", &usernames).Error
	if err != nil {
		return fmt.Errorf("failed to load docker group users: %v", err)
	}

	// Build username list
	memberUids := make([]string, 0)
	for _, username := range usernames {
		if username != "" {
			memberUids = append(memberUids, username)
		}
	}

	// Get LDAP connection
	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return fmt.Errorf("failed to get LDAP connection: %v", err)
	}
	defer conn.Close()

	dockerGroupDN := fmt.Sprintf("cn=docker,%s", config.Conf.Ldap.BaseDN)

	// Check if docker posixGroup exists
	searchRequest := ldap.NewSearchRequest(
		dockerGroupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"memberUid", "gidNumber"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		// Group doesn't exist, create it
		add := ldap.NewAddRequest(dockerGroupDN, nil)
		add.Attribute("objectClass", []string{"posixGroup", "top"})
		add.Attribute("cn", []string{"docker"})
		add.Attribute("gidNumber", []string{fmt.Sprintf("%d", dockerGroup.GIDNumber)})
		add.Attribute("description", []string{"Docker posixGroup"})
		if len(memberUids) > 0 {
			add.Attribute("memberUid", memberUids)
		}
		if err := conn.Add(add); err != nil {
			return fmt.Errorf("failed to create docker posixGroup: %v", err)
		}
		return nil
	}

	// Group exists, update memberUid
	if len(sr.Entries) > 0 {
		currentMembers := sr.Entries[0].GetAttributeValues("memberUid")

		// Check if update is needed
		needUpdate := len(currentMembers) != len(memberUids)
		if !needUpdate {
			currentSet := make(map[string]bool)
			for _, m := range currentMembers {
				currentSet[m] = true
			}
			for _, m := range memberUids {
				if !currentSet[m] {
					needUpdate = true
					break
				}
			}
			// Also check for members that should be removed
			if !needUpdate {
				newSet := make(map[string]bool)
				for _, m := range memberUids {
					newSet[m] = true
				}
				for _, m := range currentMembers {
					if !newSet[m] {
						needUpdate = true
						break
					}
				}
			}
		}

		if needUpdate {
			modify := ldap.NewModifyRequest(dockerGroupDN, nil)
			if len(memberUids) > 0 {
				modify.Replace("memberUid", memberUids)
			} else {
				// Remove all members if empty
				if len(currentMembers) > 0 {
					modify.Delete("memberUid", nil)
				}
			}
			if err := conn.Modify(modify); err != nil {
				return fmt.Errorf("failed to update docker posixGroup: %v", err)
			}
		}
	}

	return nil
}
