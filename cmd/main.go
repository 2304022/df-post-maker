package main

import (
	"df-post-maker/internal/controller"
	"df-post-maker/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	testUC := usecase.NewTestUseCase()
	testCtrl := controller.NewTestController(testUC)

	directFarmUC := usecase.NewDirectFarmUseCase()
	directFarmCtrl := controller.NewDirectFarmController(directFarmUC)

	router := gin.Default()
	testCtrl.RegisterRoutes(router)
	directFarmCtrl.RegisterRoutes(router)

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
