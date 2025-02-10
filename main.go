package main

import (
	api "APP4/api/handlers"
	"APP4/api/repository"
	services "APP4/api/services/twitter"
	"APP4/database/db"
	"APP4/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err.Error())
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error in database: %v", err.Error())
		return
	}

	router := gin.Default()
	repo := repository.NewRepository(db)

	twitterServices := services.NewTwitterServices(repo)
	OauthTwitterCtrl := api.NewOAuthTwitterHandlers(repo, twitterServices)

	OauthInstagramCtrl := api.NewOAuthInstagramHandlers(repo)
	authCommonCtrl := api.NewCommonAuthHandlers(repo)

	routes.SetUpRoutes(router, OauthInstagramCtrl, OauthTwitterCtrl, authCommonCtrl)

	router.Run(":8000")

}
