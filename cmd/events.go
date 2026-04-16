package main

import (
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// stores only those methods which belongs to Controller but related to event only
type EventsOnlyControllerStore interface {
	handleEventNotification(ctx *gin.Context)
	handleEventStreamSSE(ctx *gin.Context,client chan string)
	hanldeEventAdminNotiFHandler(ctx *gin.Context) 
}


type eventController struct {
	EventsOnlyControllerStoreIface EventsOnlyControllerStore
}

// Routes only event controller methods
func NewEventController(c Controller) eventController {
	return eventController{
		EventsOnlyControllerStoreIface:&c ,
	}
}


func (c *Controller) handleEventNotification(ctx *gin.Context) {
	orderID := ctx.Query("orderId")

	if orderID == "" {
		ctx.String(400,"Invalid OrderId")
		return
	}
	_,err := c.OrderStore.GetOrder(orderID)
	if err != nil {
		ctx.String(400,"Order not found")
		return
	}

	key := NotificationKeyOrder(orderID)
	client := make(chan string,10)
	c.NotificationManagerStore.NotiManagerIface.AddClient(key,client)

	defer func ()  {
		c.NotificationManagerStore.NotiManagerIface.RemoveClient(key,client)
		slog.Info("Customer client disconnected","orderId",orderID)
	}()

	c.handleEventStreamSSE(ctx,client)
}

// func that do the serverSideEvent
func(c *Controller) handleEventStreamSSE(ctx *gin.Context,client chan string) {
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
func (c *Controller) hanldeEventAdminNotiFHandler(ctx *gin.Context) {
	key := NotificationKeyAdminNewOrders
	client := make(chan string,10)
	c.NotificationManagerStore.NotiManagerIface.AddClient(key,client)
	
	defer func ()  {
		c.NotificationManagerStore.NotiManagerIface.RemoveClient(key,client)
		slog.Info("Admin client disconnected")
	}()

	c.handleEventStreamSSE(ctx,client)
}