package app

import (
	"cms-octo-chat-api/handler"
	"cms-octo-chat-api/helper"
	"cms-octo-chat-api/middleware"
	"cms-octo-chat-api/model"
	"encoding/json"

	"github.com/gofiber/swagger"

	_ "cms-octo-chat-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mikhail-bigun/fiberlogrus"
)

func ApiRouter(resources model.Resources) *fiber.App {

	f := fiber.New(fiber.Config{

		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	f.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "*",
	}))

	m := middleware.NewBaseMiddleware(middleware.BaseMiddleware{
		Env:  resources.Env,
		Logs: resources.Logs,
		DB:   resources.DB,
	})
	mBodyValidator, _ := m.NewBodyValidatorMiddleware("validation_messages.yaml")

	h := handler.NewBaseHandler(handler.BaseHandler{
		Env:    resources.Env,
		Logs:   resources.Logs,
		DB:     resources.DB,
		OpenAI: resources.OpenAI,
	})

	f.Use(
		fiberlogrus.New(
			fiberlogrus.Config{
				Logger: helper.MakeLogger(
					helper.Setup{
						Env:     resources.Env.LogEnv,
						Logname: resources.Env.LogPath + "/access_log",
						Display: true,
						Level:   resources.Env.LogLevel,
					}),
			}))

	f.Get("/swagger/*", swagger.HandlerDefault)

	v1 := f.Group("/v1")

	//Auth
	v1.Post("/login", mBodyValidator.Validate(&model.AuthRequest{}), h.Login)

	// users
	user := v1.Group("/user", m.JwtAuthMiddleware())
	user.Post("/", m.PermissionChecker("C"), mBodyValidator.Validate(&model.UserCreateRequest{}), h.CreateUser)
	user.Post("/:id/matrix", m.PermissionChecker("R"), mBodyValidator.Validate(&model.UserMatrixRequest{}), h.CreateUserMatrix)

	// conversation
	conversation := v1.Group("/conversation", m.JwtAuthMiddleware())
	conversation.Get("/", h.ListConversation)
	conversation.Post("/", mBodyValidator.Validate(&model.ConversationCreateReq{}), h.CreateConversation)

	//chat
	chat := v1.Group("/chat", m.JwtAuthMiddleware())
	chat.Post("/", mBodyValidator.Validate(&model.ChatRequest{}), h.SingleChat)
	chat.Post("/stream", mBodyValidator.Validate(&model.ChatRequest{}), h.StreamChat)

	// message
	message := v1.Group("/message", m.JwtAuthMiddleware())
	message.Get("/:conversation_pid", h.ListMessage)

	return f
}
