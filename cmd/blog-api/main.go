package main

import (
	"github.com/Roh-Bot/blog-api/cmd/api"
	_ "github.com/Roh-Bot/blog-api/docs"
	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/Roh-Bot/blog-api/internal/database"
	servicesv1 "github.com/Roh-Bot/blog-api/internal/services"
	"github.com/Roh-Bot/blog-api/internal/store"
	"github.com/Roh-Bot/blog-api/internal/store/cache"
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

// @host
// @BasePath /api
// @schemes https http

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

	//Initializing logger
	newLogger, err := logger.ZapNew(cfg.Get().Logger)
	if err != nil {
		log.Fatal(err)
	}
	// Flushing the logger
	defer func() {
		err := newLogger.Sync()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	// Connecting to master database
	db, err := database.NewMasterConnection(cfg.Get().Database)
	if err != nil {
		log.Fatal(err)
	}
	//Flushing database connection pool
	defer db.Flush()

	// Connection to cache database
	dbCache, err := database.NewCache(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing storage layer
	newStore := store.NewStorage(db, cfg)

	// Initializing cache layer
	newCache := cache.NewCache(dbCache, cfg)

	// Initializing Authenticators
	jwt := auth.NewJWTAuthenticator(cfg, newStore)
	aes, err := auth.NewAES(cfg.Get().Auth.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	auth2 := auth.NewAuthentication(jwt, aes)

	// Initializing Service layer
	services := servicesv1.NewService(cfg, auth2, newStore, newCache, newLogger)

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

	log.Println("Goodbye")
}
