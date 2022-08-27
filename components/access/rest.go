package access

import (
	"net/http"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) GetPreparedMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	return mux
}
