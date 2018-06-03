package storage

import (
	"tabula-rasa/services/account_service/model"
	"strings"
	"tabula-rasa/services/account_service/libs/crypto"
)

// AccountStorage gives all methods to manage Accounts
type AccountStorage struct {}


// GetOne Account
func (ac AccountStorage) GetOne(id uint) (*model.Account, error) {

	var account model.Account

	return &account, nil

}

// Insert an Account
func (ac AccountStorage) Insert(Account model.Account) error {
	Account.Email = strings.ToLower(Account.Email)
	Account.Password = crypto.CryptPrivate(Account.Password, crypto.CRYPT_SETTING)


	return nil
}

// Get Account by email
func (ac AccountStorage) GetByEmail(email string) (*model.Account, error) {

	var account model.Account

	return &account, nil

}

type AccountUpdater struct {
	accountID uint
	updaterID uint

	updates map[string]interface{}
}

// Update an Account
func (ac AccountStorage) NewAccountUpdater(accountID, updaterID uint) *AccountUpdater {
	return &AccountUpdater{
		accountID: accountID,
		updaterID: updaterID,
		updates:   make(map[string]interface{}),
	}
}
func (a *AccountUpdater) FirstName(f string) *AccountUpdater {
	a.updates["first_name"] = f
	return a
}

func (a *AccountUpdater) LastName(f string) *AccountUpdater {
	a.updates["last_name"] = f
	return a
}

func (a *AccountUpdater) Email(f string) *AccountUpdater {
	a.updates["email"] = strings.ToLower(f)
	return a
}

func (a *AccountUpdater) Active(f bool) *AccountUpdater {
	a.updates["active"] = f
	return a
}

func (a *AccountUpdater) Password(f string) *AccountUpdater {
	a.updates["password"] = crypto.CryptPrivate(f, crypto.CRYPT_SETTING)
	return a
}

func (a *AccountUpdater) Update(/*tx *gorm.DB*/) error {
	/*if tx == nil {
		tx = model.PgsqlDB
	}

	tx = tx.Model(&model.Account{Model: gorm.Model{ID: a.accountID}}).
		Updates(a.updates)
	rowsAffected, err := tx.RowsAffected, tx.Error
	if err != nil {
		if model.IsUniqueConstraintError(err, model.UniqueConstraintEmail) {
			return &model.EmailDuplicateError{}
		}
		return err
	}
	if rowsAffected == 0 {
		err = fmt.Errorf("record not found")
		return err
	}*/
	return nil

}