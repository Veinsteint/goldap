package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// RSAEncrypt encrypts data with public key
func RSAEncrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	block, _ := pem.Decode(publicBytes)
	if block == nil {
		return res, fmt.Errorf("invalid public key")
	}

	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("invalid public key: %v", err)
	}
	
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return res, fmt.Errorf("encryption failed: %v", err)
	}
	
	return []byte(EncodeStr2Base64(string(res))), nil
}

// RSADecrypt decrypts data with private key (supports PKCS#1 and PKCS#8)
func RSADecrypt(base64Data, privateBytes []byte) ([]byte, error) {
	var res []byte
	data := []byte(DecodeStrFromBase64(string(base64Data)))
	block, _ := pem.Decode(privateBytes)
	if block == nil {
		return res, fmt.Errorf("invalid private key")
	}
	
	var privateKey *rsa.PrivateKey
	var err error
	
	// Try PKCS#1 format first
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS#8 format
		keyInterface, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return res, fmt.Errorf("invalid private key (PKCS#1: %v, PKCS#8: %v)", err, err2)
		}
		var ok bool
		privateKey, ok = keyInterface.(*rsa.PrivateKey)
		if !ok {
			return res, fmt.Errorf("invalid private key type")
		}
	}
	
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return res, fmt.Errorf("decryption failed: %v", err)
	}
	return res, nil
}

// EncodeStr2Base64 encodes string to base64
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// DecodeStrFromBase64 decodes base64 string
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}
