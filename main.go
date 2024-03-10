package main

import (
	"TimeToGoDoStuff/db"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func setContentTypeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, ".css") {
            w.Header().Set("Content-Type", "text/css")
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", setContentTypeMiddleware(fs))

	database := db.InitDB("tasks.db")
	db.CreateTable(database)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case "GET":
			tasks := db.GetAllTasks(database)
			tasksJSON, err := json.Marshal(tasks)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(tasksJSON)
		case "POST":
			var task db.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			db.CreateTask(database, task.Description)
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "PUT":
			var task db.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			db.UpdateTask(database, id, task.Description, task.Completed) // Use the extracted 'id' directly
			w.WriteHeader(http.StatusOK)
		case "DELETE":
			db.DeleteTask(database, strconv.Itoa(id)) // Convert 'id' to string if necessary
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	})

	http.ListenAndServe(":8080", mux)
}

