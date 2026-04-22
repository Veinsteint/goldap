package common

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"goldap-server/config"

	ldap "github.com/go-ldap/ldap/v3"
)

var (
	ldapPool    *LdapConnPool
	ldapInit    bool
	ldapInitOne sync.Once
)

func InitLDAP() {
	if ldapInit {
		return
	}

	ldapInitOne.Do(func() {
		ldapInit = true
	})

	conn, err := ldap.DialURL(config.Conf.Ldap.Url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		Log.Panicf("LDAP connection failed: %v", err)
	}

	if config.Conf.Ldap.AllowAnonBinding {
		if err = conn.UnauthenticatedBind(""); err != nil {
			Log.Panicf("Anonymous bind failed: %v", err)
		}
	} else {
		if err = conn.Bind(config.Conf.Ldap.AdminDN, config.Conf.Ldap.AdminPass); err != nil {
			Log.Panicf("Admin bind failed: %v", err)
		}
	}

	ldapPool = &LdapConnPool{
		conns:    make([]*ldap.Conn, 0),
		reqConns: make(map[uint64]chan *ldap.Conn),
		maxOpen:  config.Conf.Ldap.MaxConn,
	}
	PutLADPConn(conn)
}

func GetLDAPConn() (*ldap.Conn, error) {
	return ldapPool.GetConnection()
}

func GetLDAPConnForModify() (*ldap.Conn, error) {
	return initLDAPConnForModify()
}

func PutLADPConn(conn *ldap.Conn) {
	ldapPool.PutConnection(conn)
}

type LdapConnPool struct {
	mu       sync.Mutex
	conns    []*ldap.Conn
	reqConns map[uint64]chan *ldap.Conn
	openConn int
	maxOpen  int
}

func (p *LdapConnPool) GetConnection() (*ldap.Conn, error) {
	p.mu.Lock()

	if len(p.conns) > 0 {
		p.openConn++
		conn := p.conns[0]
		p.conns = p.conns[1:]
		p.mu.Unlock()

		if conn.IsClosing() {
			return initLDAPConn()
		}
		return conn, nil
	}

	if p.maxOpen != 0 && p.openConn > p.maxOpen {
		req := make(chan *ldap.Conn, 1)
		reqKey := p.nextRequestKeyLocked()
		p.reqConns[reqKey] = req
		p.mu.Unlock()
		return <-req, nil
	}

	p.openConn++
	p.mu.Unlock()
	return initLDAPConn()
}

func (p *LdapConnPool) PutConnection(conn *ldap.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.reqConns) > 0 {
		var req chan *ldap.Conn
		var reqKey uint64
		for reqKey, req = range p.reqConns {
			break
		}
		delete(p.reqConns, reqKey)
		req <- conn
		return
	}

	p.openConn--
	if !conn.IsClosing() {
		p.conns = append(p.conns, conn)
	}
}

func (p *LdapConnPool) nextRequestKeyLocked() uint64 {
	for {
		reqKey := rand.Uint64()
		if _, ok := p.reqConns[reqKey]; !ok {
			return reqKey
		}
	}
}

func initLDAPConn() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(config.Conf.Ldap.Url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		return nil, err
	}

	if config.Conf.Ldap.AllowAnonBinding {
		err = conn.UnauthenticatedBind("")
	} else {
		err = conn.Bind(config.Conf.Ldap.AdminDN, config.Conf.Ldap.AdminPass)
	}
	return conn, err
}

func initLDAPConnForModify() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(config.Conf.Ldap.Url, ldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}))
	if err != nil {
		return nil, err
	}

	if err = conn.Bind(config.Conf.Ldap.AdminDN, config.Conf.Ldap.AdminPass); err != nil {
		return nil, fmt.Errorf("admin bind failed: %v", err)
	}
	return conn, nil
}
