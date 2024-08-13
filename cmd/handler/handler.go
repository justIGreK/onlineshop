package handler

import (
	"onlineshop/internal/service"

	_ "onlineshop/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getUserList)
			users.GET("/:id", h.getUser)
			users.PUT("/:id", h.changeBalance)
			users.DELETE("/", h.deleteUser)
		}
		products := api.Group("/products")
		{
			products.POST("/", h.addProduct)
			products.GET("/", h.getProductList)
			products.GET("/:id", h.getProduct)
			products.PUT("/:id", h.changeProduct)
			products.DELETE("/:id", h.deleteProduct)
		}
		cart := api.Group("/cart")
		{
			cart.GET("/", h.checkCart)
			cart.POST("/order", h.makeOrder)
			cart.POST("/add", h.addProductToCart)
		}
		order := api.Group("/orders")
		{
			order.GET("/", h.checkOrders)
			order.GET("/:id", h.getOrder)
		}
	}
	return router
}
