package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"taskmaster-api/internal/middleware"
	"taskmaster-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate schema
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully.")
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	initDB()

	r := gin.Default()

	// --- Public Endpoints ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "TaskMaster API is healthy",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		// Dummy credentials
		var body struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		if body.Username != "admin" || body.Password != "password" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
			return
		}

		secret := os.Getenv("JWT_SECRET")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": body.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login successful",
			"data":    gin.H{"token": tokenString},
		})
	})

	// --- Protected Endpoints (require JWT) ---
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// GET /api/v1/tasks - Ambil semua task
		api.GET("/tasks", func(c *gin.Context) {
			var tasks []models.Task
			db.Find(&tasks)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Tasks retrieved successfully",
				"data":    tasks,
			})
		})

		// POST /api/v1/tasks - Buat task baru
		api.POST("/tasks", func(c *gin.Context) {
			var req models.TaskRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
				return
			}

			task := models.Task{Title: req.Title, Status: "pending"}
			if req.Status != "" {
				task.Status = req.Status
			}

			db.Create(&task)
			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"message": "Task created successfully",
				"data":    task,
			})
		})

		// PUT /api/v1/tasks/:id - Update task
		api.PUT("/tasks/:id", func(c *gin.Context) {
			id := c.Param("id")
			var task models.Task
			if result := db.First(&task, id); result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Task not found"})
				return
			}

			var req models.TaskRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
				return
			}

			task.Title = req.Title
			if req.Status != "" {
				task.Status = req.Status
			}
			db.Save(&task)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Task updated successfully",
				"data":    task,
			})
		})

		// DELETE /api/v1/tasks/:id - Soft delete task
		api.DELETE("/tasks/:id", func(c *gin.Context) {
			id := c.Param("id")
			var task models.Task
			if result := db.First(&task, id); result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Task not found"})
				return
			}

			db.Delete(&task)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Task deleted successfully",
			})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
