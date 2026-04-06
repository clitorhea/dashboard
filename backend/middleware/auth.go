package middleware

import (
	"context"
	"net/http"

	"github.com/rhea/nas-dashboard/db"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	db *db.DB
}

func NewAuthMiddleware(database *db.DB) *AuthMiddleware {
	return &AuthMiddleware{db: database}
}

func (m *AuthMiddleware) Require(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		session, err := m.db.GetSession(cookie.Value)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		next(w, r.WithContext(ctx))
	}
}
