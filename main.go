package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/handler"
	"github.com/mbaitar/levenue-assignment/pkg/sb"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", public())
	router.Get("/status", handler.Make(handler.HandleStatus))
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.Make(handler.HandleLoginWithGoogle))
	router.Get("/signup", handler.Make(handler.HandleSignupIndex))
	router.Post("/signup", handler.Make(handler.HandleSignupCreate))
	router.Post("/logout", handler.Make(handler.HandleLogoutCreate))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	router.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	router.Get("/account/stripe/callback", handler.Make(handler.HandleStripeAuthCallback))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth)
		auth.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))
		auth.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))
		auth.Get("/account/setup/type", handler.Make(handler.HandleAccountSetupTypeIndex))
		auth.Post("/account/setup/type", handler.Make(handler.HandleAccountSetupTypeCreate))
		auth.Post("/account/stripe/onboarding", handler.Make(handler.HandleStripeAuth))
		auth.Get("/account/stripe/completed", handler.Make(handler.HandleStripeConnectCompleted))
	})

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth, handler.WithAccountSetup)
		auth.Get("/dashboard", handler.Make(handler.HandleDashboardIndex))
		auth.Post("/rerun", handler.Make(handler.HandleRunMetricCalculation))
		auth.Post("/trade", handler.Make(handler.HandleTradeCreate))
		auth.Get("/history", handler.Make(handler.HandleHistoryIndex))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
