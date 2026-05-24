package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var errInvalidPassword = errors.New("Неверный пароль")

func passwordHash(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func signingKey(pass string) []byte {
	sum := sha256.Sum256([]byte(pass))
	return sum[:]
}

// CreateToken формирует JWT-токен для текущего пароля.
func CreateToken() (string, error) {
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		return "", errInvalidPassword
	}

	claims := jwt.MapClaims{
		"hash": passwordHash(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey(pass))
}

// ValidateToken проверяет JWT-токен и соответствие паролю.
func ValidateToken(tokenString string) bool {
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		return true
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return signingKey(pass), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	hash, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	return hash == passwordHash(pass)
}

// Auth проверяет аутентификацию перед вызовом обработчика API.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		var jwtToken string
		if cookie, err := r.Cookie("token"); err == nil {
			jwtToken = cookie.Value
		}

		if !ValidateToken(jwtToken) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
