package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishowsagar/go-pizza-shop/internal/models"
)

//@ types declaration

type OrderFormData struct {
	PizzaTypes []string
	PizzaSizes []string
}

// type of data incoming from client - to decode it properly into store vars
type OrderRequest struct {
	Name string `form:"name" binding:"required,min=2,max=70"` //! specifying tags to check what it would be in form and binding validator on in for gin
	Phone string `form:"phone" binding:"required,min=10,max=20"` 
	Address string `form:"address" binding:"required,min=10,max=20"` 
	Sizes []string `form:"size" binding:"required,min=1,dive,pizza_valid_size"` //! binding declared validation type to let gin check for this
	PizzaTypes []string `form:"pizza" binding:"required,min=1,dive,pizza_valid_type"` 
	Instructions []string `form:"instructions" binding:"max=140"`
}


//@ methods that belongs to type -> *Controller
func(c *Controller) ServeNewOrderRequest(ctx *gin.Context) {
	//* creating template passing this data struct
	// todo --> create templalte "order.tmpl" to send res
	ctx.HTML(http.StatusOK,"order.tmpl",OrderFormData{
		PizzaTypes: models.PizzaTypes,
		PizzaSizes: models.PizzaSizes ,
	})

}

// we would be writing respons to client using ctx's methods like Json
func(c *Controller) HandleNewOrderPost(ctx *gin.Context) {

	var orderReq OrderRequest

	err := ctx.ShouldBind(&orderReq)
	
	// if hit any error
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			// gin.H sends a custom data struct like resp to client in form of a map
			"error" : err.Error(),
		})
		return
	}

	// otherwise we would be sending post req to make an order
	orderItemsReq := make([]models.OrderItem,len(orderReq.Sizes))

	// iterating over it to maje OrderItem form current iteration
	for i:= range orderItemsReq {
		// iterating over each element recieved and setting it in orderItemsReq slice 
		orderItemsReq[i] = models.OrderItem{
			Size: orderReq.Sizes[i],
			Pizza: orderReq.PizzaTypes[i],
			Instructions: orderReq.Instructions[i],
		} 
	}

	order := models.Order{
		CustomerName: orderReq.Name,
		Phone: orderReq.Phone,
		Address: orderReq.Address,
		Status: models.OrderStatus[0], //* setting placed on 1st element in slice being 0th item in index
		Items: orderItemsReq,
	}

	err = c.orderModel.CreateOrder(&order)
	if err != nil {
		slog.Error("failed to create an order -%v","error",err)
		ctx.String(http.StatusInternalServerError,"internal server error")
		return
	}

	slog.Info("order is successfully placed","orderID",order.ID,"customer",order.CustomerName)//! just pass key val pairs
	ctx.Redirect(http.StatusSeeOther,"/customer/"+order.ID) //* redirecting client to this UrlPath once order is placed successfully
}

func(c *Controller) ServeCustomer(ctx *gin.Context) {

	orderID := ctx.Param("id") // fetching id from url param
	if orderID == "" {
		slog.Error("order id is missing!.")
		ctx.String(http.StatusBadGateway,"order id is missing")
		return
	}
	 
	retrievedOrder,err := c.orderModel.GetOrder(orderID)
	if err != nil {
		slog.Error("wrong order id!.")
		ctx.String(http.StatusNotFound,"Order not found,please pass correct order id")
		return
	}

	// otherwise if successfully retrieved order --> render it in html template
	ctx.HTML(http.StatusOK,"customer.tmpl",gin.H{
	// todo --> create templalte "customer.tmpl" to send res
		// * sending this order to tmpl  
		"Order" : retrievedOrder,
		"Statuses": models.OrderStatus,
	})
}
