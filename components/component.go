package components

import (
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
	Configure() error
	GetLogger() *logrus.Entry
}

func New(name string) Default {
	return Default{
		name: name,
	}
}

type Default struct {
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

func (d *Default) GetLogger() *logrus.Entry {
	return logrus.StandardLogger().WithField("component", d.name)
}
