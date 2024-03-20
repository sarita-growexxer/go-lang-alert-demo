package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/sarita-growexx/note_with_alarm/repository"
	"github.com/sarita-growexx/note_with_alarm/utils"
	"gorm.io/gorm"
)

const standardTime = "2006-01-02T15:04:05Z"

var ErrNoteNotFound = errors.New("note not found")

type NoteServiceImpl struct {
	noteRepository repository.NoteRepository
}

func NewNoteService(noteRepository repository.NoteRepository) *NoteServiceImpl {
	return &NoteServiceImpl{
		noteRepository: noteRepository,
	}
}

func (s *NoteServiceImpl) CreateNote(note *models.Note) error {

	existingNote, err := s.noteRepository.GetNoteByTitle(note.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to check duplicate title")
	}

	if existingNote != nil {
		fmt.Println("Existing note=", existingNote)
		return errors.New("duplicate title, please choose a different title")

	}

	deadlineLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return err
	}
	deadline, err := time.ParseInLocation(standardTime, note.Deadline.Format(standardTime), deadlineLocation)
	if err != nil {
		return err
	}

	note.Deadline = deadline

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	err = s.noteRepository.Create(note)
	if err != nil {
		fmt.Println("Error creating note:", err)

		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("duplicate title, please choose a different title")
		}

		return fmt.Errorf("failed to create note: %w", err)
	}

	// Set the alarm for the note
	utils.SetAlarmForNotes([]*models.Note{note})

	return nil
}

func (s *NoteServiceImpl) UpdateNote(note *models.Note) error {
	existingNote, err := s.noteRepository.GetById(note.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing note: %w", err)
	}

	if existingNote == nil {
		return errors.New("note not found")
	}

	// existingNote, err := s.noteRepository.GetNoteByTitle(note.Title)
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return errors.New("failed to check duplicate title")
	// }

	if existingNote.ID != 0 && existingNote.ID != note.ID {
		return errors.New("duplicate title, please choose a different title")
	}

	// Parse deadline
	deadlineLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return err
	}
	deadline, err := time.ParseInLocation(standardTime, note.Deadline.Format(standardTime), deadlineLocation)
	if err != nil {
		return err
	}

	note.Deadline = deadline

	note.UpdatedAt = time.Now()

	err = s.noteRepository.Update(note)
	if err != nil {
		return errors.New("failed to update note")
	}

	utils.SetAlarmForNotes([]*models.Note{note})

	return nil
}

func (s *NoteServiceImpl) DeleteNote(id uint) error {
	err := s.noteRepository.Delete(id)
	if err != nil {
		return errors.New("failed to delete note")
	}

	return nil
}

func (s *NoteServiceImpl) GetAllNotes() ([]*models.Note, error) {
	notes, err := s.noteRepository.GetAll()
	if err != nil {
		return nil, errors.New("failed to get all notes")
	}

	return notes, nil
}

func (s *NoteServiceImpl) GetNoteById(id uint) (*models.Note, error) {
	note, err := s.noteRepository.GetById(id)
	if err != nil {
		return nil, errors.New("failed to get note by ID")
	}

	return note, nil
}

func (s *NoteServiceImpl) SearchNotes(query string) ([]*models.Note, error) {
	return s.noteRepository.Search(query)
}
