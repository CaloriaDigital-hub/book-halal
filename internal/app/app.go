package app

import (
	"book_halal/internal/middleware"
	"context"
	"log"
	"log/slog"
	"net/http"
)

type App struct {
	container  *Container
	server     *http.Server
	logger     *slog.Logger
	staticRoot string
}

func NewApp(container *Container, addr string, staticRoot string, logger *slog.Logger) *App {
	a := &App{
		container:  container,
		logger:     logger,
		staticRoot: staticRoot,
	}

	mux := http.NewServeMux()
	a.setupRoutes(mux)

	a.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return a
}

func (a *App) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/register/initiate", a.container.AuthController.InitiateRegistration)
	mux.HandleFunc("POST /api/auth/register/confirm", a.container.AuthController.ConfirmRegistration)
	mux.HandleFunc("POST /api/auth/sign-in", a.container.AuthController.SignIn)

	mux.Handle("GET /api/books/catalog", a.container.Authenticated(a.container.BooksController.GetCatalog))
	mux.Handle("GET /api/books/{id}", a.container.Authenticated(a.container.BooksController.GetBookByID))
	mux.Handle("GET /api/books/{id}/pages", a.container.Authenticated(a.container.BooksController.GetPages))

	mux.Handle("POST /api/admin/books/upload", a.container.AdminOnly(a.container.BooksController.UploadBook))

	fs := http.FileServer(middleware.NoListFileSystem(http.Dir(a.staticRoot)))
mux.Handle("GET /static/", middleware.CacheStatic(http.StripPrefix("/static/", fs)))
}

func (a *App) Run() error {
	a.logger.Info("starting server", "addr", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("shutting down server...")
	return a.server.Shutdown(ctx)
}