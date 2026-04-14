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

	// !setting up an instance of DBModel
	dbModel,err := models.ConnectionToDB(config.DBConnStr)
	
	if err != nil {
		// fmt.Errorf("failed to load config -%v",err)
		slog.Error("failed to load DbConnection","error",err)
		os.Exit(1) // exiting program
	}

	slog.Info("Db connection is successfully loaded ⚡...")
	RegisterCustomValidator() // start up validator

	//!   setting up instance of controller
	controller := NewController(dbModel)

	// @ creating router 
	router := gin.Default()
	err = LoadTemplates(router)
	if err != nil {
		slog.Error("failed to load template","error",err)
		os.Exit(1) // exiting program
	}
	
	SetupRoutes(router,controller)
	slog.Info("server has started successfully🚀...","url","http://localhost:"+config.Port)
	router.Run(":"+config.Port)
}