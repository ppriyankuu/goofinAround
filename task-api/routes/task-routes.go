package routes

import (
	"task-api/config"
	"task-api/controllers"
	"task-api/middlewares"

	"github.com/gin-gonic/gin"
)

func TaskRoutes(r *gin.Engine) {
	tasks := r.Group("/tasks")
	tasks.Use(middlewares.AuthMiddleware(config.SecretKey))
	{
		tasks.GET("/tasks", controllers.GetTasks)
		tasks.POST("/tasks", controllers.CreateTask)
		tasks.GET("/tasks/:id", controllers.GetTaskByID)
		tasks.PUT("/tasks/:id", controllers.UpdateTask)
		tasks.DELETE("/tasks/:id", controllers.DeleteTask)
	}
}
