package handlers

import (
	"context"
	"dreampic/pkg/sb"
	"dreampic/types"
	"net/http"
	"strings"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}

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
			Email:    resp.Email,
			Avatar:   "https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp",
			LoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
