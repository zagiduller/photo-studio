package main

import (
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

func init() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	collection := []components.Component{
		auth.New(),
		users.New(),
		orders.New(),
	}
	for _, c := range collection {
		start := time.Now()
		if err := c.Configure(); err != nil {
			log.Fatal(err)
		}
		c.GetLogger().
			WithField("passed", time.Since(start)).
			Info("Configured")
	}

	// Главная страница. Создание заказа
	// Храним информацию о заказе в БД
	// Менеджер получает уведомление о новом заказе

	// Менеджер видит заказ
	// Менеджер может менять статус заказа

	// Если заказ меняет статус - создаем пользователя из заказа
	// Под заказ выделяется minio контейнер для файлов заказа
}
