package openldap

import (
	"fmt"
	"strings"

	"goldap-server/config"
	"goldap-server/public/common"

	ldap "github.com/go-ldap/ldap/v3"
)

type Dept struct {
	DN       string `json:"dn"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
	ParentId string `json:"parentid"`
}

type User struct {
	Name             string   `json:"name"`
	DN               string   `json:"dn"`
	CN               string   `json:"cn"`
	SN               string   `json:"sn"`
	Mobile           string   `json:"mobile"`
	BusinessCategory string   `json:"businessCategory"`
	DepartmentNumber string   `json:"departmentNumber"`
	Description      string   `json:"description"`
	DisplayName      string   `json:"displayName"`
	Mail             string   `json:"mail"`
	EmployeeNumber   string   `json:"employeeNumber"`
	GivenName        string   `json:"givenName"`
	PostalAddress    string   `json:"postalAddress"`
	DepartmentIds    []string `json:"department_ids"`
}

// GetAllDepts returns all groups (OU, CN, posixGroup, sudoRole)
func GetAllDepts() (ret []*Dept, err error) {
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(|(objectClass=organizationalUnit)(objectClass=groupOfUniqueNames)(objectClass=posixGroup)(objectClass=sudoRole))",
		[]string{"DN", "cn", "ou", "description"},
		nil,
	)

	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return nil, err
	}

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return ret, err
	}
	if len(sr.Entries) > 0 {
		for _, v := range sr.Entries {
			if v.DN == config.Conf.Ldap.BaseDN || v.DN == config.Conf.Ldap.AdminDN {
				continue
			}
			
			// Skip user entries
			if strings.HasPrefix(strings.ToLower(v.DN), "uid=") {
				continue
			}
			
			var ele Dept
			ele.DN = v.DN
			
			// Extract group name
			ele.Name = v.GetAttributeValue("cn")
			if ele.Name == "" {
				ele.Name = v.GetAttributeValue("ou")
			}
			if ele.Name == "" {
				dnParts := strings.Split(v.DN, ",")
				if len(dnParts) > 0 {
					firstPart := strings.Split(dnParts[0], "=")
					if len(firstPart) == 2 {
						ele.Name = firstPart[1]
					}
				}
			}
			
			ele.Id = ele.Name
			ele.Remark = v.GetAttributeValue("description")
			
			// Calculate parent ID
			if len(strings.Split(v.DN, ","))-len(strings.Split(config.Conf.Ldap.BaseDN, ",")) == 1 {
				ele.ParentId = "0"
			} else {
				dnParts := strings.Split(v.DN, ",")
				if len(dnParts) > 1 {
					parentPart := strings.Split(dnParts[1], "=")
					if len(parentPart) == 2 {
						ele.ParentId = parentPart[1]
					} else {
						ele.ParentId = "0"
					}
				} else {
					ele.ParentId = "0"
				}
			}
			ret = append(ret, &ele)
		}
	}
	return
}

// GetAllUsers returns all users
func GetAllUsers() (ret []*User, err error) {
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))",
		[]string{},
		nil,
	)

	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return nil, err
	}

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return ret, err
	}
	if len(sr.Entries) > 0 {
		for _, v := range sr.Entries {
			if v.DN == config.Conf.Ldap.UserDN || !strings.Contains(v.DN, config.Conf.Ldap.UserDN) {
				continue
			}
			name := strings.Split(strings.Split(v.DN, ",")[0], "=")[1]
			deptIds, err := GetUserDeptIds(v.DN)
			if err != nil {
				return ret, err
			}
			ret = append(ret, &User{
				Name:             name,
				DN:               v.DN,
				CN:               v.GetAttributeValue("cn"),
				SN:               v.GetAttributeValue("sn"),
				Mobile:           v.GetAttributeValue("mobile"),
				BusinessCategory: v.GetAttributeValue("businessCategory"),
				DepartmentNumber: v.GetAttributeValue("departmentNumber"),
				Description:      v.GetAttributeValue("description"),
				DisplayName:      v.GetAttributeValue("displayName"),
				Mail:             v.GetAttributeValue("mail"),
				EmployeeNumber:   v.GetAttributeValue("employeeNumber"),
				GivenName:        v.GetAttributeValue("givenName"),
				PostalAddress:    v.GetAttributeValue("postalAddress"),
				DepartmentIds:    deptIds,
			})
		}
	}
	return
}

// GetUserDeptIds gets user's department IDs
func GetUserDeptIds(udn string) (ret []string, err error) {
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(|(Member=%s)(uniqueMember=%s))", udn, udn),
		[]string{},
		nil,
	)

	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return nil, err
	}

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return ret, err
	}
	if len(sr.Entries) > 0 {
		for _, v := range sr.Entries {
			ret = append(ret, strings.Split(strings.Split(v.DN, ",")[0], "=")[1])
		}
	}
	return ret, nil
}
