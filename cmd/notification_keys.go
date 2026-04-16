package main

const NotificationKeyAdminNewOrders = "admin:new-orders"

func NotificationKeyOrder(orderID string) string {
	return "order:" + orderID
}
