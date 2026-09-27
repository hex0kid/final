package api

import (
	"encoding/json"
	"errors"
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
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		writeErrorText(w, http.StatusMethodNotAllowed, "unsupported method")
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorText(w, http.StatusBadRequest, "task id is required")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if task.ID == "" {
		writeErrorText(w, http.StatusBadRequest, "task id is required")
		return
	}
	if task.Title == "" {
		writeErrorText(w, http.StatusBadRequest, "task title is required")
		return
	}
	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeErrorText(w, http.StatusBadRequest, "task id is required")
		return
	}
	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
