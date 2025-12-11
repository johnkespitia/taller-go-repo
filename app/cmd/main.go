package main

import (
	"fmt"
	"log"
	"net/http"
		events "github.com/johnkespitia/taller-go-repo/app/pkg/api/events"
		"github.com/johnkespitia/taller-go-repo/app/pkg/db"
)

func main() {
	// connect to database
	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api/health", healthHandler)

	http.HandleFunc("/api/events/", events.HandleEvents)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"healthy"}`)
}

