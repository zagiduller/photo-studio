package components

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

// @project photo-studio
// @created 10.08.2022

var (
	ErrorCodeDbIsNil = errors.New("components: db is nil")
)

type Component interface {
	GetName() string
	Configure(ctx context.Context) error
	SetLogger(*logrus.Entry)
	GetLogger() *logrus.Entry
}

func Std(name string) Default {
	return Default{
		name: name,
	}
}

type Default struct {
	Ctx  context.Context
	log  *logrus.Entry
	name string
}

func (d *Default) GetName() string {
	if d.name == "" {
		return "unknown"
	}
	return d.name
}

func (d *Default) Configure() error {
	return nil
}

func (d *Default) SetLogger(log *logrus.Entry) {
	d.log = log.WithField("component", d.GetName())
}

func (d *Default) GetLogger() *logrus.Entry {
	if d.log == nil {
		d.SetLogger(logrus.NewEntry(logrus.New()))
	}
	return d.log
}
