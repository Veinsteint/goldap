package ildap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"
	"goldap-server/public/tools"

	ldap "github.com/go-ldap/ldap/v3"
)

type UserService struct{}

// Add creates a new user in LDAP
func (x UserService) Add(user *model.User) error {
	// Prevent creating admin user with uid=admin format
	if user.Username == "admin" || strings.HasPrefix(user.UserDN, "uid=admin,") {
		return fmt.Errorf("creating admin user is not allowed")
	}

	// Extract parent DN from UserDN (e.g., uid=user,ou=People,dc=example,dc=com -> ou=People,dc=example,dc=com)
	parentDN := extractParentDN(user.UserDN)
	if parentDN == "" {
		parentDN = config.Conf.Ldap.UserDN
	}

	// Extract OU name for creating OU if needed
	ouName := extractOUName(parentDN)
	if ouName == "" {
		if strings.Contains(config.Conf.Ldap.UserDN, "ou=") {
			ouName = extractOUName(config.Conf.Ldap.UserDN)
		}
		if ouName == "" {
			ouName = "CMPLabHPC"
		}
	}

	// Check if parent DN exists using query connection
	searchConn, err := common.GetLDAPConn()
	defer common.PutLADPConn(searchConn)
	if err != nil {
		return err
	}

	// Check and create parent DN if it doesn't exist
	searchRequest := ldap.NewSearchRequest(
		parentDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	_, err = searchConn.Search(searchRequest)
	if err != nil {
		// Parent DN doesn't exist, try to create it
		modifyConn, err := common.GetLDAPConnForModify()
		if err != nil {
			return fmt.Errorf("failed to get LDAP modify connection: %v", err)
		}
		defer modifyConn.Close()

		ouAdd := ldap.NewAddRequest(parentDN, nil)
		ouAdd.Attribute("objectClass", []string{"organizationalUnit", "top"})
		ouAdd.Attribute("ou", []string{ouName})

		err = modifyConn.Add(ouAdd)
		if err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				// Parent DN already exists, continue
			} else {
				_, checkErr := searchConn.Search(searchRequest)
				if checkErr != nil {
					return fmt.Errorf("failed to create parent DN %s: %v", parentDN, err)
				}
			}
		}
	}

	// Get LDAP connection for modification (requires admin authentication)
	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	// Create user with Unix account support (posixAccount and shadowAccount)
	add := ldap.NewAddRequest(user.UserDN, nil)

	// Base objectClass (order matters: inetOrgPerson, posixAccount, shadowAccount)
	objectClasses := []string{"inetOrgPerson"}

	// Add posixAccount and shadowAccount if user has UIDNumber
	if user.UIDNumber > 0 {
		objectClasses = append(objectClasses, "posixAccount", "shadowAccount")
	}

	add.Attribute("objectClass", objectClasses)

	// Set cn (common name) from Nickname or Username
	cnValue := user.Nickname
	if cnValue == "" {
		cnValue = user.Username
	}
	add.Attribute("cn", []string{cnValue})
	add.Attribute("sn", []string{user.Nickname})

	// businessCategory: sanitize to comply with LDAP syntax
	businessCategory := sanitizeBusinessCategory(user.Departments)
	if businessCategory != "" {
		add.Attribute("businessCategory", []string{businessCategory})
	}
	add.Attribute("departmentNumber", []string{user.Position})
	add.Attribute("description", []string{user.Introduction})
	add.Attribute("displayName", []string{user.Nickname})
	add.Attribute("mail", []string{user.Mail})

	// employeeNumber: sanitize to comply with LDAP syntax
	employeeNumber := sanitizeEmployeeNumber(user.JobNumber)
	if employeeNumber != "" {
		add.Attribute("employeeNumber", []string{employeeNumber})
	}
	add.Attribute("givenName", []string{user.GivenName})
	add.Attribute("postalAddress", []string{user.PostalAddress})

	// mobile: sanitize to comply with LDAP syntax
	mobile := sanitizeMobile(user.Mobile)
	if mobile != "" {
		add.Attribute("mobile", []string{mobile})
	}
	add.Attribute("uid", []string{user.Username})

	// Unix user attributes (posixAccount) - required attributes
	if user.UIDNumber > 0 {
		add.Attribute("uidNumber", []string{fmt.Sprintf("%d", user.UIDNumber)})

		// GIDNumber: use UIDNumber if GIDNumber is 0
		gidNumber := user.GIDNumber
		if gidNumber == 0 && user.UIDNumber > 0 {
			gidNumber = user.UIDNumber
		}
		add.Attribute("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})

		// HomeDirectory: default to /home/username
		homeDir := user.HomeDirectory
		if homeDir == "" {
			homeDir = fmt.Sprintf("/home/%s", user.Username)
		}
		add.Attribute("homeDirectory", []string{homeDir})

		// LoginShell: default to /bin/bash
		loginShell := user.LoginShell
		if loginShell == "" {
			loginShell = "/bin/bash"
		}
		add.Attribute("loginShell", []string{loginShell})

		// gecos: user info field, default to cn value
		gecosValue := user.Gecos
		if gecosValue == "" {
			gecosValue = cnValue
		}
		if gecosValue != "" {
			add.Attribute("gecos", []string{gecosValue})
		}
	}

	// Password handling
	var pass string
	if user.Password == "" {
		if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
			pass = config.Conf.Ldap.UserInitPassword
		} else {
			pass = tools.EncodePass([]byte(config.Conf.Ldap.UserInitPassword))
		}
	} else {
		// Decrypt RSA-encrypted password from MySQL
		decryptedPassword := tools.NewParPasswd(user.Password)
		if decryptedPassword == "" {
			// Use default password if decryption fails
			if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
				pass = config.Conf.Ldap.UserInitPassword
			} else {
				pass = tools.EncodePass([]byte(config.Conf.Ldap.UserInitPassword))
			}
		} else {
			if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
				pass = decryptedPassword
			} else {
				pass = tools.EncodePass([]byte(decryptedPassword))
			}
		}
	}
	add.Attribute("userPassword", []string{pass})

	// shadowAccount attributes
	if user.UIDNumber > 0 {
		// shadowLastChange: days since 1970-01-01
		shadowLastChange := fmt.Sprintf("%d", int(time.Now().Unix()/86400))
		add.Attribute("shadowLastChange", []string{shadowLastChange})
		add.Attribute("shadowMin", []string{"0"})
		add.Attribute("shadowMax", []string{"99999"})
		add.Attribute("shadowWarning", []string{"7"})
		add.Attribute("shadowInactive", []string{"-1"})
		add.Attribute("shadowExpire", []string{"-1"})
		add.Attribute("shadowFlag", []string{"0"})
	}

	err = conn.Add(add)
	if err != nil {
		return err
	}

	// Ensure corresponding posixGroup exists for user's primary group
	if user.UIDNumber > 0 {
		// Use effective GIDNumber (same as UIDNumber if not specified)
		effectiveGID := user.GIDNumber
		if effectiveGID == 0 {
			effectiveGID = user.UIDNumber
		}
		if err := x.ensurePosixGroupExists(effectiveGID, user.Username); err != nil {
			common.Log.Warnf("Failed to ensure posixGroup exists (GID: %d): %v", effectiveGID, err)
		}
	}

	return nil
}

