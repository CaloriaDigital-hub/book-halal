package app

import (
	"book_halal/internal/application/books/commands/handlers"
	queryHandlersBooks "book_halal/internal/application/books/queries/handlers"
	userHandlers "book_halal/internal/application/users/commands/handlers"
	"book_halal/internal/config"
	"book_halal/internal/infrastructure/email"
	"book_halal/internal/infrastructure/pdf"
	postgresBooks "book_halal/internal/infrastructure/postgres/books"
	postgresSessions "book_halal/internal/infrastructure/postgres/sessions"
	postgresUsers "book_halal/internal/infrastructure/postgres/users"
	redisinfra "book_halal/internal/infrastructure/redis"
	"book_halal/internal/middleware"
	booksController "book_halal/internal/presentation/http/controllers/books"
	usersController "book_halal/internal/presentation/http/controllers/users"

	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type Container struct {
	AuthController  *usersController.AuthController
	BooksController *booksController.BookController

	authMiddleware func(http.Handler) http.Handler
}

// AdminOnly wraps a handler so it requires authentication and admin role.
func (c *Container) AdminOnly(handler http.HandlerFunc) http.Handler {
	return c.authMiddleware(middleware.RequireAdmin(handler))
}

// Authenticated wraps a handler so it requires authentication only.
func (c *Container) Authenticated(handler http.HandlerFunc) http.Handler {
	return c.authMiddleware(handler)
}

func NewContainer(cfg *config.Config, pgPool *pgxpool.Pool, redisClient *goredis.Client, logger *slog.Logger) *Container {
	// Infrastructure
	userRepo := postgresUsers.NewUserRepository(pgPool)
	bookRepo := postgresBooks.NewBookRepository(pgPool)
	sessionRepo := postgresSessions.NewSessionRepository(pgPool)
	verificationRepo := redisinfra.NewVerificationRepository(redisClient)
	emailSender := email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	bookProcessor := pdf.NewPDFToPPMProcessor(bookRepo, cfg.StaticRoot, logger)

	// Application — Users
	initiateHandler := userHandlers.NewInitiateRegistrationHandler(userRepo, verificationRepo, emailSender)
	confirmHandler := userHandlers.NewConfirmRegistrationHandler(userRepo, verificationRepo)
	signInHandler := userHandlers.NewSignInCommandHandler(userRepo, sessionRepo)

	// Application — Books (Commands & Queries)
	uploadBookHandler := handlers.NewUploadBookCommandHandler(bookRepo, bookProcessor)
	getBooksQuery := queryHandlersBooks.NewGetBooksQueryHandler(bookRepo)
	getBookByIDQuery := queryHandlersBooks.NewGetBookByIDQueryHandler(bookRepo)
	getBookPagesQuery := queryHandlersBooks.NewGetBookPagesQueryHandler(bookRepo)

	// Presentation
	return &Container{
		AuthController:  usersController.NewAuthController(initiateHandler, confirmHandler, signInHandler),
		BooksController: booksController.NewBookController(uploadBookHandler, getBooksQuery, getBookByIDQuery, getBookPagesQuery),

		authMiddleware: middleware.Authenticate(sessionRepo, userRepo),
	}
}