package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Flashcard struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	db, err = sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASS")+"@tcp(localhost:3306)/"+os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()
	r.Use(cors.Default())

	// Routes
	r.GET("/flashcards", getFlashcards)
	r.POST("/flashcards", createFlashcard)
	r.DELETE("/flashcards/:id", deleteFlashcard)

	// Start server
	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

// Get all flashcards
func getFlashcards(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, content FROM flashcards")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flashcards"})
		return
	}
	defer rows.Close()

	var flashcards []Flashcard
	for rows.Next() {
		var fc Flashcard
		if err := rows.Scan(&fc.ID, &fc.Title, &fc.Content); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		flashcards = append(flashcards, fc)
	}

	c.JSON(http.StatusOK, flashcards)
}

// Create a new flashcard
func createFlashcard(c *gin.Context) {
	var fc Flashcard
	if err := c.ShouldBindJSON(&fc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	_, err := db.Exec("INSERT INTO flashcards (title, content) VALUES (?, ?)", fc.Title, fc.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flashcard"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Flashcard created successfully"})
}

// Delete a flashcard
func deleteFlashcard(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM flashcards WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete flashcard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flashcard deleted successfully"})
}
