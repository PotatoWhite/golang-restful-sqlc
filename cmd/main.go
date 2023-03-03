package main

import (
	"github.com/gin-gonic/gin"
	"github.com/potatowhite/restfulapi/cmd/config"
	"github.com/potatowhite/restfulapi/pkg/database"
	"github.com/potatowhite/restfulapi/pkg/microservice/authors"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	cfg := loadConfig()
	db := connectDatabase(cfg)
	queries := initQueries(db)
	authorService := initAuthorService(queries)
	handler := initAuthorHandler(authorService)
	server := initServer(handler)

	err := server.Run(":" + cfg.Server.Port)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func loadConfig() *config.Config {
	logger.Println("Loading configuration...")
	cfg, err := config.Read()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %s", err.Error())
	}
	logger.Printf("Loaded configuration: %+v", cfg)
	return cfg
}

func connectDatabase(cfg *config.Config) *database.Postgres {
	logger.Println("Connecting to database...")
	db, err := database.NewPostgres(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %s", err.Error())
	}
	return db
}

func initQueries(db *database.Postgres) *database.Queries {
	logger.Println("Initializing queries...")
	queries := database.New(db.DB)
	return queries
}

func initAuthorService(queries *database.Queries) authors.AuthorService {
	logger.Println("Initializing author service...")
	return authors.NewAuthorService(queries)
}

func initAuthorHandler(authorService authors.AuthorService) authors.AuthorHandler {
	logger.Println("Initializing author handler...")
	return authors.NewAuthorHandler(authorService)
}

func initServer(handler authors.AuthorHandler) *gin.Engine {
	logger.Println("Initializing server...")
	router := gin.Default()
	handler.RegisterHandlers(router)
	return router
}
