package auth

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"

	"devin/auth/keys"
)

// GetJWTSignKey return private key for RSA sign.
func GetJWTSignKey() (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(keys.JWT_RSA_PRIVATE))
}

// GetJWTVerifyKey return public key for RSA verification.
func GetJWTVerifyKey() (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM([]byte(keys.JWT_RSA_PUBLIC))
}
