package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"onlineshop/internal/models"
)

type checkCart struct {
	Data []models.GetCart `json:"data"`
}

// @Summary Check your cart
// @Security BearerAuth
// @Tags cart
// @Description get cart by your id from database
// @Accept  json
// @Produce  json
// @Router /api/cart/ [get]
func (h *Handler) checkCart(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cartItems, err := h.services.GetCart(userid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if cartItems == nil {
		c.JSON(http.StatusOK, "your cart is empty")
	} else {
		c.JSON(http.StatusOK, checkCart{
			Data: cartItems,
		})
	}
}

type addProductToCart struct {
	ProductId int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// @Summary Add product to your cart
// @Security BearerAuth
// @Tags cart
// @Description add product to your cart by id and amount of product
// @Param balance body addProductToCart true "NewProduct"
// @Accept  json
// @Produce  json
// @Router /api/cart/add [post]
func (h *Handler) addProductToCart(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var product addProductToCart
	if err := c.BindJSON(&product); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Cart.AddProductToCart(user_id, product.ProductId, product.Quantity)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Make order
// @Security BearerAuth
// @Tags cart
// @Description create order from your cart
// @Accept  json
// @Produce  json
// @Router /api/cart/order [post]
func (h *Handler) makeOrder(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.Cart.MakeOrder(user_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "order is created",
	})
}
