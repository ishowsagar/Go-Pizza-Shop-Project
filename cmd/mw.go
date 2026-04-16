package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// stores only those methods which belongs to Controller but related to middleware only
type MiddlewareOnlyControllerStore interface {
	AuthorizingMw() gin.HandlerFunc
}


type middlewareController struct {
	MiddlewareOnlyControllerStoreIface MiddlewareOnlyControllerStore
}

// Routes only middleware controller methods
func NewMiddlewareController(c Controller) middlewareController {
	return middlewareController{
		MiddlewareOnlyControllerStoreIface:&c ,
	}
}




//$ method belongs to *Controller type
func(c *Controller) AuthorizingMw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := GetSessionString(ctx,"userID")
		
		// check if userID is not available --> send client to the login page
		if userID == "" {
			ctx.Redirect(http.StatusSeeOther,"/login")
			ctx.Abort()
			return
		}

		// if id is available --> get that user ( for user rltd calls --> have a sep usrMdl)
		_,err := c.UserStore.GetUserByID(userID) //! fetching user from retrieved userID from session when it was set during Post login req
		if err == gorm.ErrRecordNotFound {
			ClearSession(ctx)
			ctx.Redirect(http.StatusBadRequest,"/login")
			ctx.Abort()
			return 
		}
		ctx.Next() // call the next chained method which was expected to be loaded
	}
}
