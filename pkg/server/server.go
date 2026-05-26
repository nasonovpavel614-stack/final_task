package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/pkg/api"
)

const defaultPort = "7540"

// Port возвращает порт сервера: из TODO_PORT или значение по умолчанию.
func Port() string {
	if p := os.Getenv("TODO_PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			return p
		}
	}
	return defaultPort
}

// Run запускает HTTP-сервер, отдающий статические файлы из webDir.
func Run(webDir string) error {
	mux := http.NewServeMux()
	api.Init(mux)
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%s", Port())
	log.Printf("server listening on http://localhost%s", addr)
	return http.ListenAndServe(addr, mux)
}
