package components

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// @project photo-studio
// @created 27.08.2022

type App struct {
	ctx         context.Context
	log         *logrus.Entry
	name        string
	version     string
	collections []Component
}

func New(ctx context.Context, name, version string) *App {
	return &App{
		ctx:     ctx,
		name:    name,
		version: version,
		log: logrus.New().
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
		clog := c.GetLogger()
		start := time.Now()
		if err := c.Configure(a.ctx); err != nil {
			clog.Error(err)
			return err
		}
		clog.WithField("passed", time.Since(start)).
			Info("configured")
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
	alog.Info("starting")
	select {
	case <-a.ctx.Done():
		alog.WithFields(logrus.Fields{
			"end":    time.Now().UTC(),
			"passed": time.Since(start),
		}).Warn("stopped")
	}
	return nil
}
