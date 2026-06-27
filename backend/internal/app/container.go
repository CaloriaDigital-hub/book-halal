package app

import (
	bookCommands "book_halal/internal/application/books/commands"
	bookQueries "book_halal/internal/application/books/queries"
	readerCommands "book_halal/internal/application/reader/commands"
	readerQueries "book_halal/internal/application/reader/queries"
	userCommands "book_halal/internal/application/users/commands"
	"book_halal/internal/config"
	"book_halal/internal/infrastructure/email"
	"book_halal/internal/infrastructure/pdf"
	postgresBooks "book_halal/internal/infrastructure/postgres/books"
	postgresReader "book_halal/internal/infrastructure/postgres/reader"
	postgresSessions "book_halal/internal/infrastructure/postgres/sessions"
	postgresUsers "book_halal/internal/infrastructure/postgres/users"
	ragclient "book_halal/internal/infrastructure/ragclient"
	redisinfra "book_halal/internal/infrastructure/redis"
	"book_halal/internal/middleware"
	booksController "book_halal/internal/presentation/http/controllers/books"
	readerController "book_halal/internal/presentation/http/controllers/reader"
	usersController "book_halal/internal/presentation/http/controllers/users"

	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type Container struct {
	AuthController   *usersController.AuthController
	BooksController  *booksController.BookController
	ReaderController *readerController.ReaderController

	authMiddleware func(http.Handler) http.Handler
}

func (c *Container) AdminOnly(handler http.HandlerFunc) http.Handler {
	return c.authMiddleware(middleware.RequireAdmin(handler))
}

func (c *Container) Authenticated(handler http.HandlerFunc) http.Handler {
	return c.authMiddleware(handler)
}

func NewContainer(cfg *config.Config, pgPool *pgxpool.Pool, redisClient *goredis.Client, logger *slog.Logger) *Container {
	// Infrastructure
	userRepo := postgresUsers.NewUserRepository(pgPool)
	bookRepo := postgresBooks.NewBookRepository(pgPool)
	sessionRepo := postgresSessions.NewSessionRepository(pgPool)
	readerRepo := postgresReader.NewReaderRepository(pgPool)
	verificationRepo := redisinfra.NewVerificationRepository(redisClient)
	emailSender := email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	bookProcessor := pdf.NewPDFToPPMProcessor(bookRepo, cfg.StaticRoot, logger)
	ragClient := ragclient.New(cfg.RAGServiceURL)

	// Application — Users
	initiateHandler := userCommands.NewInitiateRegistrationHandler(userRepo, verificationRepo, emailSender)
	confirmHandler := userCommands.NewConfirmRegistrationHandler(userRepo, verificationRepo)
	signInHandler := userCommands.NewSignInHandler(userRepo, sessionRepo)

	// Application — Books (Commands & Queries)
	uploadBookHandler := bookCommands.NewUploadBookHandler(bookRepo, bookProcessor, ragClient, logger)
	getBooksQuery := bookQueries.NewGetBooksHandler(bookRepo)
	getBookByIDQuery := bookQueries.NewGetBookByIDHandler(bookRepo)
	getBookPagesQuery := bookQueries.NewGetBookPagesHandler(bookRepo)
	getPageQuery := bookQueries.NewGetPageHandler(bookRepo)

	// Application — Reader (Commands & Queries)
	saveProgressHandler := readerCommands.NewSaveProgressHandler(readerRepo)
	getProgressHandler := readerQueries.NewGetProgressHandler(readerRepo)
	addBookmarkHandler := readerCommands.NewAddBookmarkHandler(readerRepo)
	getBookmarksHandler := readerQueries.NewGetBookmarksHandler(readerRepo)
	removeBookmarkHandler := readerCommands.NewRemoveBookmarkHandler(readerRepo)

	// Presentation
	return &Container{
		AuthController: usersController.NewAuthController(initiateHandler, confirmHandler, signInHandler),
		BooksController: booksController.NewBookController(
			uploadBookHandler,
			getBooksQuery,
			getBookByIDQuery,
			getBookPagesQuery,
			getPageQuery,
			ragClient,
			cfg.BaseURL,
		),
		ReaderController: readerController.NewReaderController(
			saveProgressHandler,
			getProgressHandler,
			addBookmarkHandler,
			getBookmarksHandler,
			removeBookmarkHandler,
		),
		authMiddleware: middleware.Authenticate(sessionRepo, userRepo),
	}
}