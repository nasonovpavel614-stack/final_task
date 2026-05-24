package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": errNoID.Error()})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": errTaskNotFound.Error()})
		return
	}

	if task.Repeat == "" {
		if err = db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": errTaskNotFound.Error()})
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err = db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": errTaskNotFound.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}
