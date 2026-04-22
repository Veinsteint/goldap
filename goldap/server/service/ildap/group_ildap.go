package ildap

import (
	"errors"
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"

	ldap "github.com/go-ldap/ldap/v3"
)

type GroupService struct{}

// extractParentDN extracts parent DN from full DN
func extractParentDN(dn string) string {
	parts := strings.SplitN(dn, ",", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// extractOUName extracts OU name from DN
func extractOUName(dn string) string {
	parts := strings.Split(dn, ",")
	if len(parts) > 0 && strings.HasPrefix(parts[0], "ou=") {
		return strings.TrimPrefix(parts[0], "ou=")
	}
	return ""
}

// Add creates group in LDAP
func (x GroupService) Add(g *model.Group) error {
	if g.Remark == "" {
		g.Remark = g.GroupName
	}

	// Ensure parent DN exists
	parentDN := extractParentDN(g.GroupDN)
	if parentDN != "" && parentDN != config.Conf.Ldap.BaseDN {
		searchConn, err := common.GetLDAPConn()
		if err == nil {
			defer common.PutLADPConn(searchConn)
			searchRequest := ldap.NewSearchRequest(
				parentDN,
				ldap.ScopeBaseObject,
				ldap.NeverDerefAliases,
				0, 0, false,
				"(objectClass=*)",
				[]string{"dn"},
				nil,
			)
			_, err = searchConn.Search(searchRequest)
			if err != nil && strings.HasPrefix(parentDN, "ou=") {
				modifyConn, modErr := common.GetLDAPConnForModify()
				if modErr == nil {
					defer modifyConn.Close()
					ouName := extractOUName(parentDN)
					if ouName != "" {
						ouAdd := ldap.NewAddRequest(parentDN, nil)
						ouAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
						ouAdd.Attribute("ou", []string{ouName})
						addErr := modifyConn.Add(ouAdd)
						if addErr != nil {
							if ldapErr, ok := addErr.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
								common.Log.Warnf("Failed to create parent DN %s: %v", parentDN, addErr)
							}
						}
					}
				}
			}
		}
	}

	add := ldap.NewAddRequest(g.GroupDN, nil)
	isSudoRole := strings.Contains(g.GroupDN, "ou=sudoers,") && g.GroupType == "cn"

	if isSudoRole {
		sudoUser := "ALL"
		if len(g.Users) > 0 {
			sudoUser = g.Users[0].Username
		}
		
		sudoRule := &model.SudoRule{
			Name:        g.GroupName,
			Description: g.Remark,
			User:        sudoUser,
			Host:        "ALL",
			Command:     "ALL",
			RunAsUser:   "ALL",
			RunAsGroup:  "ALL",
			Options:     "",
			Creator:     "system",
		}
		if g.GroupName == "sudouser-nopasswd" {
			sudoRule.Options = "!authenticate"
		}
		return SudoService{}.Add(sudoRule)
	} else if g.GroupType == "posix" {
		add.Attribute("objectClass", []string{"posixGroup", "top"})
		add.Attribute("cn", []string{g.GroupName})
		add.Attribute("gidNumber", []string{fmt.Sprintf("%d", g.GIDNumber)})
		if g.Remark != "" {
			add.Attribute("description", []string{g.Remark})
		}
	} else if g.GroupType == "ou" {
		add.Attribute("objectClass", []string{"organizationalUnit", "top"})
		add.Attribute("ou", []string{g.GroupName})
		if g.Remark != "" {
			add.Attribute("description", []string{g.Remark})
		}
	} else if g.GroupType == "cn" {
		add.Attribute("objectClass", []string{"groupOfUniqueNames", "top"})
		add.Attribute("cn", []string{g.GroupName})
		if len(g.Users) > 0 {
			uniqueMembers := make([]string, 0)
			for _, u := range g.Users {
				if u.UserDN != "" {
					uniqueMembers = append(uniqueMembers, u.UserDN)
				}
			}
			if len(uniqueMembers) > 0 {
				add.Attribute("uniqueMember", uniqueMembers)
			} else {
				return fmt.Errorf("groupOfUniqueNames requires at least one uniqueMember, but no valid users found")
			}
		} else {
			return fmt.Errorf("groupOfUniqueNames requires at least one uniqueMember, but group has no users. Please add at least one user to the group first")
		}
		if g.Remark != "" {
			add.Attribute("description", []string{g.Remark})
		}
	} else {
		add.Attribute("objectClass", []string{"top"})
		add.Attribute(g.GroupType, []string{g.GroupName})
		if g.Remark != "" {
			add.Attribute("description", []string{g.Remark})
		}
	}

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Add(add)
	if err != nil {
		// If group already exists
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
			common.Log.Infof("Group %s already exists in LDAP", g.GroupDN)
			return nil
		}
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultInvalidAttributeSyntax {
			if isSudoRole {
				return fmt.Errorf("failed to create sudoRole: sudo schema may not be loaded. Error: %v", err)
			}
		}
		return err
	}

	return nil
}

