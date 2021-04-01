// Package auth provides methods to work with jwt tokens and user authorization.
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/katiasuya/audio-conversion-service/internal/server/context"
	"github.com/katiasuya/audio-conversion-service/internal/server/response"
)

// TokenManager has methods to use jwt and contains private and public keys.
type TokenManager struct {
	privateKey []byte
	publicKey  []byte
}

// New returns new token manager with the given keys.
func New(publicKey, privateKey []byte) *TokenManager {
	return &TokenManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// IsAuthorized is a middleware that checks user authorization.
func (tm *TokenManager) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			response.RespondErr(w, http.StatusUnauthorized, errors.New("malformed token"))
			return
		}

		jwtToken := authHeader[1]
		claimUserID, err := tm.ParseJWT(jwtToken)
		if err != nil {
			response.RespondErr(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.ContextWithUserID(r.Context(), claimUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ParseJWT validates and parses the given jwt access token.
func (tm *TokenManager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM(tm.publicKey)
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

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(tm.privateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
