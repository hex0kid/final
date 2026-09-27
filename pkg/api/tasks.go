package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeErrorText(w, http.StatusMethodNotAllowed, "unsupported method")
		return
	}

	var (
		tasks []*db.Task
		err   error
	)
	if search := r.FormValue("search"); search != "" {
		tasks, err = db.SearchTasks(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
