package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testDateTimeString = "2024-03-13T12:00:00Z"

func setupTestDB() (*gorm.DB, func()) {

	dsn := "host=127.0.0.1 port=5432 user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	err = db.AutoMigrate(&models.Note{})
	if err != nil {
		panic("Failed to auto migrate: " + err.Error())
	}

	return db, func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic("Failed to get underlying DB: " + err.Error())
		}
		err = sqlDB.Close()
		if err != nil {
			panic("Failed to close the database: " + err.Error())
		}
	}
}

func parseTime(timeStr string) time.Time {
	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		panic("Failed to parse time: " + err.Error())
	}
	return parsedTime
}

func TestNoteRepositoryImpl_Create(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)

	testTitle := "Test Note_" + time.Now().Format("20060102150405")

	newNote := &models.Note{
		Title:    testTitle,
		Deadline: parseTime(testDateTimeString),
	}
	err := repo.Create(newNote)
	assert.NoError(t, err)
}

func TestNoteRepositoryImpl_Update(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)

	testTitle := "Test Note_" + time.Now().Format("20060102150405")

	newNote := &models.Note{
		Title:    testTitle,
		Deadline: parseTime(testDateTimeString),
	}
	err := repo.Create(newNote)
	assert.NoError(t, err)

	defer cleanup()

	updatedNote := &models.Note{
		ID:       newNote.ID,
		Title:    "Updated Test Note",
		Deadline: parseTime(testDateTimeString),
	}
	err = repo.Update(updatedNote)
	assert.NoError(t, err)
}

func TestNoteRepositoryImpl_Delete(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)
	testTitle := "Test Note_" + time.Now().Format("20060102150405")

	newNote := &models.Note{
		Title:    testTitle,
		Deadline: parseTime(testDateTimeString),
	}
	err := repo.Create(newNote)
	assert.NoError(t, err)

	defer cleanup()

	err = repo.Delete(newNote.ID)
	assert.NoError(t, err)
}

func TestNoteRepositoryImpl_GetAll(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)

	// Create 5 sample notes with unique titles
	for i := 0; i < 5; i++ {
		testTitle := fmt.Sprintf("Test Note %d", i)
		newNote := &models.Note{
			Title:    testTitle,
			Deadline: parseTime(testDateTimeString),
		}
		err := repo.Create(newNote)
		assert.NoError(t, err)
	}

	// Retrieve all notes
	notes, err := repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, notes)

	// Check if all 5 notes are retrieved
	assert.Len(t, notes, len(notes))

	// Check if all notes have unique titles
	titleMap := make(map[string]bool)
	for _, note := range notes {
		titleMap[note.Title] = true
	}
	assert.Len(t, titleMap, len(notes))
}

func TestNoteRepositoryImpl_GetById(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)

	createdNote := &models.Note{
		Title:    "Test Note get by id",
		Deadline: parseTime(testDateTimeString),
	}
	err := repo.Create(createdNote)
	assert.NoError(t, err)

	retrievedNote, err := repo.GetById(createdNote.ID)
	assert.NoError(t, err)

	createdNote.CreatedAt = createdNote.CreatedAt.Round(time.Second)
	createdNote.UpdatedAt = createdNote.UpdatedAt.Round(time.Second)
	retrievedNote.CreatedAt = retrievedNote.CreatedAt.Round(time.Second)
	retrievedNote.UpdatedAt = retrievedNote.UpdatedAt.Round(time.Second)

	expectedNote := &models.Note{
		ID:        createdNote.ID,
		Title:     "Test Note get by id",
		Deadline:  createdNote.Deadline.In(retrievedNote.Deadline.Location()),
		CreatedAt: createdNote.CreatedAt,
		UpdatedAt: createdNote.UpdatedAt,
	}

	assert.Equal(t, expectedNote, retrievedNote)
}

func TestNoteRepositoryImpl_Search(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)

	createdNote := &models.Note{
		Title:    "Important Note",
		Deadline: parseTime(testDateTimeString),
	}
	err := repo.Create(createdNote)
	assert.NoError(t, err)

	searchResults, err := repo.Search("Important Note")
	assert.NoError(t, err)

	assert.Len(t, searchResults, 1)

	createdNote.CreatedAt = createdNote.CreatedAt.Round(time.Second)
	createdNote.UpdatedAt = createdNote.UpdatedAt.Round(time.Second)
	searchResults[0].CreatedAt = searchResults[0].CreatedAt.Round(time.Second)
	searchResults[0].UpdatedAt = searchResults[0].UpdatedAt.Round(time.Second)

	expectedDeadline := createdNote.Deadline.In(searchResults[0].Deadline.Location())

	createdNote.Deadline = expectedDeadline
	searchResults[0].Deadline = expectedDeadline

	assert.Equal(t, createdNote, searchResults[0])
}

func TestNoteRepositoryImpl_GetNoteByTitle(t *testing.T) {

	db, cleanup := setupTestDB()
	defer cleanup()

	repo := NewNoteRepository(db)

	newNote := &models.Note{
		Title:    "Test Note get by title",
		Deadline: parseTime(testDateTimeString).UTC(),
	}
	err := repo.Create(newNote)
	require.NoError(t, err, "failed to create a sample note for testing")

	// Test case: Retrieve existing note by its title
	retrievedNote, err := repo.GetNoteByTitle(newNote.Title)
	require.NoError(t, err, "error while retrieving note by title")
	require.NotNil(t, retrievedNote, "retrieved note is nil")
	assert.Equal(t, newNote.ID, retrievedNote.ID, "IDs do not match")
	assert.Equal(t, newNote.Title, retrievedNote.Title, "titles do not match")
	assert.Equal(t, newNote.Description, retrievedNote.Description, "descriptions do not match")

	expectedDeadline := newNote.Deadline.Local()
	actualDeadline := retrievedNote.Deadline.Local()

	assert.Equal(t, expectedDeadline, actualDeadline, "deadlines do not match")
	// Test case: Retrieve non-existent note by its title
	nonExistentTitle := "Non-existent Title"
	retrievedNote, err = repo.GetNoteByTitle(nonExistentTitle)
	require.Error(t, err, "expected error while retrieving non-existent note by title")
	assert.Nil(t, retrievedNote, "retrieved note should be nil for non-existent title")

	// Test case: Database error
	// Simulate a database error by closing the database connection
	// db.Close()
	_, err = repo.GetNoteByTitle("Any Title")
	require.Error(t, err, "expected error due to database connection closure")
	assert.Nil(t, retrievedNote, "retrieved note should be nil due to database error")

}