// Update modifies an existing user in LDAP (syncs all attributes including password from MySQL)
func (x UserService) Update(oldusername string, user *model.User) error {
	// Check if user exists in LDAP
	searchConn, err := common.GetLDAPConn()
	if err != nil {
		return fmt.Errorf("failed to get LDAP connection: %v", err)
	}
	defer common.PutLADPConn(searchConn)

	searchRequest := ldap.NewSearchRequest(
		user.UserDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{"businessCategory", "employeeNumber", "mobile", "gidNumber", "uid"},
		nil,
	)

	sr, err := searchConn.Search(searchRequest)
	if err != nil {
		// User doesn't exist, create it first
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			if err := x.Add(user); err != nil {
				if ldapErr2, ok2 := err.(*ldap.Error); ok2 && ldapErr2.ResultCode == ldap.LDAPResultEntryAlreadyExists {
					// Already exists, continue with update
				} else {
					return fmt.Errorf("failed to create user %s: %v", user.UserDN, err)
				}
			} else {
				return nil // Created successfully
			}
		} else {
			return fmt.Errorf("failed to check if user %s exists: %v", user.UserDN, err)
		}
	}

	// Get current attributes to determine which ones exist
	var currentAttributes map[string][]string
	if len(sr.Entries) > 0 {
		entry := sr.Entries[0]
		currentAttributes = make(map[string][]string)
		for _, attr := range entry.Attributes {
			currentAttributes[attr.Name] = attr.Values
		}
	}

	modify := ldap.NewModifyRequest(user.UserDN, nil)

	// Set cn from Nickname or Username
	cnValue := user.Nickname
	if cnValue == "" {
		cnValue = user.Username
	}
	modify.Replace("cn", []string{cnValue})
	modify.Replace("sn", []string{user.Nickname})

	// businessCategory: sanitize and handle empty values
	businessCategory := sanitizeBusinessCategory(user.Departments)
	if businessCategory != "" {
		modify.Replace("businessCategory", []string{businessCategory})
	} else {
		if currentAttributes != nil && len(currentAttributes["businessCategory"]) > 0 {
			modify.Delete("businessCategory", nil)
		}
	}
	modify.Replace("departmentNumber", []string{user.Position})
	modify.Replace("description", []string{user.Introduction})
	modify.Replace("displayName", []string{user.Nickname})
	modify.Replace("mail", []string{user.Mail})

	// employeeNumber: sanitize and handle empty values
	employeeNumber := sanitizeEmployeeNumber(user.JobNumber)
	if employeeNumber != "" {
		modify.Replace("employeeNumber", []string{employeeNumber})
	} else {
		if currentAttributes != nil && len(currentAttributes["employeeNumber"]) > 0 {
			modify.Delete("employeeNumber", nil)
		}
	}
	modify.Replace("givenName", []string{user.GivenName})
	modify.Replace("postalAddress", []string{user.PostalAddress})

	// mobile: sanitize and handle empty values
	mobile := sanitizeMobile(user.Mobile)
	if mobile != "" {
		modify.Replace("mobile", []string{mobile})
	} else {
		if currentAttributes != nil && len(currentAttributes["mobile"]) > 0 {
			modify.Delete("mobile", nil)
		}
	}

	// Update Unix user attributes if UIDNumber exists
	if user.UIDNumber > 0 {
		modify.Replace("uidNumber", []string{fmt.Sprintf("%d", user.UIDNumber)})

		gidNumber := user.GIDNumber
		if gidNumber == 0 && user.UIDNumber > 0 {
			gidNumber = user.UIDNumber
		}
		modify.Replace("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})

		homeDir := user.HomeDirectory
		if homeDir == "" {
			homeDir = fmt.Sprintf("/home/%s", user.Username)
		}
		modify.Replace("homeDirectory", []string{homeDir})

		loginShell := user.LoginShell
		if loginShell == "" {
			loginShell = "/bin/bash"
		}
		modify.Replace("loginShell", []string{loginShell})

		gecosValue := user.Gecos
		if gecosValue == "" {
			gecosValue = cnValue
		}
		if gecosValue != "" {
			modify.Replace("gecos", []string{gecosValue})
		}

		// Update shadowAccount attributes
		shadowLastChange := fmt.Sprintf("%d", int(time.Now().Unix()/86400))
		modify.Replace("shadowLastChange", []string{shadowLastChange})
		modify.Replace("shadowMin", []string{"0"})
		modify.Replace("shadowMax", []string{"99999"})
		modify.Replace("shadowWarning", []string{"7"})
		modify.Replace("shadowInactive", []string{"-1"})
		modify.Replace("shadowExpire", []string{"-1"})
		modify.Replace("shadowFlag", []string{"0"})
	}

	// Sync password from MySQL
	if user.Password != "" {
		decryptedPassword := tools.NewParPasswd(user.Password)
		if decryptedPassword != "" {
			var pass string
			if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
				pass = decryptedPassword
			} else {
				pass = tools.EncodePass([]byte(decryptedPassword))
			}
			modify.Replace("userPassword", []string{pass})
		}
	}

	// Get LDAP connection for modification
	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	err = conn.Modify(modify)
	if err != nil {
		return err
	}

	// Ensure objectClass includes required classes for Unix attributes
	if user.UIDNumber > 0 {
		searchConn, err := common.GetLDAPConn()
		if err == nil {
			searchRequest := ldap.NewSearchRequest(
				user.UserDN,
				ldap.ScopeBaseObject,
				ldap.NeverDerefAliases,
				0,
				0,
				false,
				"(objectClass=*)",
				[]string{"objectClass", "uidNumber", "gidNumber", "homeDirectory", "loginShell"},
				nil,
			)
			sr, err := searchConn.Search(searchRequest)
			common.PutLADPConn(searchConn)
			if err == nil && len(sr.Entries) > 0 {
				entry := sr.Entries[0]
				objectClasses := entry.GetAttributeValues("objectClass")
				hasInetOrgPerson := false
				hasPosixAccount := false
				hasShadowAccount := false
				for _, oc := range objectClasses {
					if oc == "inetOrgPerson" {
						hasInetOrgPerson = true
					}
					if oc == "posixAccount" {
						hasPosixAccount = true
					}
					if oc == "shadowAccount" {
						hasShadowAccount = true
					}
				}

				// Check required posixAccount attributes
				hasUIDNumber := len(entry.GetAttributeValues("uidNumber")) > 0
				hasGIDNumber := len(entry.GetAttributeValues("gidNumber")) > 0
				hasHomeDirectory := len(entry.GetAttributeValues("homeDirectory")) > 0
				hasLoginShell := len(entry.GetAttributeValues("loginShell")) > 0

				needModifyOC := false
				modifyOC := ldap.NewModifyRequest(user.UserDN, nil)

				// Add missing objectClasses
				if !hasInetOrgPerson {
					modifyOC.Add("objectClass", []string{"inetOrgPerson"})
					needModifyOC = true
				}
				if !hasPosixAccount {
					modifyOC.Add("objectClass", []string{"posixAccount"})
					needModifyOC = true
				}
				if !hasShadowAccount {
					modifyOC.Add("objectClass", []string{"shadowAccount"})
					needModifyOC = true
				}

				// Add missing required attributes
				if !hasUIDNumber {
					modifyOC.Replace("uidNumber", []string{fmt.Sprintf("%d", user.UIDNumber)})
					needModifyOC = true
				}
				if !hasGIDNumber {
					gidNumber := user.GIDNumber
					if gidNumber == 0 {
						gidNumber = user.UIDNumber
					}
					modifyOC.Replace("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})
					needModifyOC = true
				}
				if !hasHomeDirectory {
					homeDir := user.HomeDirectory
					if homeDir == "" {
						homeDir = fmt.Sprintf("/home/%s", user.Username)
					}
					modifyOC.Replace("homeDirectory", []string{homeDir})
					needModifyOC = true
				}
				if !hasLoginShell {
					loginShell := user.LoginShell
					if loginShell == "" {
						loginShell = "/bin/bash"
					}
					modifyOC.Replace("loginShell", []string{loginShell})
					needModifyOC = true
				}

				if needModifyOC {
					if err := conn.Modify(modifyOC); err != nil {
						common.Log.Warnf("Failed to update objectClass or attributes for user %s: %v", user.UserDN, err)
					}
				}
			}
		}
	}

	if config.Conf.Ldap.UserNameModify && oldusername != user.Username {
		modifyDn := ldap.NewModifyDNRequest(fmt.Sprintf("uid=%s,%s", oldusername, config.Conf.Ldap.UserDN), fmt.Sprintf("uid=%s", user.Username), true, "")
		err = conn.ModifyDN(modifyDn)
		if err != nil {
			return err
		}
		user.UserDN = fmt.Sprintf("uid=%s,%s", user.Username, extractParentDN(user.UserDN))
	}

	// Handle posixGroup for user's primary group
	if user.UIDNumber > 0 {
		effectiveGID := user.GIDNumber
		if effectiveGID == 0 {
			effectiveGID = user.UIDNumber
		}

		// Check if GID has changed - update posixGroup accordingly
		var oldGID uint
		if currentAttributes != nil {
			if gidStrs := currentAttributes["gidNumber"]; len(gidStrs) > 0 {
				if gidVal, err := strconv.ParseUint(gidStrs[0], 10, 32); err == nil {
					oldGID = uint(gidVal)
				}
			}
		}

		// Get old username for personal group management
		oldUsername := oldusername
		if oldUsername == "" {
			oldUsername = user.Username
		}

		// If GID changed or username changed, handle posixGroup update
		if oldGID > 0 && (oldGID != effectiveGID || oldUsername != user.Username) {
			// Delete old personal posixGroup
			oldGroupDN := fmt.Sprintf("cn=%s,%s", oldUsername, config.Conf.Ldap.BaseDN)
			delReq := ldap.NewDelRequest(oldGroupDN, nil)
			_ = conn.Del(delReq) // Ignore errors
		}

		// Ensure new posixGroup exists
		if err := x.ensurePosixGroupExists(effectiveGID, user.Username); err != nil {
			common.Log.Warnf("Failed to ensure posixGroup exists (GID: %d): %v", effectiveGID, err)
		}
	}

	return nil
}

