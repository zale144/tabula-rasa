package resource

import (
	"github.com/labstack/echo"
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"tabula-rasa/db"
	"tabula-rasa/services/account_service/model"
	"time"
	"fmt"
)

type Account struct {
	FirstName      string  `json:"firstName" form:"firstName" query:"firstName"`
	LastName       string  `json:"lastName" form:"lastName" query:"lastName"`
	Email          string  `json:"email" form:"email" query:"email"`
	Password       string  `json:"password" form:"password" query:"password"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" query:"confirmPassword"`
}

type AccountResource struct {}

// method for handling login requests
func (ar AccountResource) Login(c echo.Context) error {
	acc := new(Account)                 //initialize  struct Account
	if err := c.Bind(acc); err != nil { //get and bind data from request to struct Account
		err := fmt.Errorf("Invalid JSON payload")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	jsonStr, err := json.Marshal(acc)
	if err != nil {
		log.Println(err)
		return err
	}
	url := "http://localhost:8081/admin/accounts/signin" // TODO get url from args
	// send the login request to the account_service server
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return err
	}
	// get the session cookie from the response
	var cookie http.Cookie
	err = json.NewDecoder(resp.Body).Decode(&cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	c.SetCookie(&cookie)
	// get the username from the cookie, and use it as the database name
	dbName, err := getDbName(cookie.Value)
	if err != nil {
		log.Println(err)
		return err
	}
	// connect to the database with provided name
	db.ConnectDB(dbName)
	return c.Redirect(http.StatusSeeOther, "/")
}
// method for handling registration requests
func (ar AccountResource) Register(c echo.Context) error {
	acc := new(Account)                 //initialize  struct Account
	if err := c.Bind(acc); err != nil { //get and bind data from request to struct Account
		err := fmt.Errorf("Invalid JSON payload")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	jsonStr, err := json.Marshal(acc)
	if err != nil {
		log.Println(err)
		return err
	}
	url := "http://localhost:8081/admin/accounts" // TODO get url from args
	// send the registration request to the account_service server
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return err
	}
	// get the cookie from the response
	var cookie http.Cookie
	body := json.NewDecoder(resp.Body)
	err = body.Decode(&cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(cookie) // TODO properly handle error responses
	c.SetCookie(&cookie)
	// use the username from the cookie as the name for the database
	dbName, err := getDbName(cookie.Value)
	if err != nil {
		log.Println(err)
		return err
	}
	// create a database with the provided name
	err = db.CreateDatabase(dbName)
	if err != nil {
		log.Println(err)
		return err
	}
	// connect to the newly created database
	db.ConnectDB(dbName)
	return c.Redirect(http.StatusSeeOther, "/")
}
// method for handling logout requests
func (ar AccountResource) Logout(c echo.Context) error {
	// get the database name from the username in the cookie
	cookieVal := getCookieValue(&c)
	dbName, err := getDbName(cookieVal)
	if err != nil {
		if err != nil {
			log.Println(err)
			return err
		}
	}
	// expire the cookie
	cookie := &http.Cookie{
		Name:    model.CookieName,
		Expires: time.Now(),
		Path:    "/",
	}
	c.SetCookie(cookie)
	// close the database connection
	db.Disconnect(dbName)
	return nil
}
