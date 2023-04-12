package repository

import (

	"github.com/mr-keppy/go-forum/internal/models"
)

type DatabaseRepo interface{
	AllUsers() bool
	GetUserByID(id int) (models.User, error)
	Authenticate(email, password string) (int, string, error)
	UpdateUser(u models.User) ( error)
}