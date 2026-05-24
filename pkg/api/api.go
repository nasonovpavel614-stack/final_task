package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Init регистрирует обработчики API.
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/signin", signinHandler)
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", Auth(taskHandler))
	mux.HandleFunc("/api/task/done", Auth(doneTaskHandler))
	mux.HandleFunc("/api/tasks", Auth(tasksHandler))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			fmt.Fprint(w, errInvalidDate.Error())
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	next, err := NextDate(now, date, repeat)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, next)
}
