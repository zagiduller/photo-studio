package orders

import (
	"github.com/gorilla/mux"
	"net/http"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) GetPreparedMux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		// it can call only manager
		w.Write([]byte("Get orders"))
	}).Methods("GET")

	return r
}
