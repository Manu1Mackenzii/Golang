package server

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/mackenzii/freemusic/internal/controllers"
	"github.com/mackenzii/freemusic/internal/models"
	"github.com/mackenzii/freemusic/internal/routes"
	"github.com/mackenzii/freemusic/internal/services"
	"github.com/mackenzii/freemusic/internal/storage"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Run() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found", err)
	}

	// Database connection
	db, err := storage.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Table migration
	if err := db.AutoMigrate(&models.Users{}, &models.Event{}, &models.FriendRequest{}, &models.Message{}); err != nil {
		log.Printf("Error migrating database: %v", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Create a broadcast channel for notifications
	notificationBroadcast := make(chan services.Notification)

	// Initialize services and controllers
	imageService := services.NewImageService("./uploads")
	emailService := services.NewEmailService()
	authService := services.NewAuthService(db, imageService, emailService)
	webSocketService := services.NewWebSocketService()
	notificationService := services.NewNotificationService(db, redisClient, notificationBroadcast, webSocketService)
	openAIService := services.NewOpenAIService()
	friendChatService := services.NewFriendChatService(db, webSocketService)
	categoryService := services.NewCategoryService(db)
	eventService := services.NewEventService(db)

	friendService := services.NewFriendService(db, authService, webSocketService)
	friendController := controllers.NewFriendController(friendService, notificationService)
	// chatService := services.NewChatService(db, redisClient)
	// matchController := controllers.NewMatchController(matchService, authService, db, chatService, redisClient)
	matchPlayersService := services.NewMatchPlayersService(db)
	// matchPlayersController := controllers.NewMatchPlayersController(matchPlayersService, authService, db)
	// chatController := controllers.NewChatController(chatService)
	openAiController := controllers.NewOpenAiController(openAIService, matchPlayersService)
	authController := controllers.NewAuthController(authService, imageService)
	friendChatController := controllers.NewfriendChatController(friendChatService, friendService)
	categoryController := controllers.NewCategoryController(categoryService, authService, db, redisClient)
	eventController := controllers.NewEventController(eventService, authService, db, redisClient)

	// Configure Fiber app
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Apply rate limiter middleware to all routes except Swagger
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/swagger/*" {
			return c.Next()
		}
		return limiter.New(limiter.Config{
			Max:        10,
			Expiration: 30 * time.Second,
		})(c)
	})

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to TeamUp API!")
	})
	routes.SetupRoutesAuth(app, authController)
	routes.SetupRoutesCategories(app, categoryController)
	routes.SetupOpenAiRoutes(app, openAiController)
	routes.SetupFriendRoutes(app, friendController)
	routes.SetupRoutesFriendMessage(app, friendChatController)
	routes.SetupRoutesEvents(app, eventController)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		webSocketService.HandleWebSocket(c)
	}))

	// Start WebSocket broadcast
	go webSocketService.StartBroadcast()

	// Start listening for notifications
	go notificationService.ListenForNotifications()

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("Server started on port %s", port)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			// if err := matchService.UpdateMatchStatuses(); err != nil {
			// 	log.Printf("Erreur lors de la mise à jour des statuts des matchs : %v", err)
			// }
		}
	}()

	log.Fatal(app.Listen(":" + port))
}
