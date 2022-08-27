package auth

import "net/http"

// @project photo-studio
// @created 27.08.2022

type Rest struct {
	http.ServeMux
	Service
}
