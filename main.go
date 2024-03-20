package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sarita-growexx/note_with_alarm/config"
	"github.com/sarita-growexx/note_with_alarm/controllers"
	"github.com/sarita-growexx/note_with_alarm/helper"
	"github.com/sarita-growexx/note_with_alarm/models"
	"github.com/sarita-growexx/note_with_alarm/repository"
	routers "github.com/sarita-growexx/note_with_alarm/router"
	"github.com/sarita-growexx/note_with_alarm/services"
	"github.com/sarita-growexx/note_with_alarm/utils"
)

func main() {
	time.LoadLocation("Asia/Kolkata") // "Asia/Kolkata" is the IANA Time Zone identifier for IST

	fmt.Println("Current Time main.go:", time.Now())

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	//Database
	db := config.ConnectionDB(&loadConfig)

	db.AutoMigrate(&models.Note{})

	noteRepository := repository.NewNoteRepository(db)
	noteService := services.NewNoteService(noteRepository)
	noteController := controllers.NewNoteController(noteService)

	routes := routers.SetupRouter(noteController)

	// Start the background task to check for upcoming deadlines
	go func() {
		for {
			notes, err := noteService.GetAllNotes()
			if err != nil {
				log.Println("Failed to retrieve notes for checking deadlines: ", err)
			}

			utils.SetAlarmForNotes(notes)

			// Sleep for a specified duration before checking again
			time.Sleep(time.Minute * 5)
		}
	}()

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverErr := server.ListenAndServe()
	helper.ErrorPanic(serverErr)
}
