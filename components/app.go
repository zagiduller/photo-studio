package components

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// @project photo-studio
// @created 27.08.2022
// @author arthur

type (
	Startable interface {
		Start() error
		Stop() error
	}
	ComponentDependable interface {
		ConfigureDependencies([]Component)
	}

	App struct {
		ctx         context.Context
		log         *logrus.Entry
		name        string
		version     string
		startables  map[string]Startable
		collections []Component
	}
)

func New(ctx context.Context, name, version string) *App {
	return &App{
		ctx:        ctx,
		name:       name,
		version:    version,
		startables: map[string]Startable{},
		log: logrus.New().
			WithField("component", "app").
			WithField("app", name).
			WithField("version", version),
	}
}
func (a *App) add(c Component) error {
	for _, _c := range a.collections {
		if _c.GetName() == c.GetName() {
			return fmt.Errorf("component %s already exists", c.GetName())
		}
	}
	a.collections = append(a.collections, c)
	return nil
}

func (a *App) Add(cs ...Component) error {
	for _, _c := range cs {
		if err := a.add(_c); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Configure() error {
	for _, c := range a.collections {
		c.SetLogger(a.log)
		var (
			interfaces []string
			clog       = c.GetLogger()
			start      = time.Now()
		)

		if err := c.Configure(a.ctx); err != nil {
			clog.Error(err)
			return err
		}

		if depended, ok := c.(ComponentDependable); ok {
			interfaces = append(interfaces, "ComponentDependable")
			depended.ConfigureDependencies(a.collections)
		}
		if startable, ok := c.(Startable); ok {
			interfaces = append(interfaces, "Startable")
			a.startables[c.GetName()] = startable
		}

		clog.WithFields(logrus.Fields{
			"interfaces": interfaces,
			"passed":     time.Since(start),
		}).Info("component configured")
	}
	return nil
}

func (a *App) GetComponent(name string) Component {
	for _, c := range a.collections {
		if c.GetName() == name {
			return c
		}
	}
	return nil
}

func (a *App) Start() error {
	start := time.Now()
	alog := a.log.WithFields(logrus.Fields{
		"start": start.UTC(),
	})
	alog.Info("app starting")

	for name, s := range a.startables {
		go func(n string, _s Startable) {
			a.log.WithField("startable", n).Info("component started")
			if err := _s.Start(); err != nil {
				alog.Error(err)
			}
		}(name, s)
	}
	alog.WithField("passed", time.Since(start)).Info("app started")

	select {
	case <-a.ctx.Done():
		alog.WithFields(logrus.Fields{
			"end":    time.Now().UTC(),
			"passed": time.Since(start),
		}).Warn("stopped")
		for _, s := range a.startables {
			if err := s.Stop(); err != nil {
				alog.Error(err)
			}
		}
	}
	return nil
}
