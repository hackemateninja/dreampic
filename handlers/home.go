package handlers

import (
	"dreampic/templates/pages/home"
	"dreampic/templates/ui"
	"log/slog"
	"net/http"

	"github.com/angelofallars/htmx-go"
)

func HandlerHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}

// showContentAPIHandler handles an API endpoint to show content.
func ShowContentAPIHandler(w http.ResponseWriter, r *http.Request) error {

	checkHTMX(w, r)

	// // Write HTML content.
	// w.Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))

	// // Send htmx response.
	// htmx.NewResponse().Write(w)

	// Send log message.
	slog.Info("request API", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)

	return ui.Hello("herman").Render(r.Context(), w)
}

func ShowRonindevMessage(w http.ResponseWriter, r *http.Request) {
	checkHTMX(w, r)

	// Write HTML content.
	w.Write([]byte("<h1>this is a hello world message </h1>"))

	// Send htmx response.
	htmx.NewResponse().Write(w)

	// Send log message.
	slog.Info("request API", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}
