package repository

import (
	"github.com/sarita-growexx/note_with_alarm/models"
)

type NoteRepository interface {
	Create(note *models.Note) error
	Update(note *models.Note) error
	Delete(id uint) error
	GetById(id uint) (*models.Note, error)
	GetAll() ([]*models.Note, error)
	Search(query string) ([]*models.Note, error)
	GetNoteByTitle(title string) (*models.Note, error)
}
