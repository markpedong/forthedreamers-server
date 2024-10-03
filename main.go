package main

import (
	"log"
	"os"

	"github.com/forthedreamers-server/cloudinary"
	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func init() {
	database.ConnectDB()
}

func main() {
	cloudinary.Init()
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://forthedreamers-admin.vercel.app",
			"https://forthedreamers.vercel.app",
			"http://localhost:6600",
		},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	r.MaxMultipartMemory = 20 << 20

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/login?otp="),
	)
	routes.CreateRoutes(r)
	log.Fatal(r.Run(":6601"))
}
