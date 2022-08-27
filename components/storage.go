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

var db *gorm.DB
var once = sync.Once{}

func GetDB() *gorm.DB {
	once.Do(func() {
		dbType := viper.GetString("components.db.type")
		host := viper.GetString("components.db.host")
		port := viper.GetString("components.db.port")
		user := viper.GetString("components.db.user")
		dbName := viper.GetString("components.db.database")
		password := viper.GetString("components.db.password")

		var dialect gorm.Dialector
		switch dbType {
		case "postgres":
			dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbName)
			dialect = postgres.Open(dsn)
		default:
			log.Fatal("GetDB: Unknown database type")
		}

		_db, err := gorm.Open(dialect)
		if err != nil {
			log.Fatalf("GetDB: failed to connect database: %s ", err)
		}
		db = _db
		log.WithFields(log.Fields{
			"port": port,
			"host": host,
			"user": user,
			"db":   dbName,
			"type": dbType,
		}).Infof("connected")
	})

	return db
}
