package handlers

import (
	"dreampic/types"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/gorilla/sessions"
)

var SessionEnv = os.Getenv("SESSION_SECRET")

func checkHTMX(w http.ResponseWriter, r *http.Request) {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if !htmx.IsHTMX(r) {
		// If not, return HTTP 400 error.
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("request API", "method", r.Method, "status", http.StatusBadRequest, "path", r.URL.Path)
		return
	}

}

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)

	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func Make(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal serve error", "err", err, "path", r.URL.Path)
		}
	}
}

func cookieSession(r *http.Request, envKey, cookieName string) (*sessions.Session, error) {
	store := sessions.NewCookieStore([]byte(envKey))

	return store.Get(r, cookieName)
}