// Exist checks if a user exists in LDAP
func (x UserService) Exist(filter map[string]any) (bool, error) {
	filter_str := ""
	for key, value := range filter {
		filter_str += fmt.Sprintf("(%s=%s)", key, value)
	}
	search_filter := fmt.Sprintf("(&(|(objectClass=inetOrgPerson)(objectClass=simpleSecurityObject))%s)", filter_str)

	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		search_filter,
		[]string{"DN"},
		nil,
	)

	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return false, err
	}

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, err
	}
	if len(sr.Entries) > 0 {
		return true, nil
	}
	return false, nil
}

// GetUserDetails fetches user attributes from LDAP for comparison
func (x UserService) GetUserDetails(userDN string) (map[string]string, error) {
	conn, err := common.GetLDAPConn()
	if err != nil {
		return nil, err
	}
	defer common.PutLADPConn(conn)

	searchRequest := ldap.NewSearchRequest(
		userDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"cn", "sn", "mail", "mobile", "givenName", "displayName",
			"businessCategory", "departmentNumber", "description", "employeeNumber",
			"postalAddress", "uidNumber", "gidNumber", "homeDirectory", "loginShell", "gecos"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("user not found: %s", userDN)
	}

	entry := sr.Entries[0]
	attrs := make(map[string]string)
	for _, attr := range entry.Attributes {
		if len(attr.Values) > 0 {
			attrs[attr.Name] = attr.Values[0]
		}
	}
	return attrs, nil
}

