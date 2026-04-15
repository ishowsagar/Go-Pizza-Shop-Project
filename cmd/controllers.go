package main

import "github.com/ishowsagar/go-pizza-shop/internal/models"

// @types declaration
//  type that stores OrderModel
type Controller struct {
	orderModel *models.OrderModel
	UserModel *models.UserModel
	NotificationManager *NotificationManager
}

// func that returns the instance of type struct Controller which stores OrderModel
func NewController(dbModel *models.DBModel) *Controller {
	return &Controller{
		orderModel: dbModel.OrderModel, //! since that type stores this model and once that is instansiated via func that returns DBModel
		//!we could pass that to this func to enable Controller having OrderModel having the db suplied to it  
		UserModel: dbModel.UserModel,
		NotificationManager: NewNotiFManager(),
	}
}