// Update modifies group in LDAP
func (x GroupService) Update(oldGroup, newGroup *model.Group) error {
	searchConn, err := common.GetLDAPConn()
	if err != nil {
		return err
	}
	defer common.PutLADPConn(searchConn)

	searchRequest := ldap.NewSearchRequest(
		oldGroup.GroupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	_, err = searchConn.Search(searchRequest)
	if err != nil {
		// Group not found in LDAP, auto-creating...
		isSudoRole := strings.Contains(oldGroup.GroupDN, "ou=sudoers,") && oldGroup.GroupType == "cn"

		if isSudoRole {
			sudoRule := &model.SudoRule{
				Name:        oldGroup.GroupName,
				Description: oldGroup.Remark,
				User:        "",
				Host:        "ALL",
				Command:     "ALL",
				RunAsUser:   "ALL",
				RunAsGroup:  "ALL",
				Options:     "",
				Creator:     "system",
			}
			if oldGroup.GroupName == "sudouser-nopasswd" {
				sudoRule.Options = "!authenticate"
			}
			addErr := SudoService{}.Add(sudoRule)
			if addErr != nil {
				if ldapErr, ok := addErr.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
					return fmt.Errorf("failed to create sudoRole %s: %v", oldGroup.GroupDN, addErr)
				}
			}
		} else {
			createGroup := &model.Group{
				GroupDN:   oldGroup.GroupDN,
				GroupName: oldGroup.GroupName,
				GroupType: oldGroup.GroupType,
				Remark:    oldGroup.Remark,
			}
			addErr := x.Add(createGroup)
			if addErr != nil {
				if ldapErr, ok := addErr.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
					return fmt.Errorf("failed to create group %s: %v", oldGroup.GroupDN, addErr)
				}
			}
		}
	}

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	isSudoRole := strings.Contains(oldGroup.GroupDN, "ou=sudoers,") && oldGroup.GroupType == "cn"
	isPosixGroup := oldGroup.GroupType == "posix"

	if isSudoRole {
		modify1 := ldap.NewModifyRequest(oldGroup.GroupDN, nil)
		modify1.Replace("description", []string{newGroup.Remark})

		if oldGroup.GroupName == "sudouser-nopasswd" {
			modify1.Replace("sudoOption", []string{"!authenticate"})
		} else {
			modify1.Delete("sudoOption", []string{})
		}

		err = conn.Modify(modify1)
		if err != nil {
			return err
		}
	} else if isPosixGroup {
		modify1 := ldap.NewModifyRequest(oldGroup.GroupDN, nil)
		modify1.Replace("description", []string{newGroup.Remark})
		if newGroup.GIDNumber > 0 {
			modify1.Replace("gidNumber", []string{fmt.Sprintf("%d", newGroup.GIDNumber)})
		}
		err = conn.Modify(modify1)
		if err != nil {
			return err
		}
	} else {
		modify1 := ldap.NewModifyRequest(oldGroup.GroupDN, nil)
		modify1.Replace("description", []string{newGroup.Remark})
		err = conn.Modify(modify1)
		if err != nil {
			return err
		}
	}

	// Rename if allowed and changed
	if config.Conf.Ldap.GroupNameModify && newGroup.GroupName != oldGroup.GroupName {
		modify2 := ldap.NewModifyDNRequest(oldGroup.GroupDN, newGroup.GroupDN, true, "")
		return conn.ModifyDN(modify2)
	}
	return nil
}

