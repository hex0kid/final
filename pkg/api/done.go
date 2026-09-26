package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorText(w, "unsupported method")
		return
	}
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

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeError(w, err)
			return
		}
	}

	writeJSON(w, map[string]any{})
}
