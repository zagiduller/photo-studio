package access

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) FillRouter(r, _ *mux.Router) {
	r.HandleFunc("/login", s.LoginHandler)
}

type LoginRequest struct {
	Login    string
	Password string
}

func (l *LoginRequest) Validate() error {
	if l == nil {
		return errors.New("LoginRequest: empty")
	}
	if len(strings.TrimSpace(l.Login)) < 4 {
		return errors.New("LoginRequest: login less than 4")
	}
	if len(strings.TrimSpace(l.Password)) < 6 {
		return errors.New("LoginRequest: password less than 6")
	}
	return nil
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := LoginRequest{}

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
