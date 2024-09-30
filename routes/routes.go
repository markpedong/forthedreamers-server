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
		public.POST("/collections", controllers.PublicCollections)
		public.POST("/collectionsByID", controllers.GetCollectionByID)
		public.POST("/products", controllers.PublicProducts)
		public.POST("/products/details", controllers.PublicProductDetails)
		public.POST("/products/variations", controllers.PublicVariations)
		public.POST("/website", controllers.PublicWebsite)
		public.POST("/testimonials", controllers.PublicTestimonials)
	}

	api := r.Group("/api")
	api.Use(middleware.Authentication)
	{
		api.POST("/uploadImage", controllers.UploadImage)
	}

	carts := r.Group("/carts")
	carts.Use(middleware.Authentication)
	{
		carts.POST("/add", controllers.AddCartItem)
		carts.POST("/get", controllers.GetCart)
		carts.POST("/delete", controllers.DeleteCartItem)
	}

	collections := r.Group("/collections")
	collections.Use(middleware.Authentication)
	{
		collections.POST("/add", controllers.AddCollection)
		collections.POST("/get", controllers.GetCollection)
		collections.POST("/update", controllers.UpdateCollection)
		collections.POST("/delete", controllers.DeleteCollection)
		collections.POST("/toggle", controllers.ToggleCollections)
	}

	products := r.Group("/products")
	products.Use(middleware.Authentication)
	{
		products.POST("/add", controllers.AddProducts)
		products.POST("/get", controllers.GetProducts)
		products.POST("/update", controllers.UpdateProducts)
		products.POST("/delete", controllers.DeleteProducts)
		products.POST("/toggle", controllers.ToggleProducts)
	}

	testimonials := r.Group("/testimonials")
	testimonials.Use(middleware.Authentication)
	{
		testimonials.POST("/add", controllers.AddTestimonials)
		testimonials.POST("/get", controllers.GetTestimonials)
	}
	users := r.Group("/users")
	users.Use(middleware.Authentication)
	{
		users.POST("/add", controllers.AddUsers)
		users.POST("/get", controllers.GetUsers)
		users.POST("/update", controllers.UpdateUsers)
		users.POST("/delete", controllers.DeleteUsers)
		users.POST("/toggle", controllers.ToggleUsers)
	}

	variations := r.Group("/variations")
	variations.Use(middleware.Authentication)
	{
		variations.POST("/get", controllers.GetVariations)
		variations.POST("/add", controllers.AddVariations)
		variations.POST("/update", controllers.UpdateVariations)
		variations.POST("/delete", controllers.DeleteVariations)
		variations.POST("/toggle", controllers.ToggleVariations)
	}

	website := r.Group("/website")
	website.Use(middleware.Authentication)
	{
		website.GET("/get", controllers.GetWebsiteData)
		website.POST("/update", controllers.UpdateWebsiteData)
	}
}
