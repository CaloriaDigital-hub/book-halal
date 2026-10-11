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
	// Auth (public)
	mux.HandleFunc("POST /api/auth/register/initiate", a.container.AuthController.InitiateRegistration)
	mux.HandleFunc("POST /api/auth/register/confirm", a.container.AuthController.ConfirmRegistration)
	mux.HandleFunc("POST /api/auth/sign-in", a.container.AuthController.SignIn)

	// Books — каталог и карточка книги открыты всем (витрина); has_access считается, если передан токен
	mux.Handle("GET /api/books/catalog", a.container.Public(a.container.BooksController.GetCatalog))
	mux.Handle("GET /api/books/{id}", a.container.Public(a.container.BooksController.GetBookByID))
	mux.Handle("GET /api/books/{id}/pages", a.container.Authenticated(a.container.BooksController.GetPages))
	mux.Handle("GET /api/books/{id}/pages/{page}", a.container.Authenticated(a.container.BooksController.GetPage))
	mux.Handle("POST /api/books/{id}/ask", a.container.Authenticated(a.container.BooksController.AskBook))

	// Reader — progress
	mux.Handle("PUT /api/books/{id}/progress", a.container.Authenticated(a.container.ReaderController.SaveProgress))
	mux.Handle("GET /api/books/{id}/progress", a.container.Authenticated(a.container.ReaderController.GetProgress))

	// Reader — bookmarks
	mux.Handle("POST /api/books/{id}/bookmarks", a.container.Authenticated(a.container.ReaderController.AddBookmark))
	mux.Handle("GET /api/books/{id}/bookmarks", a.container.Authenticated(a.container.ReaderController.GetBookmarks))
	mux.Handle("DELETE /api/books/{id}/bookmarks/{page}", a.container.Authenticated(a.container.ReaderController.RemoveBookmark))

	// Purchases — покупатель
	mux.Handle("GET /api/payments/requisites", a.container.Authenticated(a.container.PurchaseController.GetRequisites))
	mux.Handle("POST /api/books/{id}/purchase", a.container.Authenticated(a.container.PurchaseController.CreatePurchase))
	mux.Handle("GET /api/books/{id}/purchase", a.container.Authenticated(a.container.PurchaseController.GetMyPurchase))

	// Admin
	mux.Handle("POST /api/admin/books/upload", a.container.AdminOnly(a.container.BooksController.UploadBook))
	mux.Handle("PUT /api/admin/payments/requisites", a.container.AdminOnly(a.container.PurchaseController.UpdateRequisites))
	mux.Handle("GET /api/admin/purchases", a.container.AdminOnly(a.container.PurchaseController.ListPurchases))
	mux.Handle("POST /api/admin/purchases/{id}/approve", a.container.AdminOnly(a.container.PurchaseController.ApprovePurchase))
	mux.Handle("POST /api/admin/purchases/{id}/reject", a.container.AdminOnly(a.container.PurchaseController.RejectPurchase))
	mux.Handle("GET /api/admin/purchases/{id}/receipt", a.container.AdminOnly(a.container.PurchaseController.GetReceipt))

	// Static files (covers book images and covers)
	// PDF книг не раздаём: иначе платную книгу можно скачать в обход оплаты
	fs := http.FileServer(middleware.NoListFileSystem(middleware.DenyExtensions(http.Dir(a.staticRoot), ".pdf")))
	mux.Handle("GET /static/", middleware.CacheStatic(http.StripPrefix("/static/", fs)))

	// API Documentation (Scalar UI)
	docsFs := http.FileServer(http.Dir("docs"))
	mux.Handle("GET /docs/", http.StripPrefix("/docs/", docsFs))
}

// Container — доступ к фоновым компонентам (Telegram-бот) из main.
func (a *App) Container() *Container {
	return a.container
}

func (a *App) Run() error {
	a.logger.Info("starting server", "addr", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("shutting down server...")
	return a.server.Shutdown(ctx)
}