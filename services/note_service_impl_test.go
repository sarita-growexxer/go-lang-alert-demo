package services

import (
	"testing"
	"time"

	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockNoteRepository struct {
	mock.Mock
}

func (m *mockNoteRepository) Create(note *models.Note) error {
	args := m.Called(note)
	return args.Error(0)
}

func (m *mockNoteRepository) Update(note *models.Note) error {
	args := m.Called(note)
	return args.Error(0)
}

func (m *mockNoteRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockNoteRepository) GetAll() ([]*models.Note, error) {
	args := m.Called()
	return args.Get(0).([]*models.Note), args.Error(1)
}

func (m *mockNoteRepository) GetById(id uint) (*models.Note, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *mockNoteRepository) Search(query string) ([]*models.Note, error) {
	args := m.Called(query)
	return args.Get(0).([]*models.Note), args.Error(1)
}

func (m *mockNoteRepository) GetNoteByTitle(title string) (*models.Note, error) {
	args := m.Called(title)
	return args.Get(0).(*models.Note), args.Error(1)
}

func TestNoteServiceImpl_CreateNote(t *testing.T) {

	mockRepo := new(mockNoteRepository)
	service := NewNoteService(mockRepo)

	newNote := &models.Note{
		Title:    "Test Note",
		Deadline: time.Now().Add(time.Hour),
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Note")).Return(nil).Once()

	err := service.CreateNote(newNote)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestNoteServiceImpl_UpdateNote(t *testing.T) {

	mockRepo := new(mockNoteRepository)
	service := NewNoteService(mockRepo)

	updatedNote := &models.Note{
		ID:       1,
		Title:    "Updated Test Note",
		Deadline: time.Now().Add(time.Hour),
	}

	mockRepo.On("GetById", uint(1)).Return(&models.Note{}, nil).Once()

	mockRepo.On("Update", updatedNote).Return(nil).Once()

	err := service.UpdateNote(updatedNote)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestNoteServiceImpl_DeleteNote(t *testing.T) {
	mockRepo := new(mockNoteRepository)
	service := NewNoteService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteNote(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNoteServiceImpl_GetAllNotes(t *testing.T) {
	mockRepo := new(mockNoteRepository)
	service := NewNoteService(mockRepo)

	notes := []*models.Note{
		{ID: 1, Title: "Note 1", Deadline: time.Now().Add(1 * time.Hour)},
		{ID: 2, Title: "Note 2", Deadline: time.Now().Add(2 * time.Hour)},
	}

	mockRepo.On("GetAll").Return(notes, nil)

	resultNotes, err := service.GetAllNotes()

	assert.NoError(t, err)
	assert.Equal(t, notes, resultNotes)
	mockRepo.AssertExpectations(t)
}

func TestNoteServiceImpl_GetNoteById(t *testing.T) {
	mockRepo := new(mockNoteRepository)
	service := NewNoteService(mockRepo)

	testNote := &models.Note{ID: 1, Title: "Test Note", Deadline: time.Now().Add(1 * time.Hour)}

	mockRepo.On("GetById", uint(1)).Return(testNote, nil)

	resultNote, err := service.GetNoteById(1)

	assert.NoError(t, err)
	assert.Equal(t, testNote, resultNote)
	mockRepo.AssertExpectations(t)
}

func TestNoteServiceImpl_SearchNotes(t *testing.T) {
	mockRepo := new(mockNoteRepository)
	service := NewNoteService(mockRepo)

	notes := []*models.Note{
		{ID: 1, Title: "Important Note", Deadline: time.Now().Add(1 * time.Hour)},
		{ID: 2, Title: "Random Note", Deadline: time.Now().Add(2 * time.Hour)},
	}

	mockRepo.On("Search", "Important").Return(notes[:1], nil)

	resultNotes, err := service.SearchNotes("Important")

	assert.NoError(t, err)
	assert.Equal(t, notes[:1], resultNotes)
	mockRepo.AssertExpectations(t)
}
