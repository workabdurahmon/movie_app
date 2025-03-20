package main

import (
	_ "movie_app/docs"
	"movie_app/internal/config"
	"movie_app/internal/handler"
	"movie_app/internal/middleware"
	"movie_app/internal/repository"
	"movie_app/internal/router"
	"movie_app/internal/service"

	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	uberfx "go.uber.org/fx"
)

const (
	readTimeout  = 15 * time.Second
	writeTimeout = 15 * time.Second
	idleTimeout  = 60 * time.Second
	gracePeriod  = 5 * time.Second
)

// @title           Movie API
// @version         1.0
// @description     A RESTful API for managing movies with authentication.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func ProvideRepositories() uberfx.Option {
	return uberfx.Provide(
		uberfx.Annotate(
			repository.NewMovieRepository,
			uberfx.As(new(repository.MovieRepository)),
		),
		uberfx.Annotate(
			repository.NewUserRepository,
			uberfx.As(new(repository.UserRepository)),
		),
	)
}

func ProvideServices() uberfx.Option {
	return uberfx.Provide(
		uberfx.Annotate(
			service.NewMovieService,
			uberfx.As(new(service.MovieService)),
		),
		uberfx.Annotate(
			func(repo repository.UserRepository, cfg *config.Config) service.UserService {
				return service.NewUserService(repo, cfg.JWT.Secret)
			},
			uberfx.As(new(service.UserService)),
		),
	)
}

func ProvideHandlers() uberfx.Option {
	return uberfx.Provide(
		handler.NewMovieHandler,
		handler.NewUserHandler,
	)
}

func NewAuthMiddleware(cfg *config.Config) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(cfg.JWT.Secret)
}

func main() {
	app := uberfx.New(
		uberfx.Provide(config.LoadConfig),

		uberfx.Provide(repository.NewDatabase),

		uberfx.Provide(NewAuthMiddleware),

		// Provide all dependencies
		ProvideRepositories(),
		ProvideServices(),
		ProvideHandlers(),

		// Provide HTTP server
		uberfx.Provide(router.NewRouter),

		// Invoke server start
		uberfx.Invoke(func(router *gin.Engine, cfg *config.Config) {
			srv := &http.Server{
				Addr:         ":" + cfg.Server.Port,
				Handler:      router,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
				IdleTimeout:  idleTimeout,
			}

			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()

			gracefulShutdown(srv, gracePeriod)

		}),
	)

	app.Run()
}

func gracefulShutdown(srv *http.Server, gracePeriod time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Server exited gracefully")
}
