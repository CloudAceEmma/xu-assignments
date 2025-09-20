package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "todoapp-go/backend/docs" // Import generated docs
)

// @title Todo Application API
// @version 1.0
// @description This is a sample server for a todo application.
// @host localhost:8000
// @BasePath /api/v1

// Todo model
type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	Priority    int       `json:"priority" gorm:"default:0"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

var DB *gorm.DB

func main() {
	// Database connection
	dsn := "monty:test001@tcp(127.0.0.1:3306)/todoapp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthCheck)
		v1.GET("/todos", getTodos)
		v1.POST("/todos", createTodo)
		v1.PUT("/todos/:id", updateTodo)
		v1.DELETE("/todos/:id", deleteTodo)
		v1.PATCH("/todos/:id/toggle", toggleTodoStatus)
		v1.DELETE("/todos/completed", deleteCompletedTodos)
		v1.DELETE("/todos/all", deleteAllTodos)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// healthCheck godoc
// @Summary Show the status of the server
// @Description get the status of the server
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

// getTodos godoc
// @Summary Get all todos
// @Description Get a list of all todos, with optional filtering by completion status and pagination.
// @Tags todos
// @Accept json
// @Produce json
// @Param completed query bool false "Filter by completion status"
// @Param limit query int false "Limit the number of results" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{} "Success response"
// @Router /todos [get]
func getTodos(c *gin.Context) {
	var todos []Todo
	query := DB.Model(&Todo{})

	if completedStr := c.Query("completed"); completedStr != "" {
		completed := completedStr == "true"
		query = query.Where("completed = ?", completed)
	}

	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")

	query = query.Limit(toInt(limit)).Offset(toInt(offset))

	result := query.Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to fetch todos", "detail": result.Error.Error()})
		return
	}

	var total int64
	query.Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    todos,
		"total":   total,
	})
}

// createTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item with title, description, priority, and optional due date.
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body Todo true "Todo object to be created"
// @Success 201 {object} map[string]interface{} "Todo created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Router /todos [post]
func createTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid input", "detail": err.Error()})
		return
	}

	result := DB.Create(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to create todo", "detail": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "Todo created successfully",
		"data":    todo,
	})
}

// updateTodo godoc
// @Summary Update an existing todo
// @Description Update an existing todo item by its ID.
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body Todo true "Updated todo object"
// @Success 200 {object} map[string]interface{} "Todo updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Router /todos/{id} [put]
func updateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if result := DB.First(&todo, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Todo not found", "detail": result.Error.Error()})
		return
	}

	var updatedTodo Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid input", "detail": err.Error()})
		return
	}

	// Update only the fields that are provided in the request body
	result := DB.Model(&todo).Updates(updatedTodo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to update todo", "detail": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Todo updated successfully",
		"data":    todo,
	})
}

// deleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo item by its ID.
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} map[string]interface{} "Todo deleted successfully"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Router /todos/{id} [delete]
func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if result := DB.First(&todo, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Todo not found", "detail": result.Error.Error()})
		return
	}

	result := DB.Delete(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to delete todo", "detail": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Todo deleted successfully",
	})
}

// toggleTodoStatus godoc
// @Summary Toggle todo completion status
// @Description Toggle the completion status of a todo item by its ID.
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} map[string]interface{} "Todo status toggled successfully"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Router /todos/{id}/toggle [patch]
func toggleTodoStatus(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if result := DB.First(&todo, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "Todo not found", "detail": result.Error.Error()})
		return
	}

	todo.Completed = !todo.Completed
	result := DB.Save(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to toggle todo status", "detail": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Todo status toggled successfully",
		"data": gin.H{
			"id":        todo.ID,
			"completed": todo.Completed,
			"updated_at": todo.UpdatedAt,
		},
	})
}

// deleteCompletedTodos godoc
// @Summary Delete all completed todos
// @Description Delete all todo items that are marked as completed.
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Completed todos deleted successfully"
// @Router /todos/completed [delete]
func deleteCompletedTodos(c *gin.Context) {
	result := DB.Where("completed = ?", true).Delete(&Todo{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to delete completed todos", "detail": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"message":     "Completed todos deleted successfully",
		"deleted_count": result.RowsAffected,
	})
}

// deleteAllTodos godoc
// @Summary Delete all todos
// @Description Delete all todo items from the database.
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "All todos deleted successfully"
// @Router /todos/all [delete]
func deleteAllTodos(c *gin.Context) {
	result := DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Todo{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Failed to delete all todos", "detail": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"message":     "All todos deleted successfully",
		"deleted_count": result.RowsAffected,
	})
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
