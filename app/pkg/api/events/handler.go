package events

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"github.com/google/uuid"

	models "github.com/johnkespitia/taller-go-repo/app/pkg/api/events/models"
	db "github.com/johnkespitia/taller-go-repo/app/pkg/db"
	_ "github.com/lib/pq"
)

// HandleEvents routes /api/events methods
func HandleEvents(w http.ResponseWriter, r *http.Request) {
	d := db.Get()
	if d == nil {
		http.Error(w, "database not available", http.StatusInternalServerError)
		return
	}
	repo := models.NewEventRepository(d)
	switch r.Method {
	case http.MethodGet:
		GetEventsHandler(w, r, repo)
	case http.MethodPost:
		CreateEventHandler(w, r, repo)
	case http.MethodPut:
		UpdateEventHandler(w, r, repo)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request, repo *models.EventRepository) {
	var newEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	newEvent.ID = uuid.New().String()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	fmt.Println("NewEvent:", newEvent)
	if err := repo.Create(ctx, &newEvent); err != nil {
		fmt.Println("Error creating event:", err)
		http.Error(w, "failed to create event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event created successfully"})
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request, repo *models.EventRepository) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id := r.URL.Path[len("/api/events/"):]
	if(id != ""){
		event, err := repo.GetByID(ctx, id)
		if err != nil {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(event)
	}else{
		events, err := repo.GetAll(ctx)
		if err != nil {
			http.Error(w, "failed to get events", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(events)
	}
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request, repo *models.EventRepository) {
	var ev models.Event
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if ev.ID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := repo.Update(ctx, &ev); err != nil {
		http.Error(w, "failed to update event", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Event updated successfully"})
}
