package main

import (
	api "APP4/api/handlers"
	"APP4/api/repository"
	"APP4/database/db"
	"APP4/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error in database: %v", err.Error())
		return
	}

	repo := repository.NewRepository(db)

	router := gin.Default()
	authCtrl := api.NewAuthHandlers(repo)
	routes.SetUpRoutes(router, authCtrl)

	router.Run(":8000")

}
