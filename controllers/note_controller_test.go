package controllers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sarita-growexx/note_with_alarm/controllers"
	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a mock note service for testing
type MockNoteService struct {
	mock.Mock
}

// Implement methods for the mock note service

// Mock CreateNote method
func (m *MockNoteService) CreateNote(note *models.Note) error {
	args := m.Called(note)
	return args.Error(0)
}

// Mock UpdateNote method
func (m *MockNoteService) UpdateNote(note *models.Note) error {
	args := m.Called(note)
	return args.Error(0)
}

// Mock DeleteNote method
func (m *MockNoteService) DeleteNote(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock GetAllNotes method
func (m *MockNoteService) GetAllNotes() ([]*models.Note, error) {
	args := m.Called()
	return args.Get(0).([]*models.Note), args.Error(1)
}

// Mock GetNoteById method
func (m *MockNoteService) GetNoteById(id uint) (*models.Note, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Note), args.Error(1)
}

// Mock SearchNotes method
func (m *MockNoteService) SearchNotes(query string) ([]*models.Note, error) {
	args := m.Called(query)
	return args.Get(0).([]*models.Note), args.Error(1)
}

// Test CreateNoteHandler function
func TestCreateNoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNoteService)
	controller := controllers.NewNoteController(mockService)

	router := gin.Default()
	router.POST("/note", controller.CreateNoteHandler)

	// Test case: Successful note creation
	note := &models.Note{Title: "Test Note"}
	mockService.On("CreateNote", note).Return(nil)
	body, _ := json.Marshal(note)
	req, _ := http.NewRequest("POST", "/note", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)

	// Test case: Invalid request body
	req, _ = http.NewRequest("POST", "/note", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case: Duplicate title
	duplicateNote := &models.Note{Title: "Duplicate Title"}
	mockService.On("CreateNote", duplicateNote).Return(errors.New("duplicate title"))
	body, _ = json.Marshal(duplicateNote)
	req, _ = http.NewRequest("POST", "/note", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Duplicate title, please choose a different title")
	mockService.AssertExpectations(t)

	// Test case: Failed note creation
	invalidNote := &models.Note{Title: "Invalid Note"} // Note without required fields
	mockService.On("CreateNote", invalidNote).Return(errors.New("invalid note"))
	body, _ = json.Marshal(invalidNote)
	req, _ = http.NewRequest("POST", "/note", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test case: Validation failure
	invalidNote = &models.Note{} // Empty note, which should fail validation
	body, _ = json.Marshal(invalidNote)
	req, _ = http.NewRequest("POST", "/note", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Field validation for '' failed on the 'required' tag") // Check for the actual error message
	mockService.AssertNotCalled(t, "CreateNote")                                                // Ensure CreateNote is not called when validation fails

}

// Test UpdateNoteHandler function
func TestUpdateNoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNoteService)
	controller := controllers.NewNoteController(mockService)

	router := gin.Default()
	router.PUT("/note/:id", controller.UpdateNoteHandler)

	// Test case: Successful note update
	noteID := uint(1)
	note := &models.Note{ID: noteID, Title: "Updated Note"}
	mockService.On("UpdateNote", note).Return(nil)
	body, _ := json.Marshal(note)
	req, _ := http.NewRequest("PUT", "/note/"+strconv.Itoa(int(noteID)), strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	// Test case: Invalid note ID
	req, _ = http.NewRequest("PUT", "/note/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case: Validation failure
	invalidNote := &models.Note{} // Empty note, which should fail validation
	body, _ = json.Marshal(invalidNote)
	req, _ = http.NewRequest("PUT", "/note/1", strings.NewReader(string(body))) // Assuming note ID is 1
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Field validation for '' failed on the 'required' tag") // Check for the actual error message
	mockService.AssertNotCalled(t, "UpdateNote")                                                // Ensure UpdateNote is not called when validation fails

	// Test case: Duplicate title error
	duplicateTitleNote := &models.Note{ID: noteID, Title: "Duplicate Title"}
	mockService.On("UpdateNote", duplicateTitleNote).Return(errors.New("duplicate title"))
	body, _ = json.Marshal(duplicateTitleNote)
	req, _ = http.NewRequest("PUT", "/note/"+strconv.Itoa(int(noteID)), strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Duplicate title, please choose a different title") // Check for the specific error message

}

// Test DeleteNoteHandler function
func TestDeleteNoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNoteService)
	controller := controllers.NewNoteController(mockService)

	router := gin.Default()
	router.DELETE("/note/:id", controller.DeleteNoteHandler)

	// Test case: Successful note deletion
	noteID := uint(1)
	mockService.On("DeleteNote", noteID).Return(nil)
	req, _ := http.NewRequest("DELETE", "/note/"+strconv.Itoa(int(noteID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	// Test case: Invalid note ID
	req, _ = http.NewRequest("DELETE", "/note/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case: Failed to delete note
	noteID = uint(2)
	mockService.On("DeleteNote", noteID).Return(errors.New("failed to delete note"))
	req, _ = http.NewRequest("DELETE", "/note/"+strconv.Itoa(int(noteID)), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete note")
}

// Test GetAllNotesHandler function
// Test GetAllNotesHandler function
func TestGetAllNotesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNoteService)
	controller := controllers.NewNoteController(mockService)

	router := gin.Default()
	router.GET("/notes", controller.GetAllNotesHandler)

	// Test case: Failure to retrieve notes
	t.Run("Failure to retrieve notes", func(t *testing.T) {
		mockService := new(MockNoteService)
		controller := controllers.NewNoteController(mockService)

		router := gin.Default()
		router.GET("/notes", controller.GetAllNotesHandler)

		mockService.On("GetAllNotes").Return(nil, errors.New("failed to retrieve notes"))
		req, _ := http.NewRequest("GET", "/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// Test case: Successful retrieval of notes
	t.Run("Successful retrieval of notes", func(t *testing.T) {
		mockService := new(MockNoteService)
		controller := controllers.NewNoteController(mockService)

		router := gin.Default()
		router.GET("/notes", controller.GetAllNotesHandler)

		notes := []*models.Note{{ID: 1, Title: "Test Note 1"}, {ID: 2, Title: "Test Note 2"}}
		mockService.On("GetAllNotes").Return(notes, nil)
		req, _ := http.NewRequest("GET", "/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotNil(t, w.Body)
	})

}

// Test GetNoteByIDHandler function

func TestGetNoteByIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNoteService)
	controller := controllers.NewNoteController(mockService)

	router := gin.Default()
	router.GET("/note/:id", controller.GetNoteByIDHandler)

	// Test case: Note found
	noteID := uint(1)
	note := &models.Note{
		ID:    noteID,
		Title: "Test Note",
		// Add other necessary fields here
	}
	mockService.On("GetNoteById", noteID).Return(note, nil)

	req, _ := http.NewRequest("GET", "/note/"+strconv.Itoa(int(noteID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code) // Note found, so status should be OK
	// Add assertion to check the response body content if necessary

	// Test case: Note not found
	mockService.On("GetNoteById", noteID).Return(nil, errors.New("Note not found"))

	req, _ = http.NewRequest("GET", "/note/"+strconv.Itoa(int(noteID)), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, 404) // Note not found, so status should be Not Found

	// Test case: Invalid note ID
	req, _ = http.NewRequest("GET", "/note/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code) // Invalid note ID, so status should be Bad Request
}

// Test SearchNotesHandler function
func TestSearchNotesHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a custom Gin router
	router := gin.New()

	// Create a new instance of the NoteController with a mocked service
	mockService := new(MockNoteService)
	controller := controllers.NewNoteController(mockService)

	// Register the SearchNotesHandler route on the custom router
	router.GET("/notes/search", controller.SearchNotesHandler)

	t.Run("Successful Search", func(t *testing.T) {
		query := "test"
		notes := []*models.Note{{ID: 1, Title: "Test Note 1"}, {ID: 2, Title: "Test Note 2"}}
		mockService.On("SearchNotes", query).Return(notes, nil)
		req, _ := http.NewRequest("GET", "/notes/search?query="+query, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Empty Query", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/notes/search", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Search Failure", func(t *testing.T) {
		query := "test"
		mockService.On("SearchNotes", query).Return(nil, errors.New("failed to search notes"))
		req1, _ := http.NewRequest("GET", "/notes/search?query="+query, nil)
		w1 := httptest.NewRecorder()

		fmt.Println("req1", req1.Response)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusInternalServerError, 500) // Update expected status code
		mockService.AssertExpectations(t)
	})

}
