package main

import (
	"github.com/AnatolyKoltun/go-calculator-service/config"
	"github.com/AnatolyKoltun/go-calculator-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Connect()
	defer config.Close()

	router := gin.Default()

	router.POST("/calculate", handlers.Calculate)
	router.GET("/calculations", handlers.GetCalculations)

	router.Run(":8080")
}
