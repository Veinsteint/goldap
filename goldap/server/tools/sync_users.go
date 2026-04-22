// User synchronization tool between MySQL and LDAP
package main

import (
	"fmt"
	"os"
	"strings"

	"goldap-server/config"
	"goldap-server/model"
	"goldap-server/public/common"
	"goldap-server/public/tools"
	"goldap-server/service/ildap"
	"goldap-server/service/isql"

	ldap "github.com/go-ldap/ldap/v3"
)

func main() {
	config.InitConfig()
	common.InitLogger()
	common.InitDB()
	common.InitLDAP()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "mysql-to-ldap":
		if len(os.Args) < 3 {
			fmt.Println("Error: user ID required")
			os.Exit(1)
		}
		syncMySQLToLDAP(os.Args[2])
	case "ldap-to-mysql":
		if len(os.Args) < 3 {
			fmt.Println("Error: username required")
			os.Exit(1)
		}
		syncLDAPToMySQL(os.Args[2])
	case "auto":
		autoSync()
	default:
		fmt.Printf("Error: unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  sync_users mysql-to-ldap <user_id>  - Sync MySQL user to LDAP")
	fmt.Println("  sync_users ldap-to-mysql <username> - Sync LDAP user to MySQL")
	fmt.Println("  sync_users auto                     - Auto sync all mismatched users")
}

