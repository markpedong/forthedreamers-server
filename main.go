package main

import (
	"log"

	"github.com/forthedreamers-server/cloudinary"
	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	database.ConnectDB()
}

func main() {
	cloudinary.Init()
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://forthedreamers-admin.vercel.app",
			"https://forthedreamers.vercel.app",
		},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	r.MaxMultipartMemory = 20 << 20
	routes.CreateRoutes(r)
	log.Fatal(r.Run(":6601"))
}