// NeedsUpdate compares MySQL user with LDAP attributes, returns true if update needed
func (x UserService) NeedsUpdate(user *model.User, ldapAttrs map[string]string) bool {
	// Check password first (if MySQL has password)
	if user.Password != "" {
		decryptedPassword := tools.NewParPasswd(user.Password)
		if decryptedPassword != "" {
			// Try to verify password by binding
			if !x.VerifyPassword(user.UserDN, decryptedPassword) {
				return true // Password mismatch, needs update
			}
		}
	}

	// Compare key attributes
	cnValue := user.Nickname
	if cnValue == "" {
		cnValue = user.Username
	}

	// cn (common name)
	if ldapAttrs["cn"] != cnValue {
		return true
	}

	// sn (surname)
	if ldapAttrs["sn"] != user.Nickname {
		return true
	}

	// mail
	if ldapAttrs["mail"] != user.Mail {
		return true
	}

	// mobile (sanitized)
	sanitizedMobile := sanitizeMobile(user.Mobile)
	if ldapAttrs["mobile"] != sanitizedMobile && !(ldapAttrs["mobile"] == "" && sanitizedMobile == "") {
		return true
	}

	// givenName
	if ldapAttrs["givenName"] != user.GivenName {
		return true
	}

	// displayName
	if ldapAttrs["displayName"] != user.Nickname {
		return true
	}

	// businessCategory (departments)
	sanitizedDepts := sanitizeBusinessCategory(user.Departments)
	if ldapAttrs["businessCategory"] != sanitizedDepts && !(ldapAttrs["businessCategory"] == "" && sanitizedDepts == "") {
		return true
	}

	// departmentNumber (position)
	if ldapAttrs["departmentNumber"] != user.Position {
		return true
	}

	// description (introduction)
	if ldapAttrs["description"] != user.Introduction {
		return true
	}

	// employeeNumber (jobNumber)
	sanitizedJobNum := sanitizeEmployeeNumber(user.JobNumber)
	if ldapAttrs["employeeNumber"] != sanitizedJobNum && !(ldapAttrs["employeeNumber"] == "" && sanitizedJobNum == "") {
		return true
	}

	// postalAddress
	if ldapAttrs["postalAddress"] != user.PostalAddress {
		return true
	}

	// Unix attributes (if user has UIDNumber)
	if user.UIDNumber > 0 {
		uidStr := fmt.Sprintf("%d", user.UIDNumber)
		if ldapAttrs["uidNumber"] != uidStr {
			return true
		}

		gidNumber := user.GIDNumber
		if gidNumber == 0 {
			gidNumber = user.UIDNumber
		}
		gidStr := fmt.Sprintf("%d", gidNumber)
		if ldapAttrs["gidNumber"] != gidStr {
			return true
		}

		homeDir := user.HomeDirectory
		if homeDir == "" {
			homeDir = fmt.Sprintf("/home/%s", user.Username)
		}
		if ldapAttrs["homeDirectory"] != homeDir {
			return true
		}

		loginShell := user.LoginShell
		if loginShell == "" {
			loginShell = "/bin/bash"
		}
		if ldapAttrs["loginShell"] != loginShell {
			return true
		}

		gecosValue := user.Gecos
		if gecosValue == "" {
			gecosValue = cnValue
		}
		if ldapAttrs["gecos"] != gecosValue && gecosValue != "" {
			return true
		}
	}

	return false
}

