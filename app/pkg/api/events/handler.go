
package events

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	models "app/pkg/api/events/models"
	"app/pkg/db"
	_ "github.com/lib/pq"
)

// HandleEvents is a placeholder for the /api/events endpoint handler.
func HandleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetEventsHandler(w, r)
	case http.MethodPost:
		CreateEventHandler(w, r)
	case http.MethodPut:
		UpdateEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
		
	} 
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"events":[]}`))
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var NewEvent models.Event
	err := json.NewDecoder(r.Body).Decode(&NewEvent)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	eventRepository := models.NewEventRepository(db.Get())
	eventRepository.CreateEvent(&NewEvent)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Event created successfully"}`))
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"event":{"id":1,"title":"Sample Event"}}`))
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	d := db.Get()
	if d == nil {
		http.Error(w, "database not available", http.StatusInternalServerError)
		return
	}
	repo := models.NewEventRepository(d)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_ = repo.Update(ctx, &NewEvent)
	w.Write([]byte(`{"message":"Event updated successfully"}`))
}
