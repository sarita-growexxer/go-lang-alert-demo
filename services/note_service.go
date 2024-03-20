package services

import (
	"github.com/sarita-growexx/note_with_alarm/models"
)

type NoteService interface {
	CreateNote(note *models.Note) error
	UpdateNote(note *models.Note) error
	DeleteNote(id uint) error
	GetNoteById(id uint) (*models.Note, error)
	GetAllNotes() ([]*models.Note, error)
	SearchNotes(query string) ([]*models.Note, error)
}
