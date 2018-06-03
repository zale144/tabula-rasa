package service

import (
	"net/http"
	"github.com/dchest/authcookie"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"tabula-rasa/services/account_service/model"
	"tabula-rasa/services/account_service/resource"
	"tabula-rasa/services/account_service/db"
	"log"
	"fmt"
	"flag"
)

var (
	sslcert   = flag.String("sslcert", "", "path to an SSL cert file")
	sslkey    = flag.String("sslkey", "", "path to an SSL key file")
	port      = flag.Int("port", 8081, "port to be sued by the server")
	host      = flag.String("host", "localhost:8080", "host used to reach this server")
)

// setting up the router with endpoints
func main()  {

	err := db.GetDBConnection()
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

	admin := e.Group("/a")
	admin.Use(authMiddleware)

	apiv1 := e.Group("/api/v1")

	apiv1.POST("/accounts/signin/", resource.AccountResource{}.Signin)
	apiv1.POST("/accounts", resource.AccountResource{}.Add)
	apiv1.GET("/accounts/:id/signout", resource.AccountResource{}.Signout)
	apiv1.PUT("/accounts/:id", resource.AccountResource{}.Update)

	if *sslcert != "" && *sslkey != "" {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%v", *port),
			*sslcert, *sslkey))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", *port)))
	}

}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(model.CookieName)
		if err == nil {
			login := authcookie.Login(cookie.Value, []byte(model.SECRET))
			if login == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/")
			}
			c.Request().Header.Set(model.HEADER_AUTH_USER_ID, login)
			return next(c)
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
