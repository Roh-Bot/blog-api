package main

import (
	"github.com/Roh-Bot/blog-api/cmd/api"
	_ "github.com/Roh-Bot/blog-api/docs"
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/database"
	servicesv1 "github.com/Roh-Bot/blog-api/internal/services"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/validator"
	"github.com/Roh-Bot/blog-api/pkg/global"
	"github.com/Roh-Bot/blog-api/pkg/logger"
	"log"
)

// @title Blog API
// @version 1.0
// @description Blog API with JWT authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token
func main() {
	//Parsing global flags
	global.ParseFlags()

	//Initializing Application Context
	appCtx := global.NewApplicationContext()

	//Loading configuration
	cfg, err := config.LoadConfiguration(appCtx.Context())
	if err != nil {
		log.Fatal(err)
	}

	rotateLogsWriteSyncer, err := logger.RotateLogsWriteSyncer(*cfg.Get())
	if err != nil {
		log.Fatal(err)
	}
	rotateLogsSink := logger.NewZapCore(rotateLogsWriteSyncer, cfg.Get().Logger.Level)

	//Initializing logger
	newLogger, err := logger.ZapNew(cfg.Get().Logger, rotateLogsSink)
	if err != nil {
		log.Fatal(err)
	}

	// Connecting to database
	db, err := database.New(cfg.Get().Database)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing storage layer
	newStore := store.NewStorage(db, cfg)

	// Initializing Authenticators
	jwt := auth.NewJWTAuthenticator(cfg, newStore)
	aes, err := auth.NewAES(cfg.Get().Auth.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	auth2 := auth.NewAuthentication(jwt, aes)

	// Initializing Service layer
	services := servicesv1.NewService(newStore, newLogger, cfg, auth2)

	// Initializing validator
	validator2 := validator.NewValidator()

	// Initializing server
	server := api.NewServer(cfg, services, validator2, newLogger, appCtx)

	// Firing up the server
	appCtx.Add(1)
	go server.Run()

	// Waiting for termination signal
	appCtx.HandleShutdownSignal()
	appCtx.WaitForShutdown()

	//Close db connection, flush logger, close smtp
	//database.Flush()

	log.Println("Goodbye")
}
