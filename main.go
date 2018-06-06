package main
import (
	"net/http"
	"github.com/labstack/echo"
	"tabula-rasa/services/account_service/model"
	"github.com/dchest/authcookie"
	"fmt"
	"flag"
	"tabula-rasa/resource"
	"tabula-rasa/libs/memcache"
	"github.com/labstack/echo/middleware"
	"log"
	"tabula-rasa/db"
	"strings"
)

var (
	sslcert   = flag.String("sslcert", "", "path to an SSL cert file")
	sslkey    = flag.String("sslkey", "", "path to an SSL key file")
	port      = flag.Int("port", 8080, "port to be sued by the server")
	host      = flag.String("host", "localhost:8080", "host used to reach this server")
)

// setting up the router with endpoints
func main()  {

	go memcache.MemCache()
	go db.DBConnections()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/a/index")
	})
	e.GET("/login", func(c echo.Context) error {
		return c.File("src/tabula-rasa/web/public/login.html")
	})
	e.POST("/signin", resource.AccountResource{}.Login)
	e.POST("signup", resource.AccountResource{}.Register)

	a := e.Group("/a")
	a.GET("/signout", resource.AccountResource{}.Logout)
	a.Use(authMiddleware)
	a.Static("/static", "src/tabula-rasa/web")
	a.GET("/index", func(c echo.Context) error {
		return c.File("src/tabula-rasa/web/public/index.html")
	})
	// TODO add auth
	api := e.Group("/rest")
	api.GET("/:name/:typ", resource.TableResource{}.Get)
	api.POST("/:name/:typ", resource.TableResource{}.Save)
	api.DELETE("/:name/:typ/:id", resource.TableResource{}.Delete)

	if *sslcert != "" && *sslkey != "" {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%v", *port),
			*sslcert, *sslkey))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", *port)))
	}
}

// check if user is logged in
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(model.CookieName)
		if err == nil {
			login := authcookie.Login(cookie.Value, []byte(model.SECRET))
			if login == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			c.Request().Header.Set(model.HEADER_AUTH_USER_ID, login)
			username := strings.Split(login, "@")[0]
			dbName := username
			err := db.ConnectDB(dbName)
			if err != nil {
				log.Fatalf("cannot initialize db: %v", err)
				return err
			}
			return next(c)
		}
		log.Println(err)
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}
