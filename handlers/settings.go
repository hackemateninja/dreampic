package handlers

import (
	"dreampic/templates/pages"
	"net/http"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(r, w, pages.SettingsPage(user))
}
