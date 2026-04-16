package main

import "github.com/ishowsagar/go-pizza-shop/internal/models"

// @types declaration
//  type that stores OrderModel
type Controller struct {
	OrderStore models.OrderModelStore
	UserStore models.UserModelStore
	NotificationManagerStore NotificationManagerStore
}

// func that returns the instance of type struct Controller which stores OrderModel
func NewController(dbModel *models.DBModel) *Controller {
	return &Controller{
		OrderStore: dbModel.OrderStore, //! since that type stores this model and once that is instansiated via func that returns DBModel
		//!we could pass that to this func to enable Controller having OrderModel having the db suplied to it  
		UserStore: dbModel.UserStore,
		NotificationManagerStore: NewNotificationManagerStore(NewNotiFManager()),
	}
}
