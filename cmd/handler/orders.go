package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"onlineshop/internal/models"
)

type checkOrder struct {
	Data []models.GetOrder `json:"data"`
}

// @Summary Get orderList
// @Security BearerAuth
// @Tags orders
// @Description get your orders list from database
// @Accept  json
// @Produce  json
// @Router /api/orders/ [get]
func (h *Handler) checkOrders(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	orders, err := h.services.GetOrderList(userid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if orders == nil {
		c.JSON(http.StatusOK, "You have no orders")
	} else {
		c.JSON(http.StatusOK, checkOrder{
			Data: orders,
		})
	}
}

type getOrderDetails struct {
	Data       models.GetOrder
	OrderItems []models.OrderItems
}

// @Summary Get order by id
// @Security BearerAuth
// @Tags orders
// @Description get order by id from database
// @Param id path int  true  "order ID"
// @Accept  json
// @Produce  json
// @Router /api/orders/{id} [get]
func (h *Handler) getOrder(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	searchId := c.Param("id")
	orderID, err := strconv.Atoi(searchId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	order, orderItems, err := h.services.GetOrderDetails(userID, orderID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getOrderDetails{
		Data:       order,
		OrderItems: orderItems,
	})
}
