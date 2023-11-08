package access

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zagiduller/photo-studio/components"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) FillRouter(r, p *mux.Router) {
	p.HandleFunc("/info", s.GetInfoHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", s.LoginPasswordHandler).Methods(http.MethodPost)
	r.HandleFunc("/signup", s.SingUpPasswordHandler).Methods(http.MethodPost)
}

var (
	ErrLoginIsBusy        = errors.New("Login is busy ")
	ErrInvalidClaims      = errors.New("Invalid claims ")
	ErrInvalidAccessToken = errors.New("Invalid access token ")
)

type (
	LoginPasswordRequest struct {
		Login    string
		Password string
	}
	LoginPasswordResponse struct {
		Token string
	}
	SignUpPasswordRequest struct {
		Login    string
		Password string
	}
	GetInfoHandlerResponse struct {
		ID    uint
		Type  string
		Login string
	}
)

// GetAccessByRequest

func (s *Service) GetLoginByRequest(r *http.Request) (*Login, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("token is empty")
	}
	jwt, err := s.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("GetLoginByRequest: [%w] ", err)
	}

	claims, ok := jwt.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("GetLoginByRequest: [%w] ", ErrInvalidClaims)
	}
	//claims.Subject
	loginID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("GetLoginByRequest: [%w] ", ErrInvalidClaims)
	}

	login, err := FindLoginByID(uint(loginID))
	if err != nil {
		return nil, fmt.Errorf("GetLoginByRequest: [%w] ", err)
	}
	return login, nil
}

// LoginPasswordHandler

func (s *Service) LoginPasswordHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := LoginPasswordRequest{}

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login, err := FindLoginByValue(components.GetDB(), request.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Check password
	if err := login.CheckPassword(request.Password); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	access, err := s.CreateTokenByLogin(components.GetDB(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := LoginPasswordResponse{
		Token: access.Token,
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(marshal)
	return
}

func (l *LoginPasswordRequest) Validate() error {
	if l == nil {
		return errors.New("LoginPasswordRequest: empty")
	}
	if len(strings.TrimSpace(l.Login)) < 3 {
		return errors.New("LoginPasswordRequest: login less than 3")
	}
	if len(strings.TrimSpace(l.Password)) < 6 {
		return errors.New("LoginPasswordRequest: password less than 6")
	}
	return nil
}

// SingUpPasswordHandler

func (s *Service) SingUpPasswordHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := SignUpPasswordRequest{}

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var access *Access

	if err := components.GetDB().Transaction(func(tx *gorm.DB) error {
		// create user
		user, err := s.users.NewUser()
		if err != nil {
			return err
		}
		user.SetDB(tx)
		if err := user.Save(); err != nil {
			return err
		}

		// create login
		// Check login is available
		if l, _ := FindLoginByValue(tx, request.Login); l != nil {
			return ErrLoginIsBusy
		}
		login := NewLogin(user.ID, LoginTypeEmail, request.Login)
		login.SetDB(tx)
		if err := login.Save(); err != nil {
			return err
		}

		// create password
		if err := CreateLoginPassword(tx, login, request.Password); err != nil {
			return err
		}

		_access, err := s.CreateTokenByLogin(tx, login)
		if err != nil {
			return err
		}
		access = _access
		return nil
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := LoginPasswordResponse{
		Token: access.Token,
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
	return
}

func (l *SignUpPasswordRequest) Validate() error {
	if l == nil {
		return errors.New("SignUpPasswordRequest: empty")
	}
	if len(strings.TrimSpace(l.Login)) < 3 {
		return errors.New("SignUpPasswordRequest: login less than 3")
	}
	if len(strings.TrimSpace(l.Password)) < 6 {
		return errors.New("SignUpPasswordRequest: password less than 6")
	}
	return nil
}

// GetInfoHandler

func (s *Service) GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	login, err := s.GetLoginByRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if login.GetUser() == nil {
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}
	marshal, err := json.Marshal(GetInfoHandlerResponse{
		ID:    login.GetUser().ID,
		Type:  login.Type,
		Login: login.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
	return
}
