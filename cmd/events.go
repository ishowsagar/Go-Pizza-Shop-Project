package main

import (
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (c *Controller) notificationHandler(ctx *gin.Context) {
	orderID := ctx.Query("orderId")

	if orderID == "" {
		ctx.String(400,"Invalid OrderId")
		return
	}
	_,err := c.orderModel.GetOrder(orderID)
	if err != nil {
		ctx.String(400,"Order not found")
		return
	}

	key := NotificationKeyOrder(orderID)
	client := make(chan string,10)
	c.NotificationManager.AddClient(key,client)

	defer func ()  {
		c.NotificationManager.RemoveClient(key,client)
		slog.Info("Customer client disconnected","orderId",orderID)
	}()

	c.StreamSSE(ctx,client)
}

// func that do the serverSideEvent
func(c *Controller) StreamSSE(ctx *gin.Context,client chan string) {
	ctx.Header("Content-type","text/event-stream")
	ctx.Header("Cache-Control","no-cache")
	ctx.Header("Connection","keep-alive")

	ctx.Stream(func(w io.Writer)bool {
		msg,ok := <-client //* reading from channel 
		if ok {
			ctx.SSEvent("message",msg)
			return true
		}
		return false
	})
}

// notifies the admin about the events on sse
func (c *Controller) AdminNotiFHandler(ctx *gin.Context) {
	key := NotificationKeyAdminNewOrders
	client := make(chan string,10)
	c.NotificationManager.AddClient(key,client)
	
	defer func ()  {
		c.NotificationManager.RemoveClient(key,client)
		slog.Info("Admin client disconnected")
	}()

	c.StreamSSE(ctx,client)
}