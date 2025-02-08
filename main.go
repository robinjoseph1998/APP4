package main

import (
	api "APP4/api/handlers"
	"APP4/api/repository"
	"APP4/database/db"
	"APP4/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Now environment variables will be available
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	if clientID == "" {
		log.Fatal("INSTAGRAM_CLIENT_ID is not set")
	}

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error in database: %v", err.Error())
		return
	}

	repo := repository.NewRepository(db)

	router := gin.Default()
	OauthInstagramCtrl := api.NewOAuthInstagramHandlers(repo)
	OauthTwitterCtrl := api.NewOAuthTwitterHandlers(repo)
	routes.SetUpRoutes(router, OauthInstagramCtrl, OauthTwitterCtrl)

	router.Run(":8000")

}
