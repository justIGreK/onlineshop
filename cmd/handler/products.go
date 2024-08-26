package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"onlineshop/internal/models"
)

type Product interface {
	CreateProduct(name string, cost float64, description string, amount int) (int, error)
	GetProductList() ([]models.Product, error)
	GetProductById(id int) (models.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(id int, product models.UpdateProduct) error
}

type addProduct struct {
	Name        string  `json:"name" binding:"required"`
	Cost        float64 `json:"cost" binding:"required,gt=0"`
	Description string  `json:"description"`
	Amount      int     `json:"amount" binding:"required,gt=0"`
}

// @Summary Add product
// @Security BearerAuth
// @Tags products
// @Description add product to database
// @Param name query string true "product name"
// @Param cost query float64 true "price for product"
// @Param description query string false "product description"
// @Param amount query int true "amount of product"
// @Accept  json
// @Produce  json
// @Router /api/products/ [post]
func (h *Handler) addProduct(c *gin.Context) {
	userRole, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userRole != AdministratorRole {
		newErrorResponse(c, http.StatusUnauthorized, "You have 	no permission to do this")
	}
	strCost := c.Query("cost")
	cost, err := strconv.ParseFloat(strCost, 64)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid cost")
		return
	}
	strAmount := c.Query("amount")
	amount, err := strconv.Atoi(strAmount)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid amount")
		return
	}
	input := addProduct{
		Name:        c.Query("name"),
		Cost:        float64(cost),
		Description: c.Query("description"),
		Amount:      amount,
	}
	id, err := h.Prod.CreateProduct(input.Name, input.Cost, input.Description, input.Amount)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":         id,
		"name":       input.Name,
		"cost":       input.Cost,
		"desription": input.Description,
		"amount":     input.Amount,
	})
}

type getProductListResponse struct {
	Data []models.Product `json:"data"`
}

// @Summary Get productlist
// @Security BearerAuth
// @Tags products
// @Description get product list from database
// @Accept  json
// @Produce  json
// @Router /api/products/ [get]
func (h *Handler) getProductList(c *gin.Context) {
	products, err := h.Prod.GetProductList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getProductListResponse{
		Data: products,
	})
}

// @Summary Get product by id
// @Security BearerAuth
// @Tags products
// @Description get product by id from database
// @Param id path int  true  "Product ID"
// @Accept  json
// @Produce  json
// @Router /api/products/{id} [get]
func (h *Handler) getProduct(c *gin.Context) {
	searchId := c.Param("id")
	id, err := strconv.Atoi(searchId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	product, err := h.Prod.GetProductById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.Product{
		Id:          product.Id,
		Name:        product.Name,
		Cost:        product.Cost,
		Description: product.Description,
		Amount:      product.Amount,
		IsActice:    product.IsActice,
	})
}

// @Summary  Change info about product
// @Security BearerAuth
// @Tags products
// @Description change name, cost or amount of product from database
// @Param data body models.UpdateProduct true "NewProduct"
// @Param id path int  true  "Product ID"
// @Accept  json
// @Produce  json
// @Router /api/products/{id} [put]
func (h *Handler) changeProduct(c *gin.Context) {
	userRole, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userRole != AdministratorRole {
		newErrorResponse(c, http.StatusUnauthorized, "You have no permission to do this")
	}
	var inputProduct models.UpdateProduct
	searchId := c.Param("id")
	id, err := strconv.Atoi(searchId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.BindJSON(&inputProduct); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Prod.UpdateProduct(id, inputProduct); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Delete product by id
// @Security BearerAuth
// @Tags products
// @Description delete product by id from database
// @Param id path int  true  "Product ID"
// @Accept  json
// @Produce  json
// @Router /api/products/{id} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	userRole, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userRole != AdministratorRole {
		newErrorResponse(c, http.StatusUnauthorized, "You have no permission to do this")
	}
	deleteId := c.Param("id")
	id, err := strconv.Atoi(deleteId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.Prod.DeleteProduct(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