// VerifyPassword checks if the given password matches the LDAP user's password
func (x UserService) VerifyPassword(userDN, password string) bool {
	if userDN == "" || password == "" {
		return false
	}

	// Try to bind with the user's credentials
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", config.Conf.Ldap.Url))
	if err != nil {
		return false
	}
	defer conn.Close()

	// Attempt to bind with the user's DN and password
	err = conn.Bind(userDN, password)
	return err == nil
}

// Delete removes a user from LDAP
func (x UserService) Delete(udn string) error {
	// Extract username from user DN
	username := ""
	dnParts := strings.Split(udn, ",")
	if len(dnParts) > 0 {
		firstPart := strings.Split(dnParts[0], "=")
		if len(firstPart) == 2 && firstPart[0] == "uid" {
			username = firstPart[1]
		}
	}

	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	// Delete user entry
	del := ldap.NewDelRequest(udn, nil)
	err = conn.Del(del)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			// User doesn't exist, ignore
		} else {
			return fmt.Errorf("failed to delete user entry: %v", err)
		}
	}

	// Clean up user from all groups
	if username != "" {
		// Delete user's personal posixGroup
		posixGroupDN := fmt.Sprintf("cn=%s,%s", username, config.Conf.Ldap.BaseDN)
		posixGroupDel := ldap.NewDelRequest(posixGroupDN, nil)
		_ = conn.Del(posixGroupDel) // Ignore errors - posixGroup may not exist

		// Remove user from all posixGroups (docker, sudo, etc.)
		_ = x.removeUserFromAllPosixGroups(conn, username)

		// Remove user from all sudoRoles
		_ = x.removeUserFromAllSudoRoles(conn, username)

		// Remove user from all groupOfUniqueNames
		_ = x.removeUserFromAllGroupOfUniqueNames(conn, udn)
	}

	return nil
}

