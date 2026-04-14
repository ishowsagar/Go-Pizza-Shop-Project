package main

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	router.GET("/",c.ServeNewOrderRequest)
	router.GET("/customer/:id",c.ServeCustomer)
	router.POST("/new-order",c.HandleNewOrderPost)


	// serving static
	router.Static("/static","templates/static")
}