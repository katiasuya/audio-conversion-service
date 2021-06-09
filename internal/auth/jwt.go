// Package auth provides methods to work with jwt tokens and user authorization.
package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenManager has methods to use jwt and contains private and public keys.
type TokenManager struct {
	privateKey string
	publicKey  string
}

// New returns new token manager with the given keys.
func New(publicKey, privateKey string) *TokenManager {
	return &TokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// ParseJWT validates and parses the given jwt access token.
func (tm *TokenManager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(tm.publicKey))
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("could not get user claims from token")
	}

	return claims["sub"].(string), nil
}

// NewJWT creates new JWT token based on user id and private key.
func (tm *TokenManager) NewJWT(userID string) (string, error) {
	const expTimeHrs = 24
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Hour * expTimeHrs).Unix(),
	})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(tm.privateKey))
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