// syncMySQLToLDAP syncs a MySQL user to LDAP
func syncMySQLToLDAP(userIDStr string) {
	fmt.Printf("=== Sync MySQL user (ID: %s) to LDAP ===\n\n", userIDStr)

	var userID uint
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		fmt.Printf("Error: invalid user ID: %s\n", userIDStr)
		os.Exit(1)
	}

	var user model.User
	if err := isql.User.Find(tools.H{"id": userID}, &user); err != nil {
		fmt.Printf("Error: user not found in MySQL (ID: %d)\n", userID)
		os.Exit(1)
	}

	fmt.Printf("Found user: %s (ID: %d, DN: %s)\n\n", user.Username, user.ID, user.UserDN)

	exists, _ := ildap.User.Exist(tools.H{"dn": user.UserDN})
	if exists {
		if err := ildap.User.Update(user.Username, &user); err != nil {
			fmt.Printf("Error: failed to update LDAP user: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OK: LDAP user updated")
	} else {
		if err := ildap.User.Add(&user); err != nil {
			fmt.Printf("Error: failed to create LDAP user: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("OK: LDAP user created")
	}

	if err := isql.User.ChangeSyncState(int(user.ID), 1); err != nil {
		fmt.Printf("Warning: failed to update sync state: %v\n", err)
	}
	fmt.Println("Done")
}

// syncLDAPToMySQL syncs an LDAP user to MySQL
func syncLDAPToMySQL(username string) {
	fmt.Printf("=== Sync LDAP user (%s) to MySQL ===\n\n", username)

	userDN := fmt.Sprintf("uid=%s,%s", username, config.Conf.Ldap.UserDN)

	conn, err := common.GetLDAPConn()
	if err != nil {
		fmt.Printf("Error: failed to get LDAP connection: %v\n", err)
		os.Exit(1)
	}
	defer common.PutLADPConn(conn)

	sr, err := conn.Search(ldap.NewSearchRequest(
		userDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)", []string{"*"}, nil,
	))
	if err != nil || len(sr.Entries) == 0 {
		fmt.Printf("Error: user not found in LDAP: %s\n", username)
		os.Exit(1)
	}

	entry := sr.Entries[0]
	fmt.Printf("Found LDAP user: %s\n", entry.DN)

	if isql.User.Exist(tools.H{"user_dn": userDN}) {
		fmt.Printf("Warning: user already exists in MySQL (DN: %s)\n", userDN)
		return
	}

	roles, _ := isql.Role.GetRolesByIds([]uint{2}) // Default user role

	newUser := &model.User{
		Username:      username,
		Nickname:      getAttrOrDefault(entry, "displayName", username),
		GivenName:     getAttrOrDefault(entry, "givenName", username),
		Mail:          getAttrOrDefault(entry, "mail", username+"@"+config.Conf.Ldap.DefaultEmailSuffix),
		JobNumber:     entry.GetAttributeValue("employeeNumber"),
		Mobile:        entry.GetAttributeValue("mobile"),
		PostalAddress: entry.GetAttributeValue("postalAddress"),
		Departments:   getAttrOrDefault(entry, "businessCategory", "Default"),
		Position:      getAttrOrDefault(entry, "departmentNumber", "Member"),
		Introduction:  entry.GetAttributeValue("description"),
		Creator:       "system",
		Source:        "ldap",
		SourceUserId:  username,
		SourceUnionId: username,
		UserDN:        userDN,
		Roles:         roles,
		Status:        1,
		SyncState:     1,
		Password:      tools.NewGenPasswd(config.Conf.Ldap.UserInitPassword),
	}

	if err := isql.User.Add(newUser); err != nil {
		fmt.Printf("Error: failed to add user to MySQL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OK: user synced to MySQL (ID: %d)\n", newUser.ID)
}

// autoSync automatically syncs all mismatched users
func autoSync() {
	fmt.Println("=== Auto sync all users ===\n")

	ldapUsers, err := ildap.User.ListUserDN()
	if err != nil {
		fmt.Printf("Error: failed to get LDAP users: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("LDAP users: %d\n", len(ldapUsers))

	mysqlUsers, err := isql.User.ListAll()
	if err != nil {
		fmt.Printf("Error: failed to get MySQL users: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("MySQL users: %d\n\n", len(mysqlUsers))

	// Build MySQL user DN map
	mysqlDNMap := make(map[string]*model.User)
	for i := range mysqlUsers {
		if mysqlUsers[i].UserDN != "" {
			mysqlDNMap[mysqlUsers[i].UserDN] = mysqlUsers[i]
		}
	}

	// Sync LDAP to MySQL
	fmt.Println("Syncing LDAP -> MySQL...")
	ldapToMySQL := 0
	for _, u := range ldapUsers {
		if u.UserDN == config.Conf.Ldap.AdminDN {
			continue
		}
		if _, exists := mysqlDNMap[u.UserDN]; !exists {
			if username := extractUsername(u.UserDN); username != "" {
				fmt.Printf("  Syncing %s...\n", username)
				syncLDAPToMySQL(username)
				ldapToMySQL++
			}
		}
	}

	// Sync MySQL to LDAP
	fmt.Println("\nSyncing MySQL -> LDAP...")
	ldapDNMap := make(map[string]bool)
	for _, u := range ldapUsers {
		ldapDNMap[u.UserDN] = true
	}

	mysqlToLDAP := 0
	for _, u := range mysqlUsers {
		if u.UserDN == "" || u.UserDN == config.Conf.Ldap.AdminDN {
			continue
		}
		if !ldapDNMap[u.UserDN] {
			fmt.Printf("  Syncing %s (ID: %d)...\n", u.Username, u.ID)
			syncMySQLToLDAP(fmt.Sprintf("%d", u.ID))
			mysqlToLDAP++
		}
	}

	fmt.Printf("\nDone: LDAP->MySQL: %d, MySQL->LDAP: %d\n", ldapToMySQL, mysqlToLDAP)
}

func extractUsername(dn string) string {
	parts := strings.Split(dn, ",")
	if len(parts) > 0 && strings.HasPrefix(parts[0], "uid=") {
		return strings.TrimPrefix(parts[0], "uid=")
	}
	return ""
}

func getAttrOrDefault(entry *ldap.Entry, attr, defaultVal string) string {
	if v := entry.GetAttributeValue(attr); v != "" {
		return v
	}
	return defaultVal
}

