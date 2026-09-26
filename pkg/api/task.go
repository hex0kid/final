package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeErrorText(w, "unsupported method")
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorText(w, "task id is required")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, err)
		return
	}
	if task.ID == "" {
		writeErrorText(w, "task id is required")
		return
	}
	if task.Title == "" {
		writeErrorText(w, "task title is required")
		return
	}
	if err := checkDate(&task); err != nil {
		writeError(w, err)
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorText(w, "task id is required")
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{})
}
