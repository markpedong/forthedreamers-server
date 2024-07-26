package database

import (
	"fmt"
	"log"
	"os"

	"github.com/forthedreamers-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("DB_DSN")),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		})
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = DB.AutoMigrate(
		&models.Users{},
		&models.Collection{},
		&models.Product{},
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Println("--------------------Connected to Database---------------------")
}
