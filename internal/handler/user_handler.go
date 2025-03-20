package handler

import (
	"errors"
	"movie_app/internal/domain"
	"movie_app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.RegisterRequest true "User registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var req domain.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user, err := h.service.Register(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})

			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.LoginRequest true "User credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req domain.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	result, err := h.service.Login(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetUser godoc
// @Summary Get user profile
// @Description Get the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

		return
	}

	user, err := h.service.GetUser(ctx.Request.Context(), userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, user)
}
