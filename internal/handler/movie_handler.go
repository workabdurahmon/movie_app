package handler

import (
	"movie_app/internal/domain"
	"movie_app/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	service service.MovieService
}

func NewMovieHandler(svc service.MovieService) *MovieHandler {
	return &MovieHandler{service: svc}
}

// @Summary Create a new movie
// @Description Create a new movie in the system
// @Tags movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param movie body domain.CreateMovieRequest true "Movie object"
// @Success 201 {object} domain.Movie
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies [post]
func (h *MovieHandler) CreateMovie(ctx *gin.Context) {
	var req domain.CreateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	movie := &domain.Movie{
		Title:       req.Title,
		Director:    req.Director,
		Year:        req.Year,
		Plot:        req.Plot,
		Genre:       req.Genre,
		Rating:      req.Rating,
		Duration:    req.Duration,
	}

	result, err := h.service.Create(ctx.Request.Context(), movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// @Summary Get a movie by ID
// @Description Get a movie's details by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Movie ID"
// @Success 200 {object} domain.Movie
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [get]
func (h *MovieHandler) GetMovie(ctx *gin.Context) {
	movieID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})

		return
	}

	movie, err := h.service.GetByID(ctx.Request.Context(), uint(movieID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})

		return
	}

	ctx.JSON(http.StatusOK, movie)
}

// @Summary Get all movies
// @Description Get a list of all movies
// @Tags movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} domain.Movie
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies [get]
func (h *MovieHandler) GetAllMovies(ctx *gin.Context) {
	movies, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// @Summary Update a movie
// @Description Update an existing movie's details
// @Tags movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Movie ID"
// @Param movie body domain.UpdateMovieRequest true "Movie object"
// @Success 200 {object} domain.Movie
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{id} [put]
func (h *MovieHandler) UpdateMovie(ctx *gin.Context) {
	movieID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})

		return
	}

	var req domain.UpdateMovieRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})

		return
	}

	// Get existing movie
	movie, err := h.service.GetByID(ctx.Request.Context(), uint(movieID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})

		return
	}

	if req.Title != nil {
		movie.Title = *req.Title
	}

	if req.Director != nil {
		movie.Director = *req.Director
	}

	if req.Year != nil {
		movie.Year = *req.Year
	}

	if req.Plot != nil {
		movie.Plot = *req.Plot
	}

	if req.Genre != nil {
		movie.Genre = *req.Genre
	}

	if req.Rating != nil {
		movie.Rating = *req.Rating
	}

	if req.Duration != nil {
		movie.Duration = *req.Duration
	}

	result, err := h.service.Update(ctx.Request.Context(), movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary Delete a movie
// @Description Delete a movie by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Movie ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(ctx *gin.Context) {
	movieID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})

		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(movieID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
