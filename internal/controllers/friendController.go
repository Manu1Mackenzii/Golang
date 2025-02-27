package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mackenzii/freemusic/internal/services"
)

type FriendController struct {
	FriendService       *services.FriendService
	NotificationService *services.NotificationService
}

func NewFriendController(friendService *services.FriendService, notificationService *services.NotificationService) *FriendController {
	return &FriendController{
		FriendService:       friendService,
		NotificationService: notificationService,
	}
}

func (fc *FriendController) SendFriendRequest(c *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"sender_id"`
		ReceiverId string `json:"receiver_id"`
	}
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// Vérifiez que sender_id et receiver_id ne sont pas vides
	if request.SenderId == "" || request.ReceiverId == "" {
		log.Printf("SenderId or ReceiverId is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sender_id and receiver_id cannot be empty",
		})
	}

	// Vérifiez que sender_id et receiver_id ne sont pas identiques
	if request.SenderId == request.ReceiverId {
		log.Printf("SenderId and ReceiverId are the same")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sender_id and receiver_id cannot be the same",
		})
	}

	if err := fc.FriendService.SendFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		log.Printf("Error sending friend request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "friend request sent successfully",
	})
}

func (fc *FriendController) AcceptFriendRequest(c *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"sender_id"`
		ReceiverId string `json:"receiver_id"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := fc.FriendService.AcceptFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "friend request accepted successfully",
	})
}

func (fc *FriendController) DeclineFriendRequest(c *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"sender_id"`
		ReceiverId string `json:"receiver_id"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := fc.FriendService.DeclineFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "friend request declined successfully",
	})
}

func (fc *FriendController) GetFriendRequests(c *fiber.Ctx) error {
	userID := c.Params("userID")
	friendRequests, err := fc.FriendService.GetFriendRequests(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(friendRequests)
}

func (fc *FriendController) GetFriends(c *fiber.Ctx) error {
	userID := c.Params("userID")
	friends, err := fc.FriendService.GetFriends(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(friends)
}

func (fc *FriendController) SearchUsersByUsername(c *fiber.Ctx) error {
	username := c.Query("username")
	users, err := fc.FriendService.SearchUsersByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(users)
}
