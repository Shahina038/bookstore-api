package middleware

import (
	"net/http"

	"bookstore/internal/repository"
)

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
