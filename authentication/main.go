package main

import (
	"authentication/config"
	"authentication/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.InitDB()

	routes.AuthRoutes(r)

	r.Run(":8080")
}
