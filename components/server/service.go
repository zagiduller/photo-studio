package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/zagiduller/photo-studio/components"
	"net/http"
)

// @project photo-studio
// @created 27.08.2022

type (
	RouterFillable interface {
		FillRouter(r, p *mux.Router)
	}

	Service struct {
		port     int
		address  string
		rRoot    *mux.Router
		rPublic  *mux.Router
		rPrivate *mux.Router
		server   http.Server
		components.Default
	}
)

func New() *Service {
	return &Service{
		Default: components.DefaultComponent("server"),
	}
}

func (s *Service) Configure(ctx context.Context) error {
	s.Default.Ctx = ctx

	s.rRoot = mux.NewRouter()
	s.rRoot.Use(AllowCORS)

	s.rPublic = s.rRoot.PathPrefix("/").Subrouter()
	s.rPrivate = s.rRoot.PathPrefix("/-/").Subrouter()

	s.port = viper.GetInt("components.server.port")
	if s.port == 0 {
		s.port = 8080
	}
	s.address = fmt.Sprintf("http://localhost:%d", s.port)
	s.server = http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.rRoot,
	}
	return nil
}

func (s *Service) Start() error {
	s.GetLogger().WithField("address", s.address).Info("listen")
	return s.server.ListenAndServe()
}

func (s *Service) Stop() error {
	return s.server.Close()
}

func (s *Service) ConfigureDependencies(components []components.Component) {
	for _, c := range components {
		fillable, ok := c.(RouterFillable)
		if ok {
			s.fillRouter(c.GetName(), fillable)
		}
	}
}

func (s *Service) fillRouter(prefix string, preparable RouterFillable) {
	path := fmt.Sprintf("/%s/", prefix)
	preparable.FillRouter(
		s.rPublic.PathPrefix(path).Subrouter(),
		s.rPrivate.PathPrefix(path).Subrouter(),
	)
	s.GetLogger().WithField("address", s.address+path).Info("use prepared mux")
}
