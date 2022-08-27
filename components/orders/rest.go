package orders

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// @project photo-studio
// @created 27.08.2022

func (s *Service) FillRouter(r, p *mux.Router) {
	r.HandleFunc("/", s.CreateOrderHandler).Methods("POST")
	p.HandleFunc("/", s.GetOrdersHandler).Methods("GET")
}

type GetOrdersResponse struct {
	Orders []*Order `json:"orders"`
}

func (s *Service) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
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

type CreateOrderRequest struct {
	Description string `json:"description"`
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
	order, err := s.Create(request.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Order = order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
