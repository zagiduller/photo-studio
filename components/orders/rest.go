package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

// @project photo-studio
// @created 27.08.2022
// @author zagiduller

// TODO: Если пользователь не авторизован, то создавать нового пользователя и привязывать к заказу

func (s *Service) FillRouter(r, p *mux.Router) {
	r.HandleFunc("/", s.CreateOrderHandler).Methods(http.MethodPost, http.MethodOptions)
	p.HandleFunc("/", s.GetOrdersHandler).Methods(http.MethodGet, http.MethodOptions)
}

type GetOrdersResponse struct {
	Orders []*Order `json:"orders"`
}

func (s *Service) GetOrdersHandler(w http.ResponseWriter, _ *http.Request) {
	var response GetOrdersResponse
	orders, err := s.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Orders = orders
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var rePhone = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

type CreateOrderRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
}

func (r *CreateOrderRequest) Validate() error {
	email, phone := strings.TrimSpace(r.Email), strings.TrimSpace(r.Phone)
	if email == "" && phone == "" {
		return errors.New("CreateOrderRequest: email and phone are empty")
	}
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf("CreateOrderRequest: [%w]", err)
		}
	}
	if phone != "" && rePhone.MatchString(phone) {
		return errors.New("CreateOrderRequest: phone not valid")
	}

	return nil
}

type CreateOrderResponse struct {
	Order *Order `json:"order"`
}

func (s *Service) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var request, response = CreateOrderRequest{}, CreateOrderResponse{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		request.Name = "Инкогнито"
	}

	order, err := s.Create(
		// users
		strings.TrimSpace(request.Description),
		strings.TrimSpace(request.Name),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Order = order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
