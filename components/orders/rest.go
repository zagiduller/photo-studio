package orders

import (
	"github.com/gorilla/mux"
	"net/http"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) GetPreparedMux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("New order"))
	})
	return r
}
