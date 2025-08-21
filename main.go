package main

import (
	"log"
	"net/http"
	"note-app/internal/database"
	"note-app/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	db := database.InitDB()

	r := gin.Default()

	NoteHandler := handlers.NewNoteHandler(db)
	r.GET("/health", healthCheck)
	r.POST("/notes", NoteHandler.CreateNote)
	r.GET("/notes", NoteHandler.GetAllNotes)
	r.GET("/notes/:id", NoteHandler.GetNote)
	r.PUT("/notes/:id", NoteHandler.UpdateNote)
	r.DELETE("/notes/:id", NoteHandler.DeleteNote)

	//TODO SWAGGER

	//TODO go func
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Health check
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Note service is running",
	})
}
