package jwt_validator

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

/*
LoadPublicKey загружает RSA public key из PEM-файла.

Поддерживается формат:
-----BEGIN PUBLIC KEY-----
...
-----END PUBLIC KEY-----

Ошибки:
- invalid PEM block
- not RSA public key
- parse error
*/
func LoadPublicKey(path string) (*rsa.PublicKey, error) {

	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return pubKey, nil
}
