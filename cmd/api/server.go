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
			RequestMethods:           []string{fiber.MethodGet, fiber.MethodDelete, fiber.MethodOptions, fiber.MethodPost, fiber.MethodPut},
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

func (s *Server) registerHandlers() {
	apiGroup := s.Router.Group("api/")

	// Middlewares should be applied before routes
	apiGroup.Use(recover.New())
	apiGroup.Use(s.requestLogger)
	apiGroup.Use(s.responseLogger)

	// Public routes
	public := apiGroup.Group("")
	public.Get("health", s.Health)

	authentication := public.Group("authentication/")
	authentication.Post("login", s.authLoginUser)

	// Protected routes
	protected := apiGroup.Group("")
	protected.Use(s.validateAuth)
	protected.Get("blog-post", s.postsGet)
	protected.Get("blog-post/:id", s.postsGet)
	protected.Post("blog-post/:id", s.postAdd)
	protected.Patch("blog-post/:id", s.postUpdate)
	protected.Delete("blog-post/:id", s.postDelete)
}

func (s *Server) registerMiddlewares() {
	s.Router.Use(cors.New(cors.Config{
		AllowOrigins: "",
	}))
}
