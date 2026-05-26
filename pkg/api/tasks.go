package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

const tasksLimit = 50

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	tasks, err := db.Tasks(tasksLimit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, tasksResp{Tasks: tasks})
}
