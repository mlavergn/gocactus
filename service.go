package cactus

import (
	"net/http"
	"strconv"
)

// Service export
type Service struct {
}

// NewService export
func NewService() *Service {
	return &Service{}
}

func (id *Service) handler(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	resp.Header().Set("Server", "Cactus REST API "+Version)

	resp.Header().Set("Server", "Cactus REST API "+Version)
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "close")

	resp.WriteHeader(http.StatusOK)
}

// Start export
func (id *Service) Start() {
	hostPort := ":" + strconv.Itoa(8686)
	listener := http.Server{
		Addr: hostPort,
	}
	http.Handle("/", http.HandlerFunc(id.handler))

	listener.ListenAndServe()
}
