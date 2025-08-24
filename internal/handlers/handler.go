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

// CreateNote
// @Summary      Create a note
// @Description  creating a note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param   note body models.CreateNoteRequest true "Данные заметки"
// @Success 201 {object} map[string]interface{} "Заметка создана"
// @Failure 400 {object} map[string]interface{} "Неверные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /notes [post]
func (h *NoteHandler) CreateNote(c *gin.Context) {
	var noteRequest models.CreateNoteRequest
	if err := c.ShouldBindJSON(&noteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(noteRequest.Title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	if len(noteRequest.Content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content cannot be empty"})
		return
	}
	note := models.Note{
		Title:   noteRequest.Title,
		Content: noteRequest.Content,
	}
	result := h.DB.Create(&note)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNote
// @Summary Show note
// @Description get note by ID
// @Tags notes
// @Accept  json
// @Produce  json
// @Param   id path int true "ID заметки"
// @Success 200 {object} models.Note "Заметка найдена"
// @Failure 400 {object} map[string]interface{} "Неверный ID"
// @Failure 404 {object} map[string]interface{} "Заметка не найдена"
// @Router /notes/{id} [get]
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

// GetAllNotes
// @Summary List notes
// @Description get notes
// @Tags notes
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Список заметок"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /notes [get]
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

// UpdateNote
// @Summary Update note
// @Description update note by ID
// @Tags notes
// @Accept  json
// @Produce  json
// @Param   id path int true "ID заметки"
// @Param   updateData body models.UpdateNoteRequest true "Данные для обновления"
// @Success 200 {object} map[string]interface{} "Заметка обновлена"
// @Failure 400 {object} map[string]interface{} "Неверные данные"
// @Failure 404 {object} map[string]interface{} "Заметка не найдена"
// @Router /notes/{id} [put]
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

	var updateData models.UpdateNoteRequest
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

// DeleteNote
// @Summary Delete note
// @Description delete note by ID
// @Tags notes
// @Accept  json
// @Produce  json
// @Param   id path int true "ID заметки"
// @Success 200 {object} map[string]interface{} "Заметка удалена"
// @Failure 400 {object} map[string]interface{} "Неверный ID"
// @Failure 404 {object} map[string]interface{} "Заметка не найдена"
// @Router /notes/{id} [delete]
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
