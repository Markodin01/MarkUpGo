package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=12345 dbname=marekgaj sslmode=disable"

	storage, err := NewStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	err = storage.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new router using the Gorilla Mux router
	router := mux.NewRouter()

	// Define your endpoints
	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllTasks(w, r, storage)
		case http.MethodPost:
			CreateTask(w, r, storage)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetTask(w, r, storage)
		case http.MethodPut:
			UpdateTask(w, r, storage)
		case http.MethodDelete:
			DeleteTask(w, r, storage)

		}
	})

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
	print("The sever started")

}

// Define your endpoint handlers

func GetAllTasks(w http.ResponseWriter, r *http.Request, storage *Storage) {
	// Retrieve tasks from your storage or database
	tasks, err := storage.GetAllTasks()
	if err != nil {
		// Handle the error and return an appropriate response
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Convert tasks to JSON
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		// Handle the error and return an appropriate response
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonData)
}

func CreateTask(w http.ResponseWriter, r *http.Request, storage *Storage) {
	// Parse the request body to get the task data
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Perform any necessary validation on the task data
	if task.Title == "" {
		http.Error(w, "Task title is required", http.StatusBadRequest)
		return
	}

	// Save the task to the storage or database
	err = storage.CreateTask(&task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Set the response status to 201 Created
	w.WriteHeader(http.StatusCreated)
}

func GetTask(w http.ResponseWriter, r *http.Request, storage *Storage) {

	var taskId int
	var task Task

	err := json.Decoder(r.Body).Decode(&taskId)

	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Perform any necessary validation on the task data
	if taskId == 0 {
		http.Error(w, "Task id is required", http.StatusBadRequest)
		return
	}

	// Save the task to the storage or database
	task, err = storage.GetTask(&taskId)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Set the response status to 201 Created
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func UpdateTask(w http.ResponseWriter, r *http.Request, storage *Storage) {
	// Handle the PUT /tasks/{id} endpoint
	// ...
}

func DeleteTask(w http.ResponseWriter, r *http.Request, storage *Storage) {
	// Handle the DELETE /tasks/{id} endpoint
	// ...
}
