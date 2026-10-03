package main

import (
	"log"
	"net/http"
	"os"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@postgres:5432/tasksdb?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port %s", serverPort)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Printf("Connected to database %s", databaseURL)

	taskStore := database.NewTaskStore(db)

	handler := handlers.NewHandler(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTasks, "GET"))
	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, "POST"))
	mux.HandleFunc("/tasks/", taskIDHandler(handler))

	loggedMux := loggingMiddleware(mux)

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		handlerFunc(w, r)
	}
}

func taskIDHandler(handler *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
