package access

import (
	"github.com/gorilla/mux"
	"net/http"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) FillRouter(r, p *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
}
