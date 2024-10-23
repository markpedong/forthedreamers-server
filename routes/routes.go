package routes

import (
	"net/http"

	"github.com/forthedreamers-server/controllers"
	"github.com/forthedreamers-server/middleware"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API IS WORKING",
			"success": true,
			"status":  http.StatusOK,
		})
	})

	public := r.Group("/public")
	{
		public.POST("/login", controllers.Login)
		public.GET("/googleLogin", controllers.GoogleLogin)
		public.GET("/googleCallback", controllers.GoogleCallback)
		public.POST("/signup", controllers.SignUp)
		public.POST("/requestEmailOTP", controllers.RequestEmailOTP)
		public.POST("/verifyOTP", controllers.VerifyOTP)
		public.POST("/setNewPassword", controllers.SetNewPassword)
		public.GET("/collections", controllers.PublicCollections)
		public.GET("/collectionsByID", controllers.GetCollectionByID)
		public.GET("/products", controllers.PublicProducts)
		public.GET("/products/details", controllers.PublicProductDetails)
		public.GET("/products/variations", controllers.PublicVariations)
		public.GET("/website", controllers.PublicWebsite)
		public.GET("/testimonials", controllers.PublicTestimonials)
	}

	api := r.Group("/api")
	api.Use(middleware.Authentication)
	{
		api.POST("/uploadImage", controllers.UploadImage)
	}

	addresses := r.Group("/address")
	addresses.Use(middleware.Authentication)
	{
		addresses.POST("/add", controllers.AddAddress)
		addresses.GET("/get", controllers.GetAddress)
		addresses.POST("/update", controllers.UpdateAddress)
		addresses.POST("/delete", controllers.DeleteAddress)
	}

	carts := r.Group("/carts")
	carts.Use(middleware.Authentication)
	{
		carts.POST("/add", controllers.AddCartItem)
		carts.POST("/addQuantity", controllers.AddCartItemQuantity)
		carts.GET("/get", controllers.GetCart)
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
		products.POST("/finish", controllers.FinishOrder)
	}

	reviews := r.Group("/reviews")
	reviews.Use(middleware.Authentication)
	{
		reviews.POST("/add", controllers.AddOrderReview)
		reviews.GET("/get", controllers.GetUserReview)
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
		users.POST("/info", controllers.GetUserInfo)
		users.POST("/get", controllers.GetUsers)
		users.POST("/update", controllers.UpdateUsers)
		users.POST("/delete", controllers.DeleteUsers)
		users.POST("/toggle", controllers.ToggleUsers)
		users.POST("/checkout", controllers.CheckoutOrder)
		users.GET("/orders", controllers.GetOrders)
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