// removeUserFromAllSudoRoles removes user from all sudoRole entries
func (x UserService) removeUserFromAllSudoRoles(conn *ldap.Conn, username string) error {
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=sudoers,%s", config.Conf.Ldap.BaseDN),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=sudoRole)(sudoUser=%s))", escapeLDAPFilter(username)),
		[]string{"dn", "sudoUser"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			return nil
		}
		return fmt.Errorf("failed to search sudoRoles: %v", err)
	}

	for _, entry := range sr.Entries {
		sudoRoleDN := entry.DN
		currentSudoUsers := entry.GetAttributeValues("sudoUser")

		userExists := false
		for _, u := range currentSudoUsers {
			if u == username {
				userExists = true
				break
			}
		}

		if !userExists {
			continue
		}

		// If only one user, delete the entire sudoRole (sudoUser is required)
		if len(currentSudoUsers) == 1 {
			del := ldap.NewDelRequest(sudoRoleDN, nil)
			_ = conn.Del(del)
		} else {
			modify := ldap.NewModifyRequest(sudoRoleDN, nil)
			modify.Delete("sudoUser", []string{username})
			_ = conn.Modify(modify)
		}
	}

	return nil
}

// removeUserFromAllPosixGroups removes user from all posixGroup entries (excluding user's own group)
func (x UserService) removeUserFromAllPosixGroups(conn *ldap.Conn, username string) error {
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%s))", escapeLDAPFilter(username)),
		[]string{"dn", "cn", "memberUid"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			return nil
		}
		return fmt.Errorf("failed to search posixGroups: %v", err)
	}

	for _, entry := range sr.Entries {
		groupDN := entry.DN
		groupCN := entry.GetAttributeValue("cn")

		// Skip user's own personal group (will be deleted separately)
		if groupCN == username {
			continue
		}

		modify := ldap.NewModifyRequest(groupDN, nil)
		modify.Delete("memberUid", []string{username})
		if err := conn.Modify(modify); err != nil {
			common.Log.Warnf("Failed to remove user %s from posixGroup %s: %v", username, groupDN, err)
		}
	}

	return nil
}

// removeUserFromAllGroupOfUniqueNames removes user DN from all groupOfUniqueNames entries
func (x UserService) removeUserFromAllGroupOfUniqueNames(conn *ldap.Conn, userDN string) error {
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=groupOfUniqueNames)(uniqueMember=%s))", escapeLDAPFilter(userDN)),
		[]string{"dn", "uniqueMember"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			return nil
		}
		return fmt.Errorf("failed to search groupOfUniqueNames: %v", err)
	}

	for _, entry := range sr.Entries {
		groupDN := entry.DN
		members := entry.GetAttributeValues("uniqueMember")

		// Don't remove last member (uniqueMember is required)
		if len(members) <= 1 {
			continue
		}

		modify := ldap.NewModifyRequest(groupDN, nil)
		modify.Delete("uniqueMember", []string{userDN})
		if err := conn.Modify(modify); err != nil {
			common.Log.Warnf("Failed to remove user %s from groupOfUniqueNames %s: %v", userDN, groupDN, err)
		}
	}

	return nil
}

