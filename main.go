package main

import (
	"TimeToGoDoStuff/db"
	"net/http"
)

func main() {
	database := db.InitDB("tasks.db")
	db.CreateTable(database)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
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

		case "PUT":
			var task db.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			db.UpdateTask(database, task.ID, task.Description, task.Completed)
			w.WriteHeader(http.StatusOK)

		case "DELETE":
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "Missing task ID", http.StatusBadRequest)
				return
			}
			db.DeleteTask(database, id)
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
