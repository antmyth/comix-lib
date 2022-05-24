package dao

import (
	"fmt"
	"log"

	"github.com/antmyth/comix-lib/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection(cfg config.Config) *gorm.DB {
	dbCfg := cfg.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbCfg.Host, dbCfg.Username, dbCfg.Password, dbCfg.Name, dbCfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		log.Printf("Connected to db")
	}
	return db
}
