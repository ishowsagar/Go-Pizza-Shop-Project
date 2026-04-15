package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
) 

func SetupRoutes(router *gin.Engine, c *Controller,store sessions.Store) {

	//@ tuning up mw to be in action⚡
	router.Use(sessions.Sessions("pizza-tracker",store)) // mw fnc that attaches this cookie to the client req --> that could be fetched later

	router.GET("/",c.ServeNewOrderRequest)
	router.GET("/customer/:id",c.ServeCustomer)
	router.POST("/new-order",c.HandleNewOrderPost)
	router.GET("/notifications",c.notificationHandler)

	router.GET("/login",c.HandleLoginGetRequest)
	router.POST("/login",c.HandleLoginPostRequest)
	router.POST("/logout",c.HandleLogout)

	// @ admin seperate routing group
	admin := router.Group("/admin")
	admin.Use(c.AuthorizingMw())
	{
		// ! if it has session with userID attached to it --> makes an get req and call this method on it to serve the response
		admin.GET("",c.ServeAdminDashboardPanel)
		admin.GET("/notifications",c.AdminNotiFHandler)
		admin.POST("/order/:id/update",c.HandleOrderUpdateRequest)
		admin.POST("/order/:id/delete",c.HandleOrderDelelteRequest)
	}
	// serving static
	router.Static("/static","templates/static")
}