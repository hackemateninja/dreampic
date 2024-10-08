package handlers

import (
	"context"
	"dreampic/pkg/sb"
	"dreampic/types"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func verifyStatic(next http.Handler, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/static") {
		next.ServeHTTP(w, r)
		return
	}
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		verifyStatic(next, w, r)
		user := getAuthenticatedUser(r)

		if !user.LoggedIn {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		verifyStatic(next, w, r)

		cookie, err := r.Cookie("at")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := sb.Client.Auth.User(r.Context(), cookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:       uuid.MustParse(resp.ID),
			Email:    resp.Email,
			LoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
