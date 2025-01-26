package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/youngprinnce/go-ecom/config"
	"github.com/youngprinnce/go-ecom/controller/auth"
	"github.com/youngprinnce/go-ecom/types"
	"github.com/youngprinnce/go-ecom/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/users")
	userRouter.POST("/login", h.handleLogin)
	userRouter.POST("/register", h.handleRegister)
}

func (h *Handler) handleLogin(c *gin.Context) {
	var user types.LoginUserPayload
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err,
			"email": user.Email,
		}).Error("Failed to parse login request")
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWT_SECRET)
	token, err := auth.CreateJWT(secret, u.ID, u.Role)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"email": user.Email,
	}).Info("User logged in successfully")
	utils.WriteJSON(c.Writer, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(c *gin.Context) {
	var user types.RegisterUserPayload
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to parse registration request")
		utils.WriteError(c.Writer, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// Validate role (default to "user" if not provided)
	if user.Role == "" {
		user.Role = "user"
	}

	// Check if user exists
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Create user
	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
		Role:      user.Role,
	})
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	// Log successful registration
	utils.Log.WithFields(logrus.Fields{
		"email": user.Email,
	}).Info("New user registered")
	utils.WriteJSON(c.Writer, http.StatusCreated, "user created")
}

func (h *Handler) handleGetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(c.Writer, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c.Writer, http.StatusOK, user)
}
