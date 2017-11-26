package service

import (
	"tabula-rasa/dao"
)

func Get(args ...string) (string, error) {
	return dao.Get(args...)
}

func Delete(name, id string) (string, error) {
	return dao.Delete(name, id)
}

func Create(name string, obj []byte) (string, error) {
	return dao.Create(name, obj)
}