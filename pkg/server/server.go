package server

import (
	"net/http"
	"os"
	"strings"

	"go_final_project/pkg/api"
)

func Run() error {
	api.Init()
	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	return http.ListenAndServe(port, nil)
}
