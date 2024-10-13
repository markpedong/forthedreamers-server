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

	if err := DB.AutoMigrate(
		&models.Users{},
		&models.Collection{},
		&models.Product{},
		&models.ProductVariation{},
		&models.WebsiteData{},
		&models.Testimonials{},
		&models.UserCart{},
		&models.AddressItem{},
	); err != nil {
		log.Fatal(err.Error())
		return
	}

	if err2 := DB.AutoMigrate(
		&models.CartItem{},
		&models.OrderItem{},
	); err2 != nil {
		log.Fatal(err2.Error())
		return
	}

	fmt.Println("--------------------Connected to Database---------------------")
}
