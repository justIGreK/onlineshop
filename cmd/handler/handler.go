package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "onlineshop/docs"
	"onlineshop/internal/service"
	"onlineshop/internal/ws"
)

const (
	AdministratorRole string = "admin"
	CustomerRole      string = "customer"
)

type Handler struct {
	Auth Authorization
	User UserList
	Prod Product
	Crt  Cart
	Ord  Order
}

func NewHandler(auth *service.AuthService,
	user *service.UserService,
	prod *service.ProdService,
	crt *service.CartService,
	ord *service.OrderService) *Handler {
	return &Handler{
		Auth: auth,
		User: user,
		Prod: prod,
		Crt:  crt,
		Ord:  ord,
	}
}

func InitRoutes(shopHandler *Handler, wsHandler *ws.Handler) *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", shopHandler.signUp)
		auth.POST("/sign-in", shopHandler.signIn)
	}

	api := router.Group("/api", shopHandler.userIdentity)
	{
		ws := api.Group("/ws")
		{
			ws.POST("/createRoom", wsHandler.CreateRoom)
			ws.GET("/joinRoom", wsHandler.JoinRoom)
			ws.GET("/getRooms", wsHandler.GetRooms)
			ws.GET("/getHistory/", wsHandler.GetHistory)
		}
		users := api.Group("/users")
		{
			users.GET("/", shopHandler.getUserList)
			users.GET("/:id", shopHandler.getUser)
			users.PUT("/:id", shopHandler.changeBalance)
			users.DELETE("/", shopHandler.deleteUser)
			users.POST("link/:service", shopHandler.linkAcc)
		}
		products := api.Group("/products")
		{
			products.POST("/", shopHandler.addProduct)
			products.GET("/", shopHandler.getProductList)
			products.GET("/:id", shopHandler.getProduct)
			products.PUT("/:id", shopHandler.changeProduct)
			products.DELETE("/:id", shopHandler.deleteProduct)
		}
		cart := api.Group("/cart")
		{
			cart.GET("/", shopHandler.checkCart)
			cart.POST("/order", shopHandler.makeOrder)
			cart.POST("/add", shopHandler.addProductToCart)
		}
		order := api.Group("/orders")
		{
			order.GET("/", shopHandler.checkOrders)
			order.GET("/:id", shopHandler.getOrder)
		}
	}
	return router
}
