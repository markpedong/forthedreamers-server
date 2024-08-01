package routes

import (
	"net/http"

	"github.com/forthedreamers-server/controllers"
	"github.com/forthedreamers-server/middleware"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "API IS WORKING",
			"success": true,
			"status":  http.StatusOK,
		})
	})

	public := r.Group("/public")
	{
		public.POST("/login", controllers.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.Authentication)
	{
		api.POST("/uploadImage", controllers.UploadImage)
	}

	collections := r.Group("/collections")
	collections.Use(middleware.Authentication)
	{
		collections.POST("/add", controllers.AddCollection)
		collections.POST("/get", controllers.GetCollection)
		collections.POST("/update", controllers.UpdateCollection)
		collections.POST("/delete", controllers.DeleteCollection)
	}

	products := r.Group("/products")
	products.Use(middleware.Authentication)
	{
		products.POST("/add", controllers.AddProducts)
		products.POST("/get", controllers.GetProducts)
		products.POST("/update", controllers.UpdateProducts)
		products.POST("/delete", controllers.DeleteProducts)
	}

	users := r.Group("/users")
	users.Use(middleware.Authentication)
	{
		users.POST("/add", controllers.AddUsers)
		users.POST("/get", controllers.GetUsers)
		users.POST("/update", controllers.UpdateUsers)
	}

	variations := r.Group("/variations")
	variations.Use(middleware.Authentication)
	{
		variations.POST("/get", controllers.GetVariations)
		variations.POST("/update", controllers.UpdateVariations)
		variations.POST("/toggle", controllers.ToggleVariations)
	}
}
