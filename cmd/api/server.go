package api

import (
	"context"
	"errors"
	"github.com/Roh-Bot/blog-api/internal/application"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/pkg/global"
	"github.com/Roh-Bot/blog-api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	"log"
	"net/http"
	"time"
)

type Server struct {
	//Config
	Config *config.AtomicConfig

	//Dependencies
	App       application.App
	Validator *validator.Validate
	Logger    logger.Logger
	AppCtx    *global.ApplicationContext

	// API Fields
	Router *echo.Echo
}

func NewServer(config *config.AtomicConfig, services application.App, validator *validator.Validate, logger logger.Logger, appCtx *global.ApplicationContext) *Server {
	return &Server{
		Config:    config,
		App:       services,
		Router:    echo.New(),
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
		if err := s.Router.Start(s.Config.Get().Server.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	return s.Router.Shutdown(ctx)
}

func (s *Server) registerSwagger() {
	s.Router.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Server) registerHandlers() {
	// Create the main API group
	apiGroup := s.Router.Group("/api")

	// Apply global middlewares to the API group
	apiGroup.Use(s.httpLogger)

	apiGroup.GET("/health", s.Health)

	// Authentication routes
	authGroup := apiGroup.Group("/authentication")
	authGroup.POST("/login", s.authLoginUser)

	// Protected routes
	protectedGroup := apiGroup.Group("")
	protectedGroup.Use(s.validateAuth)
	protectedGroup.GET("/blog-post", s.postsGet)
	protectedGroup.GET("/blog-post/:id", s.postsGet)
	protectedGroup.POST("/blog-post/:id", s.postAdd)
	protectedGroup.POST("/blog-post/:id", s.postUpdate)
	protectedGroup.POST("/blog-post/:id", s.postDelete)
}

func (s *Server) registerMiddlewares() {
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.CORS())
}
