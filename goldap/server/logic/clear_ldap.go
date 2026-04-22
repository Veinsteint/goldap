package logic

import (
	"fmt"
	"os/exec"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"
	"goldap-server/service/ildap"

	ldap "github.com/go-ldap/ldap/v3"
)

// deleteLDAPEntryWithCommand deletes LDAP entry using ldapdelete command (fallback)
func deleteLDAPEntryWithCommand(dn string) error {
	cmd := exec.Command("ldapdelete",
		"-x",
		"-H", config.Conf.Ldap.Url,
		"-D", config.Conf.Ldap.AdminDN,
		"-w", config.Conf.Ldap.AdminPass,
		dn,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 32 {
				return nil // Entry doesn't exist
			}
		}
		return fmt.Errorf("ldapdelete failed: %v, output: %s", err, string(output))
	}
	return nil
}

// modifyLDAPEntryWithCommand modifies LDAP entry using ldapmodify command (fallback)
func modifyLDAPEntryWithCommand(ldifContent string) error {
	cmd := exec.Command("ldapmodify",
		"-x",
		"-H", config.Conf.Ldap.Url,
		"-D", config.Conf.Ldap.AdminDN,
		"-w", config.Conf.Ldap.AdminPass,
	)

	cmd.Stdin = strings.NewReader(ldifContent)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ldapmodify failed: %v, output: %s", err, string(output))
	}
	return nil
}

// deleteLDAPEntry deletes LDAP entry, falls back to ldapdelete command on failure
func deleteLDAPEntry(dn string, entryType string) error {
	var err error
	if entryType == "user" {
		err = ildap.User.Delete(dn)
	} else if entryType == "group" {
		err = ildap.Group.Delete(dn)
	} else {
		del := ldap.NewDelRequest(dn, nil)
		conn, connErr := common.GetLDAPConnForModify()
		if connErr != nil {
			return fmt.Errorf("failed to get LDAP connection: %v", connErr)
		}
		defer conn.Close()
		err = conn.Del(del)
	}

	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			return nil
		}

		if cmdErr := deleteLDAPEntryWithCommand(dn); cmdErr != nil {
			return fmt.Errorf("both Go library and ldapdelete failed: %v, %v", err, cmdErr)
		}
		return nil
	}

	return nil
}

// ClearLDAP clears all users and groups from LDAP (preserving admin and baseDN)
func ClearLDAP() error {

	conn, err := common.GetLDAPConnForModify()
	if err != nil {
		return fmt.Errorf("failed to get LDAP connection: %v", err)
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	_, err = conn.Search(searchRequest)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			common.Log.Warnf("BaseDN %s doesn't exist, continuing...", config.Conf.Ldap.BaseDN)
		} else {
			return fmt.Errorf("failed to check baseDN: %v", err)
		}
	}

	// Delete users
	users, _ := ildap.User.ListUserDN()
	searchConn, searchErr := common.GetLDAPConnForModify()
	if searchErr == nil {
		defer searchConn.Close()

		ouSearchReq := ldap.NewSearchRequest(
			config.Conf.Ldap.BaseDN,
			ldap.ScopeSingleLevel,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=organizationalUnit)",
			[]string{"DN"},
			nil,
		)

		ouSr, _ := searchConn.Search(ouSearchReq)
		if len(ouSr.Entries) > 0 {
			for _, ouEntry := range ouSr.Entries {
				if ouEntry.DN == config.Conf.Ldap.BaseDN {
					continue
				}

				userSearchReq := ldap.NewSearchRequest(
					ouEntry.DN,
					ldap.ScopeSingleLevel,
					ldap.NeverDerefAliases,
					0, 0, false,
					"(|(objectClass=inetOrgPerson)(objectClass=simpleSecurityObject))",
					[]string{"DN"},
					nil,
				)

				userSr, _ := searchConn.Search(userSearchReq)
				for _, entry := range userSr.Entries {
					exists := false
					for _, u := range users {
						if u.UserDN == entry.DN {
							exists = true
							break
						}
					}
					if !exists {
						users = append(users, &model.User{UserDN: entry.DN})
					}
				}
			}
		}
	}

	deletedUserCount := 0
	for _, user := range users {
		if user.UserDN == config.Conf.Ldap.AdminDN {
			continue
		}
		if err := deleteLDAPEntry(user.UserDN, "user"); err == nil {
			deletedUserCount++
		}
	}

	// Delete groups
	groups, _ := ildap.Group.ListGroupDN()

	groupMap := make(map[string]*struct {
		DN       string
		Children []string
	})

	for _, group := range groups {
		groupMap[group.GroupDN] = &struct {
			DN       string
			Children []string
		}{DN: group.GroupDN, Children: []string{}}
	}

	for _, group := range groups {
		parts := strings.Split(group.GroupDN, ",")
		if len(parts) > 1 {
			parentDN := strings.Join(parts[1:], ",")
			if parent, exists := groupMap[parentDN]; exists {
				parent.Children = append(parent.Children, group.GroupDN)
			}
		}
	}

	deletedGroupCount := 0
	maxIterations := len(groups) + 10
	for iteration := 0; iteration < maxIterations && len(groupMap) > 0; iteration++ {
		leafNodes := []string{}
		for dn, group := range groupMap {
			if len(group.Children) == 0 {
				leafNodes = append(leafNodes, dn)
			}
		}

		if len(leafNodes) == 0 {
			for dn := range groupMap {
				leafNodes = append(leafNodes, dn)
			}
		}

		for _, dn := range leafNodes {
			if dn == config.Conf.Ldap.BaseDN {
				delete(groupMap, dn)
				continue
			}

			if deleteLDAPEntry(dn, "group") == nil {
				deletedGroupCount++
			}

			delete(groupMap, dn)

			parts := strings.Split(dn, ",")
			if len(parts) > 1 {
				parentDN := strings.Join(parts[1:], ",")
				if parent, exists := groupMap[parentDN]; exists {
					newChildren := []string{}
					for _, child := range parent.Children {
						if child != dn {
							newChildren = append(newChildren, child)
						}
					}
					parent.Children = newChildren
				}
			}
		}
	}


	// Delete sudo rules
	sudoRules, err := ildap.Sudo.List()
	if err == nil {
		deletedSudoCount := 0
		for _, rule := range sudoRules {
			sudoDN := fmt.Sprintf("cn=%s,ou=sudoers,%s", rule.Name, config.Conf.Ldap.BaseDN)
			if deleteLDAPEntry(sudoDN, "sudo") == nil {
				deletedSudoCount++
			}
		}
	}

	return nil
}
