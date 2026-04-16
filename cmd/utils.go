package main

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"os"

	"github.com/gin-contrib/sessions"
	gormSessions "github.com/gin-contrib/sessions/gorm"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @types declaration for utilities needed throught the server
type Config struct {
	Port string
	DBConnStr string
	SessionSecretKey string
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
	port := os.Getenv("PORT")
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	sessionSecretKey := os.Getenv("Session_Secret_Key")
	config :=Config{
		Port: port,
		DBConnStr: dbConnStr,
		SessionSecretKey: sessionSecretKey,
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

//func that sets up session's store
func SetupStore(db *gorm.DB,secretKey []byte) sessions.Store{
	store := gormSessions.NewStore(db,true,secretKey)
	store.Options(sessions.Options{
		Path: "/", // cookie available for all paths
		MaxAge: 86400,
		HttpOnly: true,
		Secure: true,
		SameSite: 4,
	})
	return store
}

// func that sets up key-valeu on the req's session
func SetSessionValue(ctx *gin.Context,key string,value interface{}) error {
	session := sessions.Default(ctx) //* getting default session
	session.Set(key,value)
	return session.Save() //* saves this context on the req
}

// func that returns the stored sessions string in the request
func GetSessionString(ctx *gin.Context,key string) string {
	session := sessions.Default(ctx)
	val:= session.Get(key)

	// ! if key has nil value --> zero value 
	if val == nil {
		return ""
	}

	str,_ := val.(string)
	return str
}



// func that clears/destroys the stored sessions string in the request
func ClearSession(ctx *gin.Context) error {
	session := sessions.Default(ctx)
	session.Clear() // deletes all key-val in the session
	return session.Save()
}




