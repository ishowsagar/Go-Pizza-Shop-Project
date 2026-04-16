package main

// @ Base for all types that have been storing Interfaces
type MasterController struct {
	adminStore adminController //& stores interface which has all methods for admin,like as all belongs to Controller too
	customerStore customerController
	eventStore eventController
	middlewareStore middlewareController
}

//! returns instance of type that stores all interfaces for all corresponding needy methods 
// @ Since all  type that stores ifaces need controller type would be feeded by this func to enable all stores
func NewMasterController(c Controller) MasterController {
	return MasterController{
		adminStore: NewAdminController(c),
		customerStore : NewCustomerController(c),
		eventStore : NewEventController(c),
		middlewareStore : NewMiddlewareController(c),

	}
}