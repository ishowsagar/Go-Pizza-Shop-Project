package models

import (
	"time"

	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

// @ Declaring Var which hold these vals
var OrderStatus = []string{"Order Placed","Preparing","Baking","Quality Check","Ready"}
var PizzaTypes = []string {"Margherita", "Pepperoni", "Vegetarian", "Hawaiian", "Bbq Chicken", "Meat Lovers", "Buffalo Chicken", "Supreme", "Truffle Mushroom", "Four Cheese"}
var PizzaSizes = []string{"Small","Medium","Large","X-Large"}

// @ types declration

// holds db connection of type *gorm.db
type OrderModel struct {
	DB *gorm.DB
} 

// holds data for data type that holds ordered item
type OrderItem struct {
	ID string `gorm:"primaryKey;size:14" json:"id"` //! gorm field tags to convert to sql table
	OrderID string `gorm:"index;size:14;not null" json:"orderId"`
	Size string `gorm:"not null" json:"size"`
	Pizza string `gorm:"not null" json:"pizza"`
	Instructions string `json:"instructions"`
}

// holds order type of data struct~ure
type Order struct {
	ID string `gorm:"primaryKey;size:14" json:"id"` //! gorm field tags to convert to sql table
	Status string `gorm:"not null" json:"status"` 
	CustomerName string `gorm:"not null" json:"customerName"` 
	Phone string `gorm:"not null" json:"phone"` 
	Address string `gorm:"not null" json:"address"`
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"pizzas"` //! refs to order id from another table
	CreatedAt time.Time `json:"createdAt"` 
}

// @ interface --> stores all methods that belongs and implemented by the type --> *Order

// @methods that belongs to the type -> *Order
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	
	// checking if order has id or not
	if o.ID == "" {
		o.ID = shortid.MustGenerate() // assigns a id to it
	}
	
	return nil 
}


// @methods that belongs to the type -> *OrderItem
func(oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	
	// assigning id to it if it did not have
	if oi.ID == "" {
		oi.ID = shortid.MustGenerate()
	}
	return nil
}


// @methods that belongs to the type -> *OrderModel
func (om *OrderModel) CreateOrder(order *Order) error {
	return om.DB.Create(order).Error // returns error if there was creating an order
}

//  func that retrieves order from provided id to the method
func (om *OrderModel) GetOrder(orderID string) (*Order,error) {
	
	var order Order // store data type of Order struct --> preloading into this var addr 👇👇
	err := om.DB.Preload("Items").First(&order,"id = ?",orderID).Error // loading items table --> return one row only --> find entry where id = ? this, provide placeholder val
	return &order,err// returns error if there was creating an order

}

// method that returns all orders
func (o *OrderModel) GetAllOrders() ([]Order,error) {
	var orders []Order
	tx := o.DB.Preload("Items").Order("created_at desc").Find(&orders) //* finding all entries and ordering by these cols and populating into &orders Slice
	if tx.Error != nil {
		return nil,tx.Error
	}
	return orders,nil
}

// update order by id
func(o *OrderModel) UpdateOrderStatus(id string,status string) error {
	err := o.DB.Model(&Order{}).Where("id=?",id).UpdateColumn("status",status).Error
	if err!= nil {
		return err
	}
	return nil
}

//  delete order by id
func(o *OrderModel) DeleteOrder(id string) error {
	err := o.DB.Select("Items").Delete(&Order{ID: id}).Error
	if err!= nil {
		return err
	}
	return nil
}





