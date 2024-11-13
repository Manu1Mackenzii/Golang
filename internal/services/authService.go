package services

import (
	"errors"

	"math/rand"
	"os"
	"time"

	"github.com/mackenzii/freemusic/helpers"
	middlewares "github.com/mackenzii/freemusic/internal/middleware"
	"github.com/mackenzii/freemusic/internal/models"
	"github.com/oklog/ulid/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// AuthService fournit des services d'authentification
type AuthService struct {
	DB                *gorm.DB
	GoogleOauthConfig *oauth2.Config
	ImageService      *ImageService
	EmailService      *EmailService
}

// NewAuthService crée une nouvelle instance de AuthService
func NewAuthService(db *gorm.DB, imageService *ImageService, emailService *EmailService) *AuthService {
	googleOauthConfig := &oauth2.Config{
		ClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectURL: os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:    google.Endpoint,
	}

	return &AuthService{
		DB:                db,
		GoogleOauthConfig: googleOauthConfig,
		ImageService:      imageService,
		EmailService:      emailService,
	}
}

// RegisterUser enregistre un nouvel utilisateur et envoie un email de confirmation
func (s *AuthService) RegisterUser(userInfo models.Users) (models.Users, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	newID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	// Hash le mot de passe
	hashedPassword, err := helpers.HashPassword(userInfo.PasswordHash)
	if err != nil {
		return models.Users{}, err
	}

	// Générer un jeton de confirmation unique
	confirmationToken := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// Créer l'utilisateur
	user := models.Users{
		ID:                newID.String(),
		Username:          userInfo.Username,
		Email:             userInfo.Email,
		PasswordHash:      hashedPassword,
		IsConfirmed:       false,
		ConfirmationToken: confirmationToken,
	}

	// Sauvegarder l'utilisateur dans la base de données
	if err := s.DB.Create(&user).Error; err != nil {
		return models.Users{}, err
	}

	// Envoyer un email de confirmation avec le jeton
	if err := s.EmailService.SendConfirmationEmail(user.Email, user.ConfirmationToken); err != nil {
		return models.Users{}, err
	}

	return user, nil
}

// Login authentifie un utilisateur et retourne un accessToken et un refreshToken
func (s *AuthService) Login(email, password string) (string, string, error) {
	var user models.Users
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", err
	}

	if !helpers.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", ErrInvalidCredentials
	}

	userID, err := ulid.Parse(user.ID)
	if err != nil {
		return "", "", err
	}

	// Générer le token d'accès et le refresh token
	accessToken, err := middlewares.GenerateToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := middlewares.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// UpdateRefreshToken met à jour le refresh token d'un utilisateur
func (s *AuthService) UpdateRefreshToken(email, refreshToken string) error {
	var user models.Users
	// Rechercher l'utilisateur par email
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	// Mettre à jour le refresh token
	user.RefreshToken = refreshToken
	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// GetUserByEmail récupère un utilisateur par son adresse email
func (s *AuthService) GetUserByEmail(email string) (models.Users, error) {
	var user models.Users
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

// Refresh génère un nouveau accessToken à partir d'un refreshToken valide
func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	claims, err := middlewares.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return "", "", errors.New("refresh token expired")
	}

	accessToken, err := middlewares.GenerateToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := middlewares.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// GetUserByID récupère un utilisateur par son ID
func (s *AuthService) GetUserByID(id string) (models.Users, error) {
	var user models.Users
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

// UpdateUser met à jour les informations de l'utilisateur
func (s *AuthService) UpdateUser(id string, userInfo models.Users) (models.Users, error) {
	var user models.Users
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return models.Users{}, err
	}

	// Mise à jour des champs de l'utilisateur
	user.Username = userInfo.Username
	user.Email = userInfo.Email
	user.ProfilePhoto = userInfo.ProfilePhoto
	user.FavoriteSport = userInfo.FavoriteSport
	user.Location = userInfo.Location
	user.Bio = userInfo.Bio

	// Mettre à jour le mot de passe si fourni
	if userInfo.PasswordHash != "" {
		hashedPassword, err := helpers.HashPassword(userInfo.PasswordHash)
		if err != nil {
			return models.Users{}, err
		}
		user.PasswordHash = hashedPassword
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

// DeleteUser supprime un utilisateur et ses données associées
func (s *AuthService) DeleteUser(id string) error {
	var user models.Users
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	// Supprimer les amis et demandes d'amis
	if err := s.DB.Where("sender_id = ? OR receiver_id = ?", id, id).Delete(&models.FriendRequest{}).Error; err != nil {
		return err
	}

	// Supprimer les concerts organisés par cet utilisateur
	if err := s.DB.Where("organizer_id = ?", id).Delete(&models.Event{}).Error; err != nil {
		return err
	}

	// Supprimer l'utilisateur
	if err := s.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

// GetAllUsers récupère tous les utilisateurs
func (s *AuthService) GetAllUsers() ([]models.Users, error) {
	var users []models.Users
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// SendConfirmationEmail utilise le service Email pour envoyer un lien de confirmation à l'utilisateur
func (s *AuthService) SendConfirmationEmail(email, confirmationToken string) error {
	return s.EmailService.SendConfirmationEmail(email, confirmationToken)
}

// Erreurs spécifiques
var ErrInvalidCredentials = errors.New("invalid credentials")

// GetPublicUserInfoByID retrieves public user information by ID

func (s *AuthService) GetPublicUserInfoByID(userID string) (*models.Users, error) {

	var user models.Users

	if err := s.DB.Select("username", "email", "location", "role").Where("id = ?", userID).First(&user).Error; err != nil {

		return nil, err

	}

	return &user, nil

}

func (s *AuthService) UpdateUserStatistics(userID string, matchesPlayed, matchesWon, goalsScored, behaviorScore int) error {

	return s.DB.Model(&models.Users{}).Where("id = ?", userID).Updates(models.Users{}).Error

}

func (s *AuthService) DeleteUserAndRelatedData(userID string) error {

	if err := s.DB.Where("id = ?", userID).Delete(&models.Users{}).Error; err != nil {

		return err

	}

	return nil

}
