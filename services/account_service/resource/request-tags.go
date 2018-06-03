package resource

type Account struct {
	FirstName      string  `json:"firstName" form:"firstName" query:"firstName"`
	LastName       string  `json:"lastName" form:"lastName" query:"lastName"`
	Email          string  `json:"email" form:"email" query:"email"`
	Password       string  `json:"password" form:"password" query:"password"`
	Active		   bool    `json:"active" form:"active" query:"active"`
}