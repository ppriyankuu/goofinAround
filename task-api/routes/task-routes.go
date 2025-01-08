package routes

import (
	"net/http"
	"strconv"
	"task-api/database"
	"task-api/models"

	"github.com/gin-gonic/gin"
)

func TaskRoutes(r *gin.Engine) {
	r.GET("/tasks", GetTasks)
	r.POST("/tasks", CreateTask)
	r.GET("/tasks/:id", GetTaskByID)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)
}

func GetTasks(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10
	offset := (page - 1) * limit

	var tasks []models.Task
	query := "SELECT * FROM tasks WHERE ($1 = '' OR status = $1) LIMIT $2 OFFSET $3"
	err := database.DB.Select(&tasks, query, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO tasks (title, description, status, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id"
	err := database.DB.QueryRow(query, task.Title, task.Description, task.Status).Scan(&task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	query := "SELECT * FROM tasks WHERE id = $1"
	err := database.DB.Get(&task, query, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4"
	_, err := database.DB.Exec(query, task.Title, task.Description, task.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM tasks WHERE id = $1"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
