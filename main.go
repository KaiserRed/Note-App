package main

import (
	"log"
	"net/http"
	_ "note-app/docs"
	"note-app/internal/database"
	"note-app/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	db := database.InitDB()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	{
		NoteHandler := handlers.NewNoteHandler(db)
		v1.GET("/health", healthCheck)
		v1.POST("/notes", NoteHandler.CreateNote)
		v1.GET("/notes", NoteHandler.GetAllNotes)
		v1.GET("/notes/:id", NoteHandler.GetNote)
		v1.PUT("/notes/:id", NoteHandler.UpdateNote)
		v1.DELETE("/notes/:id", NoteHandler.DeleteNote)
	}

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

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Note service is running",
	})
}
