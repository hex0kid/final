package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, err)
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, TasksResp{Tasks: tasks})
}
