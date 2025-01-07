package controllers

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	middlewares "github.com/mackenzii/freemusic/internal/middleware"
	"github.com/mackenzii/freemusic/internal/models"
	"github.com/mackenzii/freemusic/internal/services"
	"github.com/oklog/ulid/v2"
	"golang.org/x/oauth2"
)

type AuthController struct {
	AuthService  *services.AuthService
	ImageService *services.ImageService
}

func NewAuthController(authService *services.AuthService, imageService *services.ImageService) *AuthController {
	return &AuthController{
		AuthService:  authService,
		ImageService: imageService,
	}
}

func (ctrl *AuthController) RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Role == "" {
		req.Role = "Fan"
	}

	userInfo := models.Users{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Location:     req.Location,
		ConfirmationToken: ulid.MustNew(ulid.Timestamp(time.Now()),
			ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)).String(),
	}

	user, err := ctrl.AuthService.RegisterUser(userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = ctrl.AuthService.EmailService.SendConfirmationEmail(user.Email, user.ConfirmationToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send confirmation email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Registration successful, please check your email"})
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Location string `json:"location"`
	Role     string `json:"role"`
}

type UserResponse struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	Location          string `json:"location"`
	Role              string `json:"role"`
	ConfirmationToken string `json:"confirmationToken"`
}

func (ctrl *AuthController) UserHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := ctrl.AuthService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (ctrl *AuthController) UserUpdate(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Location string `json:"location"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := ctrl.AuthService.UpdateUser(userIDStr, models.Users{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Location:     req.Location,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (ctrl *AuthController) LoginHandler(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := ctrl.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (ctrl *AuthController) RefreshHandler(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newAccessToken, newRefreshToken, err := ctrl.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}

// GoogleLogin redirige l'utilisateur vers la page de connexion Google
// @Summary Rediriger l'utilisateur vers la page de connexion Google
// @Description Rediriger l'utilisateur vers la page de connexion Google
// @Tags Auth
// @Produce json
// @Success 302 {string} string
// @Router /api/auth/google [get]
// GoogleLogin redirige l'utilisateur vers la page de connexion Google
// GoogleLogin redirige l'utilisateur vers la page de connexion Google
func (ctrl *AuthController) GoogleLogin(c *fiber.Ctx) error {
	url := ctrl.AuthService.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

// GoogleCallback gère le callback de Google après l'authentification
func (ctrl *AuthController) GoogleCallback(c *fiber.Ctx) error {
	var req struct {
		IDToken string `json:"idToken"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.IDToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing token"})
	}

	token, err := ctrl.AuthService.GoogleOauthConfig.TokenSource(context.Background(), &oauth2.Token{
		AccessToken: req.IDToken,
	}).Token()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	client := ctrl.AuthService.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode user info"})
	}

	// Rechercher l'utilisateur dans la base de données
	user, err := ctrl.AuthService.GetUserByEmail(userInfo.Email)
	if err != nil {
		// Si l'utilisateur n'existe pas, créez un nouvel utilisateur
		user, err = ctrl.AuthService.RegisterUser(models.Users{
			Email: userInfo.Email,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
		}
	}

	// Générer un token JWT pour l'utilisateur
	userID, err := ulid.Parse(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user ID"})
	}
	accessToken, err := middlewares.GenerateToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken})
}

func (ctrl *AuthController) ConfirmEmailHandler(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token manquant"})
	}

	var user models.Users
	if err := ctrl.AuthService.DB.Where("confirmation_token = ?", token).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Utilisateur non trouvé"})
	}

	user.IsConfirmed = true
	user.ConfirmationToken = ""

	if err := ctrl.AuthService.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur lors de la confirmation de l'email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email confirmé avec succès"})
}

func (ctrl *AuthController) GetUsersHandler(c *fiber.Ctx) error {
	users, err := ctrl.AuthService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (ctrl *AuthController) GetPublicUserInfoHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	publicInfo, err := ctrl.AuthService.GetPublicUserInfoByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(publicInfo)
}

// DeleteUserHandler gère la demande de suppression d'un utilisateur
func (ctrl *AuthController) DeleteUserHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err := ctrl.AuthService.DeleteUserAndRelatedData(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Utilisateur supprimé avec succès"})
}
