package api

import (
	"context"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/config"
	servicesv1 "github.com/Roh-Bot/blog-api/internal/services"
	"github.com/Roh-Bot/blog-api/pkg/global"
	"github.com/Roh-Bot/blog-api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"log"
	"net/http"
	"time"
)

type Server struct {
	//Config
	Config *config.AtomicConfig

	//Dependencies
	Services  *servicesv1.Service
	Validator *validator.Validate
	Logger    logger.Logger
	AppCtx    *global.ApplicationContext

	// API Fields
	Router *fiber.App
}

func NewServer(config *config.AtomicConfig, services *servicesv1.Service, validator *validator.Validate, logger logger.Logger, appCtx *global.ApplicationContext) *Server {
	return &Server{
		Config:   config,
		Services: services,
		Router: fiber.New(fiber.Config{
			BodyLimit:                512,
			EnableSplittingOnParsers: false,
		}),
		Validator: validator,
		Logger:    logger,
		AppCtx:    appCtx,
	}
}

func (s *Server) Run() {
	defer s.AppCtx.Done()
	go func() {
		s.registerMiddlewares()
		s.registerSwagger()
		s.registerHandlers()
		if err := s.Router.Listen(s.Config.Get().Server.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	for range s.AppCtx.Context().Done() {
		if err := s.Shutdown(); err != nil {
			log.Println("An error occurred while shutting down the server")
		}
		break
	}
	log.Println("Server shutdown completed")
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return s.Router.ShutdownWithContext(ctx)
}

func (s *Server) registerSwagger() {
	s.Router.Get("/swagger/*", swagger.HandlerDefault)
}

func (s *Server) registerHandlers() {
	// Create the main API group
	apiGroup := s.Router.Group("/api")

	// Apply global middlewares to the API group
	apiGroup.Use(recover.New())
	apiGroup.Use(s.requestLogger)
	apiGroup.Use(s.responseLogger)

	apiGroup.Get("/health", s.Health)

	// Authentication routes
	authGroup := apiGroup.Group("/authentication")
	authGroup.Post("/login", s.authLoginUser)

	// Protected routes
	protectedGroup := apiGroup.Group("")
	protectedGroup.Use(s.validateAuth)
	protectedGroup.Get("/blog-post", s.postsGet)
	protectedGroup.Get("/blog-post/:id", s.postsGet)
	protectedGroup.Post("/blog-post/:id", s.postAdd)
	protectedGroup.Patch("/blog-post/:id", s.postUpdate)
	protectedGroup.Delete("/blog-post/:id", s.postDelete)
}

func (s *Server) registerMiddlewares() {
	s.Router.Use(cors.New(cors.Config{
		AllowOrigins: "",
	}))
}
