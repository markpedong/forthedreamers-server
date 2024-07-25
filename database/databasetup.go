package database

import (
	"fmt"
	"log"
	"os"

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
	// &models.Users{},
	// &models.Links{},
	// &models.Services{},
	// &models.Messages{},
	// &models.Portfolios{},
	// &models.Blogs{},
	// &models.Education{},
	// &models.Experiences{},
	// &models.ExpSkill{},
	// &models.EduSkill{},
	// &models.Testimonials{},
	// &models.Website{},
	// &models.Color{},
	// &models.ExpDesc{},
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Println("--------------------Connected to Database---------------------")
}