// Delete removes group from LDAP
func (x GroupService) Delete(gdn string) error {
	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	del := ldap.NewDelRequest(gdn, nil)
	return conn.Del(del)
}

// AddUserToGroup adds user to group (supports groupOfUniqueNames, posixGroup, sudoRole)
func (x GroupService) AddUserToGroup(dn, udn string) error {
	if len(dn) < 3 || dn[:3] == "ou=" {
		return errors.New("cannot add user to OU")
	}

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"objectClass", "sudoUser", "memberUid", "uniqueMember"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil || len(sr.Entries) == 0 {
		// Group doesn't exist, try to create it
		dnParts := strings.Split(dn, ",")
		if len(dnParts) > 0 {
			firstPart := strings.Split(dnParts[0], "=")
			if len(firstPart) == 2 {
				groupType := firstPart[0]
				groupName := firstPart[1]
				
				// Handle sudoRole
				if strings.Contains(dn, "ou=sudoers,") && groupType == "cn" {
					username := extractUsername(udn)
					if username != "" {
						x.ensureSudoersOU(conn)

						sudoAdd := ldap.NewAddRequest(dn, nil)
						sudoAdd.Attribute("objectClass", []string{"top", "sudoRole"})
						sudoAdd.Attribute("cn", []string{groupName})
						sudoAdd.Attribute("sudoUser", []string{username})
						sudoAdd.Attribute("sudoHost", []string{"ALL"})
						sudoAdd.Attribute("sudoCommand", []string{"ALL"})
						sudoAdd.Attribute("sudoRunAsUser", []string{"ALL"})
						if groupName == "sudouser-nopasswd" {
							sudoAdd.Attribute("sudoOption", []string{"!authenticate"})
						}

						addErr := conn.Add(sudoAdd)
						if addErr != nil {
							return fmt.Errorf("failed to create sudoRole: %v", addErr)
						}
						return nil
					}
				}
				
				// Handle groupOfUniqueNames (cn type groups)
				if groupType == "cn" && !strings.Contains(dn, "ou=sudoers,") {
					// Ensure parent OU exists
					parentDN := extractParentDN(dn)
					if parentDN != "" && parentDN != config.Conf.Ldap.BaseDN {
						parentSearchRequest := ldap.NewSearchRequest(parentDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
						_, parentErr := conn.Search(parentSearchRequest)
						if parentErr != nil {
							ouName := extractOUName(parentDN)
							if ouName != "" {
								parentAdd := ldap.NewAddRequest(parentDN, nil)
								parentAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
								parentAdd.Attribute("ou", []string{ouName})
								if addErr := conn.Add(parentAdd); addErr != nil {
									if ldapErr, ok := addErr.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
										common.Log.Warnf("Failed to create parent DN %s: %v", parentDN, addErr)
									}
								}
							}
						}
					}
					
					// Create groupOfUniqueNames with the user as first member
					groupAdd := ldap.NewAddRequest(dn, nil)
					groupAdd.Attribute("objectClass", []string{"groupOfUniqueNames", "top"})
					groupAdd.Attribute("cn", []string{groupName})
					groupAdd.Attribute("uniqueMember", []string{udn})
					
					addErr := conn.Add(groupAdd)
					if addErr != nil {
						if ldapErr, ok := addErr.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
							// Group was created by another process, continue with normal flow
							// Re-search to get the group
							sr, err = conn.Search(searchRequest)
							if err != nil || len(sr.Entries) == 0 {
								return fmt.Errorf("failed to query group after creation: %v", err)
							}
						} else {
							return fmt.Errorf("failed to create groupOfUniqueNames: %v", addErr)
						}
					} else {
						// Group created successfully, return early
						return nil
					}
				}
			}
		}
		
		// If we reach here and group still doesn't exist, return error
		if err != nil {
			return fmt.Errorf("failed to query group: %v", err)
		}
		if len(sr.Entries) == 0 {
			return errors.New("group not found and could not be created")
		}
	}

	// Determine group type
	objectClasses := sr.Entries[0].GetAttributeValues("objectClass")
	isPosixGroup := false
	isSudoRole := false
	for _, oc := range objectClasses {
		if oc == "posixGroup" {
			isPosixGroup = true
		}
		if oc == "sudoRole" {
			isSudoRole = true
		}
	}

	username := extractUsername(udn)
	if (isPosixGroup || isSudoRole) && username == "" {
		return errors.New("cannot extract username from DN")
	}

	newmr := ldap.NewModifyRequest(dn, nil)
	if isSudoRole {
		currentSudoUsers := sr.Entries[0].GetAttributeValues("sudoUser")
		for _, u := range currentSudoUsers {
			if u == username {
				return nil // Already exists
			}
		}
		if len(currentSudoUsers) == 1 && currentSudoUsers[0] == "ALL" {
			newmr.Replace("sudoUser", []string{username})
		} else {
			newmr.Add("sudoUser", []string{username})
		}
		// Also add user to sudo posixGroup (GID=27) for 'id' command display
		_ = x.ensureSudoPosixGroup(conn)
		_ = x.addUserToSudoPosixGroup(conn, username)
	} else if isPosixGroup {
		// Check if user already in group
		currentMembers := sr.Entries[0].GetAttributeValues("memberUid")
		for _, m := range currentMembers {
			if m == username {
				return nil 
			}
		}
		newmr.Add("memberUid", []string{username})
	} else {
		// Check if user already in groupOfUniqueNames
		currentMembers := sr.Entries[0].GetAttributeValues("uniqueMember")
		for _, m := range currentMembers {
			if m == udn {
				return nil
			}
		}
		newmr.Add("uniqueMember", []string{udn})
	}

	return conn.Modify(newmr)
}

