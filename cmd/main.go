package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"photostudio/components"
	"photostudio/components/auth"
	"photostudio/components/orders"
	"photostudio/components/users"
	"time"
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
	))
	must(app.Configure())

	go app.Start()

	time.Sleep(1 * time.Second)
	// Остановка приложения
	cancel()
	time.Sleep(1 * time.Second)

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
