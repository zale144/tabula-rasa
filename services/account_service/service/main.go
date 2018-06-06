package main

import (
	"github.com/labstack/echo"
	"tabula-rasa/services/account_service/model"
	"tabula-rasa/services/account_service/resource"
	"log"
	"fmt"
	"flag"
	"github.com/labstack/echo/middleware"
	"tabula-rasa/services/account_service/db"
)

var (
	sslcert   = flag.String("sslcert", "", "path to an SSL cert file")
	sslkey    = flag.String("sslkey", "", "path to an SSL key file")
	port      = flag.Int("port", 8081, "port to be sued by the server")
	host      = flag.String("host", "localhost:8080", "host used to reach this server")
)

// setting up the router with endpoints
func main()  {

	err := dba.GetDBConnection()
	if err != nil {
		log.Fatalf("cannot initialize db: %v", err)
		return
	}

	if *sslcert != "" && *sslkey != "" {
		model.AppURL = fmt.Sprintf("https://%v", *host)
	} else {
		model.AppURL = fmt.Sprintf("http://%v", *host)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	admin := e.Group("/admin")

	admin.GET("/accounts/:email", resource.AccountResource{}.GetByEmail)
	admin.POST("/accounts/signin", resource.AccountResource{}.Signin)
	admin.POST("/accounts", resource.AccountResource{}.Add)
	admin.PUT("/accounts/:id", resource.AccountResource{}.Update)

	if *sslcert != "" && *sslkey != "" {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%v", *port),
			*sslcert, *sslkey))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", *port)))
	}
}
