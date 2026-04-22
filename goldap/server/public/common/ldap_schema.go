package common

import (
	"goldap-server/config"

	ldap "github.com/go-ldap/ldap/v3"
)

// EnsureSudoSchema checks if sudo schema is loaded by verifying sudoers OU exists
// The actual schema loading is handled by ldap-init container in docker-compose
func EnsureSudoSchema() error {
	conn, err := GetLDAPConn()
	if err != nil {
		return err
	}
	defer PutLADPConn(conn)

	// Check if ou=sudoers exists (indicates sudo schema is loaded and structure is initialized)
	sudoersOU := "ou=sudoers," + config.Conf.Ldap.BaseDN
	searchRequest := ldap.NewSearchRequest(
		sudoersOU,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=organizationalUnit)",
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil // ou=sudoers doesn't exist yet, ldap-init will create it
	}

	if len(result.Entries) == 0 {
		return nil // Not found, but not an error
	}

	return nil // Sudo structure verified
}

// CheckSudoSchemaExists checks if sudo schema exists by searching for sudoRole entries
func CheckSudoSchemaExists() (bool, error) {
	conn, err := GetLDAPConn()
	if err != nil {
		return false, err
	}
	defer PutLADPConn(conn)

	// Try to search for sudoRole objects - if it works, schema is loaded
	sudoersOU := "ou=sudoers," + config.Conf.Ldap.BaseDN
	searchRequest := ldap.NewSearchRequest(
		sudoersOU,
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases,
		1, 0, false,
		"(objectClass=sudoRole)",
		[]string{"dn"},
		nil,
	)

	_, err = conn.Search(searchRequest)
	if err != nil {
		return false, nil // Schema not loaded or ou=sudoers doesn't exist
	}

	return true, nil
}
