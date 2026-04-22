package ildap

import (
	"fmt"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"

	ldap "github.com/go-ldap/ldap/v3"
)

type SudoService struct{}

// Add creates sudo rule in LDAP
func (s SudoService) Add(rule *model.SudoRule) error {
	sudoDN := fmt.Sprintf("cn=%s,ou=sudoers,%s", rule.Name, config.Conf.Ldap.BaseDN)

	// Ensure ou=sudoers exists
	sudoersDN := fmt.Sprintf("ou=sudoers,%s", config.Conf.Ldap.BaseDN)
	searchRequest := ldap.NewSearchRequest(
		sudoersDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	searchConn, err := common.GetLDAPConn()
	defer common.PutLADPConn(searchConn)
	if err != nil {
		return err
	}

	_, err = searchConn.Search(searchRequest)
	if err != nil {
		modifyConn, err := common.GetLDAPConnForModify()
		if err != nil {
			return fmt.Errorf("failed to get LDAP modify connection: %v", err)
		}
		defer modifyConn.Close()
		
		ouAdd := ldap.NewAddRequest(sudoersDN, nil)
		ouAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
		ouAdd.Attribute("ou", []string{"sudoers"})
		err = modifyConn.Add(ouAdd)
		if err != nil {
			if ldapErr, ok := err.(*ldap.Error); !ok || ldapErr.ResultCode != ldap.LDAPResultEntryAlreadyExists {
				return fmt.Errorf("failed to create ou=sudoers: %v", err)
			}
		}
	}

	// Build attributes with defaults
	sudoUser := rule.User
	sudoHost := rule.Host
	if sudoHost == "" {
		sudoHost = "ALL"
	}
	sudoCommand := rule.Command
	if sudoCommand == "" {
		sudoCommand = "ALL"
	}
	sudoRunAsUser := rule.RunAsUser
	if sudoRunAsUser == "" {
		sudoRunAsUser = "ALL"
	}
	sudoRunAsGroup := rule.RunAsGroup
	if sudoRunAsGroup == "" {
		sudoRunAsGroup = "ALL"
	}
	sudoOption := rule.Options

	// Create sudo rule entry
	add := ldap.NewAddRequest(sudoDN, nil)
	add.Attribute("objectClass", []string{"top", "sudoRole"})
	add.Attribute("cn", []string{rule.Name})
	if rule.Description != "" {
		add.Attribute("description", []string{rule.Description})
	}
	
	if sudoUser == "" {
		return fmt.Errorf("sudoUser is required")
	}
	add.Attribute("sudoUser", []string{sudoUser})
	add.Attribute("sudoHost", []string{sudoHost})
	add.Attribute("sudoCommand", []string{sudoCommand})
	add.Attribute("sudoRunAsUser", []string{sudoRunAsUser})
	if sudoRunAsGroup != "ALL" {
		add.Attribute("sudoRunAsGroup", []string{sudoRunAsGroup})
	}
	if sudoOption != "" {
		add.Attribute("sudoOption", []string{sudoOption})
	}

	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	return conn.Add(add)
}

// Update modifies sudo rule
func (s SudoService) Update(rule *model.SudoRule) error {
	sudoDN := fmt.Sprintf("cn=%s,ou=sudoers,%s", rule.Name, config.Conf.Ldap.BaseDN)
	modify := ldap.NewModifyRequest(sudoDN, nil)

	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	// Build attributes
	sudoUser := rule.User
	sudoHost := rule.Host
	if sudoHost == "" {
		sudoHost = "ALL"
	}
	sudoCommand := rule.Command
	if sudoCommand == "" {
		sudoCommand = "ALL"
	}
	sudoRunAsUser := rule.RunAsUser
	if sudoRunAsUser == "" {
		sudoRunAsUser = "ALL"
	}
	sudoRunAsGroup := rule.RunAsGroup
	if sudoRunAsGroup == "" {
		sudoRunAsGroup = "ALL"
	}
	sudoOption := rule.Options

	if rule.Description != "" {
		modify.Replace("description", []string{rule.Description})
	}
	if sudoUser != "" {
		modify.Replace("sudoUser", []string{sudoUser})
	}
	modify.Replace("sudoHost", []string{sudoHost})
	modify.Replace("sudoCommand", []string{sudoCommand})
	modify.Replace("sudoRunAsUser", []string{sudoRunAsUser})
	if sudoRunAsGroup != "ALL" {
		modify.Replace("sudoRunAsGroup", []string{sudoRunAsGroup})
	}
	if sudoOption != "" {
		modify.Replace("sudoOption", []string{sudoOption})
	} else {
		modify.Delete("sudoOption", []string{})
	}

	return conn.Modify(modify)
}

// Delete removes sudo rule
func (s SudoService) Delete(ruleName string) error {
	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	sudoDN := fmt.Sprintf("cn=%s,ou=sudoers,%s", ruleName, config.Conf.Ldap.BaseDN)
	del := ldap.NewDelRequest(sudoDN, nil)
	return conn.Del(del)
}

// List returns all sudo rules
func (s SudoService) List() ([]*model.SudoRule, error) {
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return nil, err
	}

	sudoersDN := fmt.Sprintf("ou=sudoers,%s", config.Conf.Ldap.BaseDN)
	searchRequest := ldap.NewSearchRequest(
		sudoersDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=sudoRole)",
		[]string{"*"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var rules []*model.SudoRule
	for _, entry := range sr.Entries {
		rule := &model.SudoRule{
			Name:        entry.GetAttributeValue("cn"),
			Description: entry.GetAttributeValue("description"),
			User:        entry.GetAttributeValue("sudoUser"),
			Host:        entry.GetAttributeValue("sudoHost"),
			Command:     entry.GetAttributeValue("sudoCommand"),
			RunAsUser:   entry.GetAttributeValue("sudoRunAsUser"),
			RunAsGroup:  entry.GetAttributeValue("sudoRunAsGroup"),
			Options:     entry.GetAttributeValue("sudoOption"),
		}
		// Set defaults for empty values
		if rule.User == "" {
			rule.User = "ALL"
		}
		if rule.Host == "" {
			rule.Host = "ALL"
		}
		if rule.Command == "" {
			rule.Command = "ALL"
		}
		if rule.RunAsUser == "" {
			rule.RunAsUser = "ALL"
		}
		if rule.RunAsGroup == "" {
			rule.RunAsGroup = "ALL"
		}
		rules = append(rules, rule)
	}

	return rules, nil
}