// ensureSudoersOU ensures ou=sudoers exists
func (x GroupService) ensureSudoersOU(conn *ldap.Conn) {
	sudoersDN := fmt.Sprintf("ou=sudoers,%s", config.Conf.Ldap.BaseDN)
	sudoersSearchRequest := ldap.NewSearchRequest(
		sudoersDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)
	_, err := conn.Search(sudoersSearchRequest)
	if err != nil {
		sudoersAdd := ldap.NewAddRequest(sudoersDN, nil)
		sudoersAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
		sudoersAdd.Attribute("ou", []string{"sudoers"})
		sudoersAdd.Attribute("description", []string{"Sudo rules container"})
		_ = conn.Add(sudoersAdd)
	}
}

// ensureSudoPosixGroup ensures cn=sudo posixGroup exists (GID=27)
func (x GroupService) ensureSudoPosixGroup(conn *ldap.Conn) error {
	sudoGroupDN := fmt.Sprintf("cn=sudo,%s", config.Conf.Ldap.BaseDN)
	searchRequest := ldap.NewSearchRequest(
		sudoGroupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"dn"},
		nil,
	)
	_, err := conn.Search(searchRequest)
	if err != nil {
		// Create sudo posixGroup
		add := ldap.NewAddRequest(sudoGroupDN, nil)
		add.Attribute("objectClass", []string{"posixGroup", "top"})
		add.Attribute("cn", []string{"sudo"})
		add.Attribute("gidNumber", []string{"27"})
		add.Attribute("description", []string{"Sudo用户组(posixGroup GID=27)"})
		return conn.Add(add)
	}
	return nil
}

// addUserToSudoPosixGroup adds user to cn=sudo posixGroup
func (x GroupService) addUserToSudoPosixGroup(conn *ldap.Conn, username string) error {
	sudoGroupDN := fmt.Sprintf("cn=sudo,%s", config.Conf.Ldap.BaseDN)
	
	// Check if user already in group
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
	if err != nil {
		return err
	}
	if len(sr.Entries) > 0 {
		members := sr.Entries[0].GetAttributeValues("memberUid")
		for _, m := range members {
			if m == username {
				return nil // Already exists
			}
		}
	}
	
	modify := ldap.NewModifyRequest(sudoGroupDN, nil)
	modify.Add("memberUid", []string{username})
	return conn.Modify(modify)
}

