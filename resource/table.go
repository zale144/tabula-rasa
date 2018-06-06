package resource

import (
	"net/http"
	"tabula-rasa/storage"
	"io/ioutil"
	"github.com/labstack/echo"
	"tabula-rasa/services/account_service/model"
	"fmt"
	"strings"
	"log"
	"github.com/dchest/authcookie"
)

// table resource
type TableResource struct{}

// method for retrieving resources
func (tr TableResource) Get(c echo.Context) error {
	name := c.Param("name")
	spec := c.Param("typ")

	// get the cookie value
	cookieVal := getCookieValue(&c)
	// get the database name from the cookie value
	dbName, err := getDbName(cookieVal)
	if err != nil {
		log.Println(err.Error())
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	entity, err := storage.TableStorage{}.Get(name, "", spec, dbName)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, entity)
}

// method for saving resources
func (tr TableResource) Save(c echo.Context) error {
	typ := c.Param("typ")
	name := c.Param("name")
	if name == "" {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, "name is mandatory"))
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	// get the cookie value
	cookieVal := getCookieValue(&c)
	// get the database name from the cookie value
	dbName, err := getDbName(cookieVal)
	if err != nil {
		log.Println(err.Error())
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	defer c.Request().Body.Close()
	out, err := storage.TableStorage{}.Save(name, typ,  dbName, body)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, out)
}

// method for deleting resources
func (tr TableResource) Delete(c echo.Context) error {
	typ := c.Param("typ")
	id := c.Param("id")
	if id == "" {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, "id is mandatory"))
	}
	name := c.Param("name")
	if name == "" {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, "name is mandatory"))
	}
	// get the cookie value
	cookieVal := getCookieValue(&c)
	// get the database name from the cookie value
	dbName, err := getDbName(cookieVal)
	if err != nil {
		log.Println(err.Error())
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	err = storage.TableStorage{}.Delete(name, id, typ, dbName)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
// get the username/database name from the cookie value
func getDbName(value string) (string, error) {
	email := authcookie.Login(value, []byte(model.SECRET))
	if email == "" {
		err := fmt.Errorf("no user authenticated")
		log.Println(err.Error())
		return "", err
	}
	username := strings.Split(email, "@")[0]
	return username, nil
}
// get the value from the cookie
func getCookieValue(cp *echo.Context) string {
	c := *cp
	headers := c.Request().Header
	cookieStr := headers.Get("cookie")
	if cookieStr == "" {
		err := fmt.Errorf("empty cookie")
		log.Println(err.Error())
	}
	value := strings.Replace(cookieStr, model.CookieName+"=", "", -1)
	fmt.Println(value)
	return value
}