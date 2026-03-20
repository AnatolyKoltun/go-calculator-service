package main

import (
	"github.com/AnatolyKoltun/go-calculator-service/config"
	"github.com/AnatolyKoltun/go-calculator-service/database"
	"github.com/AnatolyKoltun/go-calculator-service/handlers"
	"github.com/gin-gonic/gin"
)

func setupAndRunServer() {
	router := gin.Default()

	router.POST("/calculate", handlers.Calculate)
	router.GET("/calculations", handlers.GetCalculations)

	router.Run(":8080")
}

func connectToDB() {
	dsn := new(config.DataSourceName)
	dsn.GetDatabaseURL()

	database.Connect(dsn.DatabaseURL)
}

func main() {
	defer database.Close()

	connectToDB()
	setupAndRunServer()
}
