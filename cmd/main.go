package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zagiduller/photo-studio/components"
	"github.com/zagiduller/photo-studio/components/auth"
	"github.com/zagiduller/photo-studio/components/orders"
	"github.com/zagiduller/photo-studio/components/server"
	"github.com/zagiduller/photo-studio/components/users"
	"os"
	"os/signal"
)

// @project photo-studio
// @created 27.07.2022

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app := components.New(
		ctx,
		"photo-studio",
		"0.0.1",
	)
	must(app.Add(
		auth.New(),
		users.New(),
		orders.New(),
		server.New(),
	))
	must(app.Configure())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.WithField("signal", sig.String()).Warn("shutting down...")
			cancel()
		}
	}()

	log.Error(app.Start())

	// Главная страница. Создание заказа
	// Храним информацию о заказе в БД
	// Менеджер получает уведомление о новом заказе

	// Менеджер видит заказ
	// Менеджер может менять статус заказа

	// Если заказ меняет статус - создаем пользователя из заказа
	// Под заказ выделяется minio контейнер для файлов заказа
}

func init() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