// removeUserFromSudoPosixGroup removes user from cn=sudo posixGroup
// If the group becomes empty, delete the entire group
func (x GroupService) removeUserFromSudoPosixGroup(conn *ldap.Conn, username string) error {
	sudoGroupDN := fmt.Sprintf("cn=sudo,%s", config.Conf.Ldap.BaseDN)
	
	// First check current members
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
	if err != nil {
		// Group doesn't exist, nothing to do
		return nil
	}
	
	if len(sr.Entries) == 0 {
		return nil
	}
	
	currentMembers := sr.Entries[0].GetAttributeValues("memberUid")
	
	// If this is the last member, delete the entire group
	if len(currentMembers) == 1 && currentMembers[0] == username {
		del := ldap.NewDelRequest(sudoGroupDN, nil)
		if delErr := conn.Del(del); delErr != nil {
			// If group doesn't exist, that's fine
			if ldapErr, ok := delErr.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
				return nil
			}
			return delErr
		}
		return nil
	}
	
	// Remove user from group
	modify := ldap.NewModifyRequest(sudoGroupDN, nil)
	modify.Delete("memberUid", []string{username})
	return conn.Modify(modify)
}

// isUserInAnySudoRole checks if user is in any sudoRole
func (x GroupService) isUserInAnySudoRole(conn *ldap.Conn, username string) bool {
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=sudoers,%s", config.Conf.Ldap.BaseDN),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=sudoRole)(sudoUser=%s))", username),
		[]string{"dn"},
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false
	}
	return len(sr.Entries) > 0
}

// extractUsername extracts username from user DN
func extractUsername(udn string) string {
	dnParts := strings.Split(udn, ",")
	if len(dnParts) > 0 {
		firstPart := strings.Split(dnParts[0], "=")
		if len(firstPart) == 2 && firstPart[0] == "uid" {
			return firstPart[1]
		}
	}
	return ""
}

// RemoveUserFromGroup removes user from group
func (x GroupService) RemoveUserFromGroup(gdn, udn string) error {
	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		gdn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"objectClass", "sudoUser"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		// This can happen if group was deleted from LDAP but still exists in MySQL
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			common.Log.Warnf("Group %s does not exist in LDAP, skipping LDAP removal", gdn)
			return nil
		}
		return fmt.Errorf("failed to query group: %v", err)
	}
	if len(sr.Entries) == 0 {
		common.Log.Warnf("Group %s not found in LDAP, skipping LDAP removal", gdn)
		return nil
	}

	objectClasses := sr.Entries[0].GetAttributeValues("objectClass")
	isPosixGroup := false
	isSudoRole := false
	for _, oc := range objectClasses {
		if oc == "posixGroup" {
			isPosixGroup = true
		}
		if oc == "sudoRole" {
			isSudoRole = true
		}
	}

	username := extractUsername(udn)
	if (isPosixGroup || isSudoRole) && username == "" {
		return errors.New("cannot extract username from DN")
	}

	newmr := ldap.NewModifyRequest(gdn, nil)
	if isSudoRole {
		currentSudoUsers := sr.Entries[0].GetAttributeValues("sudoUser")
		userExists := false
		for _, u := range currentSudoUsers {
			if u == username {
				userExists = true
				break
			}
		}
		if userExists {
			if len(currentSudoUsers) == 1 {
				del := ldap.NewDelRequest(gdn, nil)
				delErr := conn.Del(del)
				if delErr != nil {
					return fmt.Errorf("failed to delete sudoRole: %v", delErr)
				}
			} else {
				newmr.Delete("sudoUser", []string{username})
				if err := conn.Modify(newmr); err != nil {
					return err
				}
			}
			// Check if user is still in any other sudoRole, if not, remove from sudo posixGroup
			if !x.isUserInAnySudoRole(conn, username) {
				_ = x.removeUserFromSudoPosixGroup(conn, username)
			}
			return nil
		} else {
			return nil
		}
	} else if isPosixGroup {
		// Check if user is actually in the group before trying to remove
		currentMembers := sr.Entries[0].GetAttributeValues("memberUid")
		userInGroup := false
		for _, m := range currentMembers {
			if m == username {
				userInGroup = true
				break
			}
		}
		if !userInGroup {
			// User not in group, already in desired state
			return nil
		}
		newmr.Delete("memberUid", []string{username})
	} else {
		// Check if user is actually in the group before trying to remove
		currentMembers := sr.Entries[0].GetAttributeValues("uniqueMember")
		userInGroup := false
		for _, m := range currentMembers {
			if m == udn {
				userInGroup = true
				break
			}
		}
		if !userInGroup {
			return nil
		}
		if len(currentMembers) <= 1 {
			common.Log.Warnf("Cannot remove last member from groupOfUniqueNames %s, keeping existing member", gdn)
			return nil
		}
		newmr.Delete("uniqueMember", []string{udn})
	}

	if err := conn.Modify(newmr); err != nil {
		// If attribute doesn't exist, user was already removed, treat as success
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchAttribute {
			return nil
		}
		return err
	}
	return nil
}

