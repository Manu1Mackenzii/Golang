package routes

import (
	"github.com/mackenzii/freemusic/internal/controllers"
	middlewares "github.com/mackenzii/freemusic/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRoutesAuth configure les routes pour l'authentification des utilisateurs.
func SetupRoutesAuth(app *fiber.App, controller *controllers.AuthController) {
	api := app.Group("/api")

	// Routes ouvertes (sans authentification)
	api.Post("/register", controller.RegisterHandler)
	api.Post("/login", controller.LoginHandler)
	api.Post("/refresh", controller.RefreshHandler)
	api.Get("/auth/google", controller.GoogleLogin)
	api.Get("/auth/google/callback", controller.GoogleCallback)
	api.Get("/confirm_email", controller.ConfirmEmailHandler)

	// Routes protégées (requièrent authentification)
	api.Use(middlewares.JWTMiddleware)
	api.Put("/userUpdate", controller.UserUpdate)
	// api.Get("/userInfo", controller.GetUserInfoHandler)
	api.Get("/users", controller.GetUsersHandler)
	api.Get("/users/:id/public", controller.GetPublicUserInfoHandler)
	api.Delete("/deleteMyAccount", controller.DeleteUserHandler)
	// api.Post("/UpdateUserStatistics", controller.UpdateUserStatistics)
}

// SetupRoutesEvents configure les routes pour gérer les événements.
func SetupRoutesEvents(app *fiber.App, controller *controllers.EventController) {
	api := app.Group("/api/events")
	api.Use(middlewares.JWTMiddleware)

	api.Get("/", controller.GetAllEvents)             // Récupérer tous les événements
	api.Post("/createEvent/", controller.CreateEvent) // Créer un nouvel événement
	api.Delete("/event/:id", controller.DeleteEvent)  // Supprimer un événement
	api.Get("/:event_id", controller.GetEventByID)
	api.Put("/:id", controller.UpdateEvent)
}

// SetupRoutesCategories configure les routes pour gérer les catégories.
func SetupRoutesCategories(app *fiber.App, controller *controllers.CategoryController) {
	api := app.Group("/api/categories")
	api.Use(middlewares.JWTMiddleware)

	api.Get("/", controller.GetCategories)
	api.Get("/:id", controller.GetCategory)
	api.Post("/createCategory", controller.CreateCategory)
	api.Delete("/:id", controller.DeleteCategory)
	api.Put("/:id", controller.UpdateCategory)
}

// SetupFriendRoutes configure les routes pour gérer les relations d'amis.
func SetupFriendRoutes(app *fiber.App, friendController *controllers.FriendController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)

	api.Post("/friend/send", friendController.SendFriendRequest)            // Envoyer une demande d'ami
	api.Post("/friend/accept", friendController.AcceptFriendRequest)        // Accepter une demande d'ami
	api.Post("/friend/decline", friendController.DeclineFriendRequest)      // Refuser une demande d'ami
	api.Get("/friend/requests/:userID", friendController.GetFriendRequests) // Obtenir les demandes d'amis
	api.Get("/friend/:userID", friendController.GetFriends)                 // Récupérer les amis d'un utilisateur
	api.Get("/friend/search", friendController.SearchUsersByUsername)       // Rechercher des utilisateurs par nom
}

// SetupRoutesFriendMessage configure les routes pour gérer les messages entre amis.
func SetupRoutesFriendMessage(app *fiber.App, friendChatController *controllers.FriendChatController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)

	api.Post("/message/send", friendChatController.SendMessage)                          // Envoyer un message
	api.Get("/message/messages/:senderID/:receiverID", friendChatController.GetMessages) // Obtenir les messages entre deux utilisateurs
}

// SetupRoutesWebSocket configure les routes pour les WebSocket.
func SetupRoutesWebSocket(app *fiber.App, controller *controllers.WebSocketController) {
	api := app.Group("/api")
	api.Get("/updates", websocket.New(controller.WebSocketHandler)) // WebSocket pour les mises à jour en temps réel
}

// SetupOpenAiRoutes configure les routes pour utiliser les services OpenAI.
func SetupOpenAiRoutes(app *fiber.App, controller *controllers.OpenAiController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)

	// api.Get("/openai/suggestions/:event_id", controller.SuggestConcerts)
}

// SetupRoutes configure toutes les routes de l'application.
func SetupRoutes(app *fiber.App, authController *controllers.AuthController, eventController *controllers.EventController, categoryController *controllers.CategoryController, friendController *controllers.FriendController, friendChatController *controllers.FriendChatController, wsController *controllers.WebSocketController, aiController *controllers.OpenAiController) {
	SetupRoutesAuth(app, authController)
	SetupRoutesEvents(app, eventController)
	SetupRoutesCategories(app, categoryController)
	SetupFriendRoutes(app, friendController)
	SetupRoutesFriendMessage(app, friendChatController)
	SetupRoutesWebSocket(app, wsController)
	SetupOpenAiRoutes(app, aiController)

}
