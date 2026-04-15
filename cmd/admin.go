package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishowsagar/go-pizza-shop/internal/models"
)

// type that stores loginData
type LoginData struct {
	Error string
}

type AdminDashboardData struct {
	Username string
	Orders []models.Order
	Statuses []string
}

// method that belongs to the *Controller
func(c *Controller) HandleLoginGetRequest(ctx *gin.Context) {
	//! loading template of loginPage in handler/controller methods

	ctx.HTML(http.StatusOK,"login.tmpl",LoginData{})		
}

// method that makes post req to send data of loggin in client
func(c *Controller) HandleLoginPostRequest(ctx *gin.Context) {
	var userloginPayload struct {
		Username string `form:"username" binding:"required,min=4,max=27"`
		Password string `form:"password" binding:"required,min=4"`
	} 

	err := ctx.ShouldBind(&userloginPayload) // validates & populates incoming payoload like r.Body into this var
	if err!= nil {
		ctx.HTML(http.StatusOK,"login.tmpl",LoginData{
			Error: "Invalid input. Username must be 4-27 characters and password must be at least 4 characters.", 
		})
		return
	}

	// if there was no validator errors when did the binding process thing
	user,err := c.UserModel.AuthenticateUser(userloginPayload.Username,userloginPayload.Password)
	if err!= nil {
		ctx.HTML(http.StatusOK,"login.tmpl",LoginData{
			Error: "Invalid username or password.", 
		})
		return
	}

	// * if we get user is authenticated -> then we set session values onto the req and redirect to "/admin" panel
	SetSessionValue(ctx,"userID",fmt.Sprintf("%v",user.ID)) //! sets session with provided key-val pairs on the req
	SetSessionValue(ctx,"username",user.Username)
	ctx.Redirect(http.StatusSeeOther,"/admin")
}

func(c *Controller) HandleLogout(ctx *gin.Context) {
	err := ClearSession(ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError,"error",err.Error())
		return
	}

	//! if session is cleared then -> redirect client to the login page
	ctx.Redirect(http.StatusSeeOther,"/login") //@ Servind admin panel template
}

func(c *Controller) ServeAdminDashboardPanel(ctx *gin.Context) {

	// * fetching all orders in this panel for admin login
	orders,err := c.orderModel.GetAllOrders()
	if err != nil {
		ctx.String(http.StatusInternalServerError,err.Error())
		return
	}

	// get sesseion 	
	retrievedUsername := GetSessionString(ctx,"username") //* would be a method that fetches info from req by using the passed string
	ctx.HTML(http.StatusOK,"admin.tmpl",AdminDashboardData{
		// @ Data sending to the above mentioned "admin.tmpl" template --> "." refer to this all data as whole
		Username:retrievedUsername,
		Orders: orders,
		Statuses: models.OrderStatus,
	})	
}

// update order status
func(c *Controller) HandleOrderUpdateRequest(ctx *gin.Context) {
	orderID := ctx.Param("id")
	newStatus := ctx.PostForm("status") //* fetching new status val that has to be updated with

	err := c.orderModel.UpdateOrderStatus(orderID,newStatus)
	if err != nil {
		ctx.String(http.StatusInternalServerError,err.Error())
		return
	}

	c.NotificationManager.Notify(NotificationKeyOrder(orderID),"order_updated")

	ctx.Redirect(http.StatusSeeOther,"/admin")
}

// deletes order from the db
func(c *Controller) HandleOrderDelelteRequest(ctx *gin.Context) {
	// fetch id of the order that is going to be deleted
	orderID := ctx.Param("id")
	if orderID == "" {
		ctx.String(http.StatusInternalServerError,"Order id not found!.")
	}
	// call order delete meth that is attached on O.M
	err := c.orderModel.DeleteOrder(orderID)
	if err != nil {
		ctx.String(http.StatusInternalServerError,err.Error())
	}

	ctx.Redirect(http.StatusSeeOther,"/admin")
}