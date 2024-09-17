package handlers

import (
	"dreampic/templates/pages/settings"
	"net/http"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(r, w, settings.AccountPage(user))
}