// ListGroupDN lists all LDAP groups
func (x GroupService) ListGroupDN() (groups []*model.Group, err error) {
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(|(objectClass=organizationalUnit)(objectClass=groupOfUniqueNames)(objectClass=posixGroup)(objectClass=sudoRole))",
		[]string{"DN", "cn", "ou", "description"},
		nil,
	)

	conn, err := common.GetLDAPConn()
	if err != nil {
		return groups, err
	}
	defer common.PutLADPConn(conn)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	for _, v := range sr.Entries {
		groupName := v.GetAttributeValue("cn")
		if groupName == "" {
			groupName = v.GetAttributeValue("ou")
		}
		if groupName == "" {
			dnParts := strings.Split(v.DN, ",")
			if len(dnParts) > 0 {
				firstPart := strings.Split(dnParts[0], "=")
				if len(firstPart) == 2 {
					groupName = firstPart[1]
				}
			}
		}

		groups = append(groups, &model.Group{
			GroupDN:   v.DN,
			GroupName: groupName,
			Remark:    v.GetAttributeValue("description"),
		})
	}
	return
}

// AddPosixGroup creates posixGroup for Linux system groups
func (x GroupService) AddPosixGroup(groupName string, gidNumber uint, memberUids []string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.Conf.Ldap.BaseDN)

	searchConn, err := common.GetLDAPConn()
	if err != nil {
		return err
	}
	defer common.PutLADPConn(searchConn)

	searchRequest := ldap.NewSearchRequest(
		groupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	_, err = searchConn.Search(searchRequest)
	if err == nil {
		return x.UpdatePosixGroup(groupName, gidNumber, memberUids)
	}

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	add := ldap.NewAddRequest(groupDN, nil)
	add.Attribute("objectClass", []string{"posixGroup", "top"})
	add.Attribute("cn", []string{groupName})
	add.Attribute("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})

	if len(memberUids) > 0 {
		add.Attribute("memberUid", memberUids)
	}

	return conn.Add(add)
}

// UpdatePosixGroup updates posixGroup
func (x GroupService) UpdatePosixGroup(groupName string, gidNumber uint, memberUids []string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.Conf.Ldap.BaseDN)

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	modify := ldap.NewModifyRequest(groupDN, nil)
	modify.Replace("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})

	if len(memberUids) > 0 {
		modify.Replace("memberUid", memberUids)
	} else {
		modify.Delete("memberUid", nil)
	}

	return conn.Modify(modify)
}

// AddUserToPosixGroup adds user to posixGroup
func (x GroupService) AddUserToPosixGroup(groupName, username string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.Conf.Ldap.BaseDN)

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	modify := ldap.NewModifyRequest(groupDN, nil)
	modify.Add("memberUid", []string{username})

	return conn.Modify(modify)
}

// RemoveUserFromPosixGroup removes user from posixGroup
func (x GroupService) RemoveUserFromPosixGroup(groupName, username string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.Conf.Ldap.BaseDN)

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return err
	}
	defer conn.Close()

	modify := ldap.NewModifyRequest(groupDN, nil)
	modify.Delete("memberUid", []string{username})

	return conn.Modify(modify)
}
