package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/sarita-growexx/note_with_alarm/services"
)

type NoteController struct {
	noteService services.NoteService
}

const invalidIDErr = "Invalid note ID"

var validate *validator.Validate

func NewNoteController(noteService services.NoteService) *NoteController {
	return &NoteController{
		noteService: noteService,
	}
}

func (c *NoteController) CreateNoteHandler(ctx *gin.Context) {
	var note models.Note

	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateNoteFields(&note); err != nil {
		fmt.Println("Error:: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.noteService.CreateNote(&note)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate title") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate title, please choose a different title"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Note created successfully", "note": note})
}

func (c *NoteController) UpdateNoteHandler(ctx *gin.Context) {
	noteID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDErr})
		return
	}

	var updatedNote models.Note
	if err := ctx.ShouldBindJSON(&updatedNote); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate note fields
	if err := validateNoteFields(&updatedNote); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the ID of the note
	updatedNote.ID = uint(noteID)

	err = c.noteService.UpdateNote(&updatedNote)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
	// 	return
	// }

	if err != nil {
		// Check for specific error messages
		if strings.Contains(err.Error(), "duplicate title") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate title, please choose a different title"})
			return
		}

		// Handle other error cases
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Note updated successfully", "note": updatedNote})
}

func (c *NoteController) DeleteNoteHandler(ctx *gin.Context) {
	noteID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDErr})
		return
	}

	err = c.noteService.DeleteNote(uint(noteID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (c *NoteController) GetAllNotesHandler(ctx *gin.Context) {
	notes, err := c.noteService.GetAllNotes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notes"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (c *NoteController) GetNoteByIDHandler(ctx *gin.Context) {
	noteID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": invalidIDErr})
		return
	}

	note, err := c.noteService.GetNoteById(uint(noteID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"note": note})
}

func (nc *NoteController) SearchNotesHandler(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'query' is required"})
		return
	}

	notes, err := nc.noteService.SearchNotes(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search notes"})
		return
	}

	ctx.JSON(http.StatusOK, notes)
}

func validateNoteFields(note *models.Note) error {
	validate = validator.New()

	// Validate title
	if err := validate.Var(note.Title, "required,min=3,max=50"); err != nil {
		return err
	}

	// // Validate deadline format
	// if err := validate.Var(note.Deadline, "required,rfc3339"); err != nil {
	//     return err
	// }

	return nil
}
