package middleware

import (
	"context"
	"net/http"
	"strings"

	"bookstore/internal/repository"
	"bookstore/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				jsonError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				jsonError(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			claims, err := utils.ValidateToken(parts[1], jwtSecret)
			if err != nil {
				jsonError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(userRepo repository.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r)
			if userID == "" {
				jsonError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			user, err := userRepo.GetByID(userID)
			if err != nil {
				jsonError(w, http.StatusUnauthorized, "User not found")
				return
			}

			if !user.IsAdmin {
				jsonError(w, http.StatusForbidden, "Admin access required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(UserIDKey).(string)
	return userID
}
