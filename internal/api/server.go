package api

import (
	"encoding/json"
	"net/http"

	"github.com/ruhulfbr/mini-k8s/orchestrator"
)

type Server struct {
	ctrl *orchestrator.Controller
}

func NewServer(ctrl *orchestrator.Controller) *Server {
	return &Server{ctrl: ctrl}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/scale", s.handleScale)

	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleScale(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req orchestrator.ScaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.ctrl.Scale(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write([]byte("scaling operation queued"))
	if err != nil {
		return
	}
}
