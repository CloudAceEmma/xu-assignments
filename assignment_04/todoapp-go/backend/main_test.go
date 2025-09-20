package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	// Migrate the schema
	db.AutoMigrate(&Todo{})
	return db
}

func setupRouter(testDB *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	DB = testDB // Use the test database

	// CORS middleware (simplified for testing)
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
	return router
}

func TestHealthCheck(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "UP")
}

func TestCreateTodo(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	newTodo := Todo{Title: "Test Todo", Description: "Description for test todo", Priority: 1}
	jsonValue, _ := json.Marshal(newTodo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusCreated), response["code"])
	assert.Equal(t, "Todo created successfully", response["message"])
	assert.NotNil(t, response["data"].(map[string]interface{})["id"])

	// Verify todo is in DB
	var createdTodo Todo
	id := uint(response["data"].(map[string]interface{})["id"].(float64))
	db.First(&createdTodo, id)
	assert.Equal(t, newTodo.Title, createdTodo.Title)
}

func TestGetTodos(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create some todos
	db.Create(&Todo{Title: "Todo 1", Completed: false})
	db.Create(&Todo{Title: "Todo 2", Completed: true})

	// Test get all todos
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/todos", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.Len(t, response["data"].([]interface{}), 2)
	assert.Equal(t, float64(2), response["total"])

	// Test get incomplete todos
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/todos?completed=false", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response["data"].([]interface{}), 1)
	assert.Equal(t, float64(1), response["total"])

	// Test get completed todos
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/todos?completed=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response["data"].([]interface{}), 1)
	assert.Equal(t, float64(1), response["total"])
}

func TestUpdateTodo(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create a todo
	initialTodo := Todo{Title: "Initial Todo", Completed: false}
	db.Create(&initialTodo)

	updatedFields := Todo{Title: "Updated Todo", Completed: true, Priority: 5}
	jsonValue, _ := json.Marshal(updatedFields)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/todos/"+strconv.Itoa(int(initialTodo.ID)), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "Todo updated successfully", response["message"])

	// Verify todo is updated in DB
	var updatedTodoInDB Todo
	db.First(&updatedTodoInDB, initialTodo.ID)
	assert.Equal(t, updatedFields.Title, updatedTodoInDB.Title)
	assert.Equal(t, updatedFields.Completed, updatedTodoInDB.Completed)
	assert.Equal(t, updatedFields.Priority, updatedTodoInDB.Priority)
}

func TestDeleteTodo(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create a todo
	initialTodo := Todo{Title: "Todo to delete"}
	db.Create(&initialTodo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/todos/"+strconv.Itoa(int(initialTodo.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "Todo deleted successfully", response["message"])

	// Verify todo is deleted from DB
	var deletedTodo Todo
	result := db.First(&deletedTodo, initialTodo.ID)
	assert.Error(t, result.Error)
	assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func TestToggleTodoStatus(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create a todo
	initialTodo := Todo{Title: "Todo to toggle", Completed: false}
	db.Create(&initialTodo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/todos/"+strconv.Itoa(int(initialTodo.ID))+"/toggle", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "Todo status toggled successfully", response["message"])
	assert.Equal(t, true, response["data"].(map[string]interface{})["completed"])

	// Verify status is toggled in DB
	var toggledTodo Todo
	db.First(&toggledTodo, initialTodo.ID)
	assert.True(t, toggledTodo.Completed)

	// Toggle again
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/api/v1/todos/"+strconv.Itoa(int(initialTodo.ID))+"/toggle", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["data"].(map[string]interface{})["completed"])

	var toggledTodoAgain Todo
	db.First(&toggledTodoAgain, initialTodo.ID)
	assert.False(t, toggledTodoAgain.Completed)
}

func TestDeleteCompletedTodos(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create some todos
	db.Create(&Todo{Title: "Todo 1", Completed: true})
	db.Create(&Todo{Title: "Todo 2", Completed: false})
	db.Create(&Todo{Title: "Todo 3", Completed: true})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/todos/completed", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "Completed todos deleted successfully", response["message"])
	assert.Equal(t, float64(2), response["deleted_count"])

	// Verify only incomplete todos remain
	var remainingTodos []Todo
	db.Find(&remainingTodos)
	assert.Len(t, remainingTodos, 1)
	assert.Equal(t, "Todo 2", remainingTodos[0].Title)
}

func TestDeleteAllTodos(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create some todos
	db.Create(&Todo{Title: "Todo 1"})
	db.Create(&Todo{Title: "Todo 2"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/todos/all", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Equal(t, "All todos deleted successfully", response["message"])
	assert.Equal(t, float64(2), response["deleted_count"])

	// Verify no todos remain
	var remainingTodos []Todo
	db.Find(&remainingTodos)
	assert.Len(t, remainingTodos, 0)
}