// ChangePwd changes user password
// Uses SSHA format for Linux system compatibility
func (x UserService) ChangePwd(udn, oldpasswd, newpasswd string) error {
	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	// Try PasswordModify extension first (standard method)
	modifyPass := ldap.NewPasswordModifyRequest(udn, oldpasswd, newpasswd)
	_, err = conn.PasswordModify(modifyPass)
	if err == nil {
		return nil
	}

	// Fallback to manual SSHA password setting for Linux compatibility
	var pass string
	if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
		pass = newpasswd
	} else {
		pass = tools.EncodePass([]byte(newpasswd))
	}

	modify := ldap.NewModifyRequest(udn, nil)
	modify.Replace("userPassword", []string{pass})

	err = conn.Modify(modify)
	if err != nil {
		if udn == config.Conf.Ldap.AdminDN {
			fallbackDN := fmt.Sprintf("uid=admin,%s", config.Conf.Ldap.UserDN)
			modifyFallback := ldap.NewModifyRequest(fallbackDN, nil)
			modifyFallback.Replace("userPassword", []string{pass})
			if errFallback := conn.Modify(modifyFallback); errFallback == nil {
				return nil
			}
			return fmt.Errorf("password modify failed for %s: %v", udn, err)
		}
		return fmt.Errorf("password modify failed for %s: %v", udn, err)
	}

	return nil
}

// NewPwd generates a new password for user
// If userDN is provided, it takes precedence; otherwise DN is constructed from username
func (x UserService) NewPwd(username string, userDN ...string) (string, error) {
	var udn string
	if len(userDN) > 0 && userDN[0] != "" {
		udn = userDN[0]
	} else {
		udn = fmt.Sprintf("uid=%s,%s", username, config.Conf.Ldap.UserDN)
		if username == "admin" {
			udn = config.Conf.Ldap.AdminDN
		}
	}
	modifyPass := ldap.NewPasswordModifyRequest(udn, "", "")

	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return "", err
	}

	newpass, err := conn.PasswordModify(modifyPass)
	if err != nil {
		if (udn == config.Conf.Ldap.AdminDN || username == "admin") && udn != fmt.Sprintf("uid=admin,%s", config.Conf.Ldap.UserDN) {
			fallbackDN := fmt.Sprintf("uid=admin,%s", config.Conf.Ldap.UserDN)
			modifyPassFallback := ldap.NewPasswordModifyRequest(fallbackDN, "", "")
			newpassFallback, errFallback := conn.PasswordModify(modifyPassFallback)
			if errFallback == nil {
				return newpassFallback.GeneratedPassword, nil
			}
			return "", fmt.Errorf("password modify failed for %s: %v", username, errFallback)
		}
		return "", fmt.Errorf("password modify failed for %s: %v", username, err)
	}
	return newpass.GeneratedPassword, nil
}

// UpdateSSHKeys updates user's SSH public keys in LDAP
func (x UserService) UpdateSSHKeys(userDN string, sshKeys []string) error {
	conn, err := common.GetLDAPConnForModify()
	defer conn.Close()
	if err != nil {
		return err
	}

	modify := ldap.NewModifyRequest(userDN, nil)

	if len(sshKeys) == 0 {
		modify.Delete("sshPublicKey", []string{})
	} else {
		modify.Replace("sshPublicKey", sshKeys)
	}

	return conn.Modify(modify)
}

// ListUserDN returns all user DNs from LDAP
func (x UserService) ListUserDN() (users []*model.User, err error) {
	var conn *ldap.Conn
	var useAnonConn bool

	if config.Conf.Ldap.AllowAnonBinding {
		conn, err = common.GetLDAPConn()
		if err == nil {
			useAnonConn = true
			defer common.PutLADPConn(conn)
		}
	}

	if !useAnonConn {
		conn, err = common.GetLDAPConnForModify()
		if err != nil {
			return []*model.User{}, nil
		}
		defer conn.Close()
	}

	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(|(objectClass=inetOrgPerson)(objectClass=simpleSecurityObject))",
		[]string{"DN"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		if useAnonConn {
			// Fallback to admin connection
			common.PutLADPConn(conn)
			adminConn, adminErr := common.GetLDAPConnForModify()
			if adminErr != nil {
				return []*model.User{}, nil
			}
			defer adminConn.Close()
			sr, err = adminConn.Search(searchRequest)
			if err != nil {
				return []*model.User{}, nil
			}
		} else {
			return []*model.User{}, nil
		}
	}
	if len(sr.Entries) > 0 {
		for _, v := range sr.Entries {
			users = append(users, &model.User{
				UserDN: v.DN,
			})
		}
	}
	return
}

