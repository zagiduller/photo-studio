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
	MuxPreparable interface {
		GetPreparedMux() http.Handler
	}

	Service struct {
		port       int
		address    string
		rootRouter *mux.Router
		server     http.Server
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
	s.rootRouter = mux.NewRouter()

	s.port = viper.GetInt("components.server.port")
	if s.port == 0 {
		s.port = 8080
	}
	s.address = fmt.Sprintf("http://localhost:%d", s.port)
	s.server = http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.rootRouter,
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
		preparable, ok := c.(MuxPreparable)
		if ok {
			s.UsePreparable(c.GetName(), preparable)
		}
	}
}

func (s *Service) UsePreparable(prefix string, preparable MuxPreparable) {
	path := fmt.Sprintf("/%s/", prefix)
	sub := s.rootRouter.PathPrefix(path).Subrouter()
	sub.Handle("/", preparable.GetPreparedMux())
	s.GetLogger().WithField("address", s.address+path).Info("use prepared mux")
}
