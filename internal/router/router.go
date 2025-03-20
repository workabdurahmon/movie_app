package router

import (
	"movie_app/internal/handler"
	"movie_app/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	uberfx "go.uber.org/fx"
)

type RouterParams struct {
	uberfx.In

	MovieHandler   *handler.MovieHandler
	UserHandler    *handler.UserHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func NewRouter(p RouterParams) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/api/v1/auth/register", p.UserHandler.Register)
	router.POST("/api/v1/auth/login", p.UserHandler.Login)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(p.AuthMiddleware.Authenticate())

	// User routes
	protected.GET("/users/me", p.UserHandler.GetUser)

	// Movie routes
	protected.POST("/movies", p.MovieHandler.CreateMovie)
	protected.GET("/movies/:id", p.MovieHandler.GetMovie)
	protected.GET("/movies", p.MovieHandler.GetAllMovies)
	protected.PUT("/movies/:id", p.MovieHandler.UpdateMovie)
	protected.DELETE("/movies/:id", p.MovieHandler.DeleteMovie)

	return router
}
