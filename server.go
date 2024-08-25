package main

import (
	"dreampic/handlers"
	"dreampic/pkg/sb"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	gowebly "github.com/gowebly/helpers"
)

//go:embed all:static
var static embed.FS

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return sb.Init()
}

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {

	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	// Validate environment variables.
	port := os.Getenv("BACKEND_PORT")

	// Create a new chi router.
	router := chi.NewRouter()

	// Use chi middlewares.
	router.Use(middleware.Logger)
	router.Use(handlers.WithUser)

	// Handle static files from the embed FS (with a custom handler).
	router.Handle("/static/*", gowebly.StaticFileServerHandler(http.FS(static)))

	// Handle API endpoints.
	router.Get("/api/hello-world", handlers.Make(handlers.ShowContentAPIHandler))
	router.Get("/api/show-ronin", handlers.ShowRonindevMessage)

	// Handle index page view.
	router.Get("/", handlers.Make(handlers.HandlerHomeIndex))

	// Handle signin page view.
	router.Get("/signin", handlers.Make(handlers.HandleSignIn))
	router.Post("/signin", handlers.Make(handlers.HandleSignInCreate))

	router.Get("/signup", handlers.Make(handlers.HandleSignUp))
	router.Post("/signup", handlers.Make(handlers.HandleSignUpCreate))

	router.Get("/auth/callback", handlers.Make(handlers.HandleAuthCallback))

	router.Post("/logout", handlers.Make(handlers.HandleLogoutCreate))

	router.Get("/settings", handlers.Make(handlers.HandleSettings))

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	server := &http.Server{
		Addr:         port,
		Handler:      router, // handle all chi routes
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Send log message.
	slog.Info("Starting server...", "port", port)

	return server.ListenAndServe()
}
