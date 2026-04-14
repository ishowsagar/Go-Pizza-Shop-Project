package main

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @types declaration for utilities needed throught the server
type Config struct {
	Port string
	DBConnStr string
}

// fnc that returns instance of type Config which stores port & dbconnstr
func LoadConfig() (*Config,error) {

	// loadinv .env file to use protected vars
	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load .env file","error",err)
		return nil,err
	}

	// loading var
	PORT := os.Getenv("PORT")
	DbConnStr := os.Getenv("DB_CONNECTION_STRING")
	config :=Config{
		Port: PORT,
		DBConnStr: DbConnStr,
	}

	return &config,nil
}


// loads template
func LoadTemplates(router *gin.Engine) error {
	functions := template.FuncMap{
		"add" : func(a,b int) int {return a+b},
		"json": func(v any) (template.JS, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(b), nil
		},
		} 
	tmpl,err := template.New("").Funcs(functions).ParseGlob("templates/*/*.tmpl")
	if err != nil {
		return err
	}
	router.SetHTMLTemplate(tmpl)
	return nil
}