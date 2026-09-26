package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	next := ""
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, taskDate) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, err)
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}
