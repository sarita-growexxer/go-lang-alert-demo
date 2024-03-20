package utils

import (
	"testing"
	"time"

	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/stretchr/testify/assert"
)

func TestSetAlarmForNotes(t *testing.T) {
	// Create a sample note with a deadline 30 minutes from now
	note := &models.Note{
		Title:    "Test Note",
		Deadline: time.Now().Add(30 * time.Minute),
	}

	// Mock the alertTriggered map to capture whether an alert is triggered
	alertTriggered = make(map[uint]bool)

	// Call SetAlarmForNotes with the sample note
	SetAlarmForNotes([]*models.Note{note})

	// Wait for a short duration to allow the goroutine to execute
	time.Sleep(1 * time.Second)

	// Assertions
	assert.True(t, alertTriggered[note.ID], "Alert should be triggered for the sample note")
}

// Add more test cases as needed
