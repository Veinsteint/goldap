package tools

import (
	"crypto/rand"
	"math/big"

	"goldap-server/config"
)

// NewGenPasswd encrypts password with RSA
func NewGenPasswd(passwd string) string {
	pass, _ := RSAEncrypt([]byte(passwd), config.Conf.System.RSAPublicBytes)
	return string(pass)
}

// NewParPasswd decrypts password with RSA
func NewParPasswd(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}

const (
	passwordLength = 8
	letters        = "abcdefghijklmnopqrstu@vwxyzABCDEFGHIJKL#MNOP*QRSTUVWXYZ0123456789"
	lettersLength  = len(letters)
)

// GenerateRandomPassword generates random password
func GenerateRandomPassword() string {
	password := make([]byte, passwordLength)
	for i := range password {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(lettersLength)))
		password[i] = letters[index.Int64()]
	}
	return string(password)
}
