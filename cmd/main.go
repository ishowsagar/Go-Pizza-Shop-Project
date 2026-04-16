package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ishowsagar/go-pizza-shop/internal/models"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		slog.Error("failed to load config","error",err)
		return
	}

	// initializing a logger of type slog.log
	logger := slog.New(slog.NewTextHandler(os.Stdout,nil))
	slog.SetDefault(logger) // default logger for whole app 

	// !setting up an instance of DBModel :- UserModelStore,OrderModelStore --> Db's
	dbModel,err := models.ConnectionToDB(config.DBConnStr)
	
	if err != nil {
		// fmt.Errorf("failed to load config -%v",err)
		slog.Error("failed to load DbConnection","error",err)
		os.Exit(1) // exiting program
	}


	slog.Info("Db connection is successfully loaded ⚡...")
	RegisterCustomValidator() // start up validator

	// ! setting up session store
	sessionStore := SetupStore(dbModel.DbConn,[]byte(config.SessionSecretKey)) //* NOw u can see we can't directly expose db of type *gorm.db but interface is acting s security
	//!   setting up instance of controller
	controller := NewController(dbModel)

	// * instead of passing controller --> we can pass instance of type that stores iface where all meths belongs to the controller
	masterRoutingController := NewMasterController(*controller)

	// @ creating router 
	router := gin.Default()
	err = LoadTemplates(router)
	if err != nil {
		slog.Error("failed to load template","error",err)
		os.Exit(1) // exiting program
	}
	
	SetupRoutes(router,&masterRoutingController,sessionStore)
	slog.Info("server has started successfully🚀...","url","http://localhost:"+config.Port)
	router.Run(":"+config.Port)
}