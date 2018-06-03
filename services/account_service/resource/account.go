package resource

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"tabula-rasa/services/account_service/model"
	"github.com/dchest/authcookie"
	"time"
	"tabula-rasa/services/account_service/libs/email"
	"tabula-rasa/services/account_service/storage"
	"strings"
	"tabula-rasa/services/account_service/libs/crypto"
	"strconv"
)

// AccountResource
type AccountResource struct{}


// Add new account
func (ac AccountResource) Add(c echo.Context) error {

	type AddAccountRequest struct {
		LastName        string `json:"lastName"`
		FirstName       string `json:"firstName"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	account := new(AddAccountRequest) //initialize  struct
	if err := c.Bind(account); err != nil { //get and bind data from request to struct account
		err = fmt.Errorf("Invalid data: %v", err)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if account.FirstName == "" || account.LastName == "" {
		err := fmt.Errorf("First Name and Last Name are required")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if account.Email == "" {
		err := fmt.Errorf("Email is required")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	if err := email.ValidateFormat(account.Email); err != nil {
		err = fmt.Errorf("Invalid email")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if len(account.Password) < 6 {
		err := fmt.Errorf("Password must be at least six characters long")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if account.Password != account.ConfirmPassword {
		err := fmt.Errorf("Passwords don't match")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	var AccountModel model.Account // initialize model
	//bind Data from struct "account" to model.account
	AccountModel.LastName = account.LastName
	AccountModel.FirstName = account.FirstName
	AccountModel.Email = account.Email
	AccountModel.Password = account.Password

	err := storage.AccountStorage{}.Insert(AccountModel) //Call Insert with model account To add the new record
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	// else here it is OK => correct user
	cookie := &http.Cookie{
		Name:  model.CookieName,
		Value: authcookie.NewSinceNow(AccountModel.Email, 24*time.Hour, []byte(model.SECRET)),
		Path:  "/",
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, "Created")
}

// signin an account
func (ac AccountResource) Signin(c echo.Context) error {
	account := new(Account)                 //initialize  struct Account
	if err := c.Bind(account); err != nil { //get and bind data from request to struct Account
		err := fmt.Errorf("Invalid JSON payload")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	if account.Email == "" {
		err := fmt.Errorf("Email is required")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	if err := email.ValidateFormat(account.Email); err != nil {
		err = fmt.Errorf("Invalid email")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	if account.Password == "" {
		err := fmt.Errorf("Password is required")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	accountModel, err := storage.AccountStorage{}.GetByEmail(strings.ToLower(account.Email))
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if !crypto.PortableHashCheck(account.Password, accountModel.Password) {
		err = fmt.Errorf("wrong password for user %q", account.Email)
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	// else here it is OK => correct user
	cookie := &http.Cookie{
		Name:  model.CookieName,
		Value: authcookie.NewSinceNow(accountModel.Email, 24*time.Hour, []byte(model.SECRET)),
		Path:  "/",
	}

	c.SetCookie(cookie)

	return nil
}

// signout from an account
func (ac AccountResource) Signout(c echo.Context) error {

	cookie := &http.Cookie{
		Name:    model.CookieName,
		Expires: time.Now(),
		Path:    "/",
	}

	c.SetCookie(cookie)
	return nil
}

// update an account
func (ac AccountResource) Update(c echo.Context) error {
	type UpdateAccountRequest struct {
		FirstName       string `json:"firstName" form:"first-name"`
		LastName        string `json:"lastName" form:"last-name"`
		Email           string `json:"email" form:"email"`
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirmPassword" form:"confirm-password"`
		Active          bool   `json:"active" form:"active"`
	}

	updateAccountReq := new(UpdateAccountRequest) //initialize  struct updateAccountReq

	if err := c.Bind(updateAccountReq); err != nil { //get and bind data from request to struct updateAccountReq

		err := fmt.Errorf("Invalid JSON payload")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	AccountID, err := strconv.Atoi(c.Param("id"))
	if err != nil { //get updateAccountReq Id from query
		err := fmt.Errorf("updateAccountReq ID is invalid")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if updateAccountReq.FirstName == "" || updateAccountReq.LastName == "" {
		err := fmt.Errorf("First Name and Last Name are required")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if updateAccountReq.Email == "" {
		err := fmt.Errorf("Email is required")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}

	if err := email.ValidateFormat(updateAccountReq.Email); err != nil {
		err = fmt.Errorf("Invalid email")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	accountUpdater := storage.AccountStorage{}.NewAccountUpdater(uint(AccountID), uint(AccountID)).
		FirstName(updateAccountReq.FirstName).
		LastName(updateAccountReq.LastName).
		Email(updateAccountReq.Email)

	redirectFlag := "Updated"

	// we are not always updating the password, so since the password
	// is not returned to the form when it loads, here we check if user
	// inputted something, and make sure that it's not less than 6 chars long
	if len(updateAccountReq.Password) > 0 {
		if len(updateAccountReq.Password) < 6 {
			err := fmt.Errorf("Password must be at least six characters long")
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
			return err
		}
		// and that password and confirm password match
		if updateAccountReq.Password != updateAccountReq.ConfirmPassword {
			err := fmt.Errorf("Passwords don't match")
			c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
			return err
		}
		accountUpdater = accountUpdater.Password(updateAccountReq.Password)
	}

	// if status is being updated
	if updateAccountReq.Active {
		accountUpdater = accountUpdater.Active(updateAccountReq.Active)
	}

	err = accountUpdater.Update(/*nil*/)
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}

	//Update Echo Cookie
	UpdatedUser, err := storage.AccountStorage{}.GetOne(uint(AccountID))
	if err != nil {
		c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
		return err
	}

	cookie := &http.Cookie{
		Name:  model.CookieName,
		Value: authcookie.NewSinceNow(UpdatedUser.Email, 24*time.Hour, []byte(model.SECRET)),
		Path:  "/",
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, redirectFlag)
}