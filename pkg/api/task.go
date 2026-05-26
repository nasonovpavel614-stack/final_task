package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go_final_project/pkg/db"
)

var (
	errNoID         = errors.New("Не указан идентификатор")
	errTaskNotFound = errors.New("Задача не найдена")
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(body, &task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": errNoID.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": errEmptyTitle.Error()})
		return
	}

	if err = checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err = db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": errTaskNotFound.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": errNoID.Error()})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": errTaskNotFound.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}