// ensurePosixGroupExists ensures a posixGroup exists for the given GID
func (x UserService) ensurePosixGroupExists(gidNumber uint, username string) error {
	groupName := username
	groupDN := fmt.Sprintf("cn=%s,%s", groupName, config.Conf.Ldap.BaseDN)

	searchConn, err := common.GetLDAPConn()
	if err != nil {
		return fmt.Errorf("failed to get LDAP connection: %v", err)
	}
	defer common.PutLADPConn(searchConn)

	searchRequest := ldap.NewSearchRequest(
		groupDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=posixGroup)",
		[]string{"dn", "gidNumber", "memberUid"},
		nil,
	)

	sr, err := searchConn.Search(searchRequest)
	if err == nil && len(sr.Entries) > 0 {
		// Group exists, check GID and membership
		entry := sr.Entries[0]

		// Check if GID needs to be updated
		currentGIDStrs := entry.GetAttributeValues("gidNumber")
		if len(currentGIDStrs) > 0 {
			if currentGID, parseErr := strconv.ParseUint(currentGIDStrs[0], 10, 32); parseErr == nil {
				if uint(currentGID) != gidNumber {
					// GID changed, update it
					modifyConn, modErr := common.GetLDAPConnForModify()
					if modErr == nil {
						defer modifyConn.Close()
						modify := ldap.NewModifyRequest(groupDN, nil)
						modify.Replace("gidNumber", []string{fmt.Sprintf("%d", gidNumber)})
						_ = modifyConn.Modify(modify)
					}
				}
			}
		}

		// Check if user is a member
		memberUids := entry.GetAttributeValues("memberUid")
		userInGroup := false
		for _, uid := range memberUids {
			if uid == username {
				userInGroup = true
				break
			}
		}

		if !userInGroup {
			_ = Group.AddUserToPosixGroup(groupName, username)
		}
		return nil
	}

	// Group doesn't exist, create it
	if err := Group.AddPosixGroup(groupName, gidNumber, []string{username}); err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
			// Already exists, try to add user
			_ = Group.AddUserToPosixGroup(groupName, username)
			return nil
		}
		return fmt.Errorf("failed to create posixGroup [%s]: %v", groupName, err)
	}

	return nil
}

// sanitizeMobile sanitizes mobile value for LDAP syntax compliance
func sanitizeMobile(mobile string) string {
	if mobile == "" {
		return ""
	}

	mobile = strings.TrimSpace(mobile)
	if mobile == "" {
		return ""
	}

	// Allow: digits, plus sign, hyphen, space, parentheses, period
	var result strings.Builder
	for _, r := range mobile {
		if (r >= '0' && r <= '9') ||
			r == '+' || r == '-' || r == ' ' ||
			r == '(' || r == ')' || r == '.' {
			result.WriteRune(r)
		}
	}

	cleaned := strings.TrimSpace(result.String())

	// Limit length to 20 characters
	if len(cleaned) > 20 {
		cleaned = cleaned[:20]
	}

	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

// sanitizeEmployeeNumber sanitizes employeeNumber value for LDAP syntax compliance
func sanitizeEmployeeNumber(jobNumber string) string {
	if jobNumber == "" {
		return ""
	}

	jobNumber = strings.TrimSpace(jobNumber)
	if jobNumber == "" {
		return ""
	}

	// Allow: letters, digits, hyphen, underscore
	var result strings.Builder
	for _, r := range jobNumber {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' {
			result.WriteRune(r)
		} else if r == ' ' {
			result.WriteRune('_')
		}
	}

	cleaned := strings.TrimSpace(result.String())

	// Limit length to 64 characters
	if len(cleaned) > 64 {
		cleaned = cleaned[:64]
	}

	return cleaned
}

// sanitizeBusinessCategory sanitizes businessCategory value for LDAP syntax compliance
func sanitizeBusinessCategory(departments string) string {
	if departments == "" {
		return ""
	}

	departments = strings.TrimSpace(departments)
	if departments == "" {
		return ""
	}

	// Take only the first department (businessCategory doesn't support multiple values)
	parts := strings.Split(departments, ",")
	firstDept := strings.TrimSpace(parts[0])
	if firstDept == "" {
		return ""
	}

	// Allow: letters, digits, space, hyphen, underscore, period
	var result strings.Builder
	for _, r := range firstDept {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == ' ' || r == '-' || r == '_' || r == '.' {
			result.WriteRune(r)
		} else {
			if r < 128 {
				result.WriteRune('_')
			}
		}
	}

	cleaned := strings.TrimSpace(result.String())

	// Limit length to 64 characters
	if len(cleaned) > 64 {
		cleaned = cleaned[:64]
	}

	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

// escapeLDAPFilter escapes special characters in LDAP filter values
func escapeLDAPFilter(value string) string {
	result := strings.Builder{}
	for _, r := range value {
		switch r {
		case '*':
			result.WriteString("\\2a")
		case '(':
			result.WriteString("\\28")
		case ')':
			result.WriteString("\\29")
		case '\\':
			result.WriteString("\\5c")
		case '/':
			result.WriteString("\\2f")
		case '\x00':
			result.WriteString("\\00")
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}
