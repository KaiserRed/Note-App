package handlers

import (
	"net/http"
	"note-app/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoteHandler struct {
	DB *gorm.DB
}

func NewNoteHandler(db *gorm.DB) *NoteHandler {
	return &NoteHandler{
		DB: db,
	}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(note.Title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	if len(note.Content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content cannot be empty"})
		return
	}

	result := h.DB.Create(&note)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      note.ID,
		"message": "Note created successfully",
	})
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var note models.Note
	result := h.DB.First(&note, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	var notes []models.Note
	result := h.DB.Find(&notes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notes: " + result.Error.Error()})
		return
	}

	if len(notes) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"notes":   []string{},
			"message": "No notes found",
			"count":   0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
		"count": len(notes),
	})
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var existingNote models.Note
	if result := h.DB.First(&existingNote, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var updateData models.UpdateNoteInput
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if updateData.Title == nil && updateData.Content == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if updateData.Title != nil {
		if len(*updateData.Title) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
			return
		}
	}

	result := h.DB.Model(&existingNote).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var Note models.Note
	if result := h.DB.First(&Note, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	result := h.DB.Delete(&models.Note{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
