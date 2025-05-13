package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Task represents a to-do item
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// In-memory storage for tasks
var tasks []Task
var nextID int = 1

// Handler to create a new task
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTask Task

	// Decode JSON request body into Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set unique ID and mark task as not completed
	newTask.ID = nextID
	newTask.Completed = false
	nextID++

	// Add task to in-memory slice
	tasks = append(tasks, newTask)

	// Return created task as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// Handler to list all tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return tasks as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Handler to get a task by ID
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Find the task by ID
	for _, task := range tasks {
		if task.ID == id {
			// Return task as JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	// Task not found
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Handler to delete a task by ID
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Find the task by ID and remove it
	for index, task := range tasks {
		if task.ID == id {
			// Remove task from the slice
			tasks = append(tasks[:index], tasks[index+1:]...)

			// Return success message
			w.WriteHeader(http.StatusNoContent) // No Content response for successful deletion
			return
		}
	}

	// Task not found
	http.Error(w, "Task not found", http.StatusNotFound)
}

// taskRouter handles all requests to /tasks and routes them to the appropriate handler
func taskRouter(w http.ResponseWriter, r *http.Request) {
	// Handle /tasks endpoint (with no ID)
	if r.URL.Path == "/tasks" {
		switch r.Method {
		case http.MethodGet:
			getTasksHandler(w, r)
		case http.MethodPost:
			createTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle /tasks/{id} endpoints
	if strings.HasPrefix(r.URL.Path, "/tasks/") {
		switch r.Method {
		case http.MethodGet:
			getTaskByIDHandler(w, r)
		case http.MethodDelete:
			deleteTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// If we got here, the path is not recognized
	http.NotFound(w, r)
}

func main() {
	// Welcome route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only respond to exact root path
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "Welcome to the Go To-Do CRUD API!")
	})

	// Register single router for all /tasks routes
	http.HandleFunc("/tasks", taskRouter)
	http.HandleFunc("/tasks/", taskRouter)

	// Start server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
