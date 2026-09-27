package server

import (
	"net/http"
	"strings"

	"go_final_project/pkg/api"
)

func Run(port, password string) error {
	api.Init(password)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	if port == "" {
		port = "7540"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	return http.ListenAndServe(port, nil)
}
