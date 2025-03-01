package main

import (
	"form-server/internals/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, formHandler *handlers.FormHandler) {
	api := router.Group("/api")
	{
		api.POST("/forms", formHandler.SubmitForm)
		api.GET("/forms/:id", formHandler.GetFormByID)
		api.PUT("/forms/:id", formHandler.UpdateForm)
		api.DELETE("/forms/:id", formHandler.DeleteForm)
	}
}
