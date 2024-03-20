package repository

import (
	"errors"
	"fmt"

	"github.com/sarita-growexx/note_with_alarm/models"
	"gorm.io/gorm"
)

var ErrNoteNotFound = errors.New("note not found")

type NoteRepositoryImpl struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &NoteRepositoryImpl{db: db}
}

// Create implements NoteRepository.
func (n *NoteRepositoryImpl) Create(note *models.Note) error {
	return n.db.Create(note).Error
}

// Delete implements NoteRepository.
func (n *NoteRepositoryImpl) Delete(id uint) error {
	return n.db.Delete(&models.Note{}, id).Error
}

// GetAll implements NoteRepository.
func (r *NoteRepositoryImpl) GetAll() ([]*models.Note, error) {
	var notes []*models.Note
	if err := r.db.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// GetById implements NoteRepository.
func (n *NoteRepositoryImpl) GetById(id uint) (*models.Note, error) {
	var note models.Note
	err := n.db.First(&note, id).Error
	return &note, err
}

// Search implements NoteRepository.
func (n *NoteRepositoryImpl) Search(query string) ([]*models.Note, error) {
	var notes []*models.Note
	err := n.db.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&notes).Error
	fmt.Println("Notes :: ", notes)
	return notes, err
}

// Update implements NoteRepository.
func (n *NoteRepositoryImpl) Update(note *models.Note) error {
	return n.db.Save(note).Error
}

func (r *NoteRepositoryImpl) GetNoteByTitle(title string) (*models.Note, error) {
	var note models.Note
	err := r.db.Where("title = ?", title).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}
