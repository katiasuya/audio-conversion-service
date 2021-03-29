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

// TokenManager has mathods to use jwt and contains secret key.
type TokenManager struct {
	secretKey string
}

// New returns new token manager with the given secret key.
func New(secretKey string) *TokenManager {
	return &TokenManager{secretKey: secretKey}
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

		ctx := context.SetWithUserID(r.Context(), claimUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ParseJWT validates and parses the given jwt access token.
func (tm *TokenManager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
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

// NewJWT creates new JWT token based on user id and secret key.
func (tm *TokenManager) NewJWT(userID string) (string, error) {
	const expTimeHrs = 24
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Hour * expTimeHrs).Unix(),
	})

	return token.SignedString([]byte(tm.secretKey))
}
