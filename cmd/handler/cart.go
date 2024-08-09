package handler

import (
	"net/http"
	"onlineshop/internal/models"

	"github.com/gin-gonic/gin"
)

type checkCart struct {
	Data []models.GetCart `json:"data"`
}

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

	c.JSON(http.StatusOK, checkCart{
		Data: cartItems,
	})

}

type addProductToCart struct {
	ProductId int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

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
