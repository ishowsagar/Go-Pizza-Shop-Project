package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// func takes in routerController interface to serve controllers methods
func SetupRoutes(router *gin.Engine, masterController *MasterController,store sessions.Store) {

	//@ tuning up mw to be in action⚡
	router.Use(sessions.Sessions("pizza-tracker",store)) // mw fnc that attaches this cookie to the client req --> that could be fetched later

	router.GET("/", masterController.customerStore.customerOnlysRouteControllerStore.ServeCustomerNewOrderRequest)
	router.GET("/customer/:id",masterController.customerStore.customerOnlysRouteControllerStore.ServeCustomer)
	router.POST("/new-order",masterController.customerStore.customerOnlysRouteControllerStore.HandleCustomerNewOrderPost)
	router.GET("/notifications",masterController.eventStore.EventsOnlyControllerStoreIface.handleEventNotification)

	router.GET("/login",masterController.adminStore.adminRoutesControllerIface.HandleAdminLoginGetRequest)
	router.POST("/login",masterController.adminStore.adminRoutesControllerIface.HandleAdminLoginPostRequest)
	router.POST("/logout",masterController.adminStore.adminRoutesControllerIface.HandleAdminLogout)

	// @ admin seperate routing group
	admin := router.Group("/admin")
	admin.Use(masterController.middlewareStore.MiddlewareOnlyControllerStoreIface.AuthorizingMw())
	{
		// ! if it has session with userID attached to it --> makes an get req and call this method on it to serve the response
		admin.GET("",masterController.adminStore.adminRoutesControllerIface.ServeAdminDashboardPanel)
		admin.GET("/notifications",masterController.eventStore.EventsOnlyControllerStoreIface.hanldeEventAdminNotiFHandler)
		admin.POST("/order/:id/update",masterController.adminStore.adminRoutesControllerIface.HandleAdminOrderUpdateRequest)
		admin.POST("/order/:id/delete",masterController.adminStore.adminRoutesControllerIface.HandleAdminOrderDelelteRequest)
	}
	// serving static
	router.Static("/static","templates/static")
}