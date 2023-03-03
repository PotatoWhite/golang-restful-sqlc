package main

import (
	"github.com/gin-gonic/gin"
	"github.com/potatowhite/restfulapi/cmd/config"
	"github.com/potatowhite/restfulapi/pkg/database"
	"github.com/potatowhite/restfulapi/pkg/service"
	"log"
)

func main() {
	// Read configuration
	cfg := initConfig()

	// Instantiate the database
	postgres := initDatabase(cfg)

	// Instantiate data access layer
	queries := initDataAccessLayer(postgres)

	// Instantiate the service
	authorService := initService(queries)

	// Register our service handlers to the router
	router := initServer(authorService)

	// Start the server
	err := router.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Server started")
}

func initServer(authorService *service.Service) *gin.Engine {
	router := gin.Default()
	authorService.RegisterHandlers(router)
	return router
}

func initService(queries *database.Queries) *service.Service {
	authorService := service.NewService(queries)
	log.Printf("Author service instantiated")
	return authorService
}

func initDataAccessLayer(postgres *database.Postgres) *database.Queries {
	queries := database.New(postgres.DB)
	log.Printf("Data access layer instantiated")
	return queries
}

func initDatabase(cfg *config.Config) *database.Postgres {
	postgres, err := database.NewPostgres(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Database connection established")
	return postgres
}

func initConfig() *config.Config {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Configuration loaded")
	return cfg
}
