package components

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

// @project photo-studio
// @created 10.08.2022

var once = sync.Once{}
var db *gorm.DB

func GetDB() *gorm.DB {
	once.Do(func() {
		dbType := viper.GetString("components.db.type")
		host := viper.GetString("components.db.host")
		port := viper.GetString("components.db.port")
		user := viper.GetString("components.db.user")
		dbName := viper.GetString("components.db.dbname")
		password := viper.GetString("components.db.password")

		var dialect gorm.Dialector
		switch dbType {
		case "postgres":
			dialect = postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbName, password))
		default:
			log.Fatal("GetDB: Unknown database type")
		}

		_db, err := gorm.Open(dialect, &gorm.Config{})
		if err != nil {
			log.Fatal("GetDB: failed to connect database")
		}
		db = _db
	})

	return db
}
