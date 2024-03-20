package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sarita-growexx/note_with_alarm/controllers"
)

func SetupRouter(noteController *controllers.NoteController) *gin.Engine {
    router := gin.Default()

    // Initialize routes
    api := router.Group("/api")
    {
        notes := api.Group("/notes")
        {
            notes.POST("/", noteController.CreateNoteHandler)
            notes.PUT("/:id", noteController.UpdateNoteHandler)
            notes.DELETE("/:id", noteController.DeleteNoteHandler)
            notes.GET("/:id", noteController.GetNoteByIDHandler)
            notes.GET("/", noteController.GetAllNotesHandler)
            notes.GET("/search", noteController.SearchNotesHandler)
        }
    }

    return router
}
