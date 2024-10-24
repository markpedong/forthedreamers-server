package main

import (
	"fmt"
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
	goth.UseProviders(google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), fmt.Sprintf("%s/public/googleCallback", os.Getenv("WEB_CALLBACK")), "email", "profile"))
}

func main() {
	cloudinary.Init()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://forthedreamers-admin.vercel.app",
			"https://forthedreamers.vercel.app",
			os.Getenv("WEB_DOMAIN"),
		},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	r.MaxMultipartMemory = 20 << 20

	routes.CreateRoutes(r)
	log.Fatal(r.Run(":6601"))
}
