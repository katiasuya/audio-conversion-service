// Package auth provides methods to work with jwt tokens and user auth.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
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
func (tm *TokenManager) IsAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			RespondErr(w, http.StatusUnauthorized, errors.New("malformed token"))
			return
		}

		jwtToken := authHeader[1]
		claimUserID, err := tm.ParseToken(jwtToken)
		if err != nil {
			RespondErr(w, http.StatusUnauthorized, err)
			return
		}

		context.Set(r, "userID", claimUserID)
		next(w, r)
	})
}

// ParseToken validates and parses the given jwt access token.
func (tm *TokenManager) ParseToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("could not get user claims from token")
	}

	return claims["sub"].(string), nil
}

// NewJWT creates new JWT token based on user id and secret key.
func (tm *TokenManager) NewJWT(userID string) (string, error) {
	const expTimeMin = 15
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Minute * expTimeMin).Unix(),
	})

	return token.SignedString([]byte(tm.secretKey))
}

// Respond is a function to make http responses.
func Respond(w http.ResponseWriter, code int, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "can't marshal the given payload", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

// RespondErr is a function to make http error responses.
func RespondErr(w http.ResponseWriter, code int, err error) {
	type error struct {
		Code    int
		Message string
	}

	respErr := error{
		Code:    code,
		Message: err.Error(),
	}
	Respond(w, code, respErr)
}
