package models

import "time"

type Note struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title" binding:"required"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type UpdateNoteInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
