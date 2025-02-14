package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	ContextKeyUsername = contextKey("username")
)

// AuthMiddleware проверяет заголовок Authorization, извлекает
// и валидирует JWT-токен и добавляет username в контекст запроса
func AuthMiddleware(jwtSecret []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("неожиданный метод подписи")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok || username == "" {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		// Добавляем username в контекст запроса
		ctx := context.WithValue(r.Context(), ContextKeyUsername, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
