package models

import (
	"time"
)

type User struct{
	ID int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevel int
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Question struct{
	ID int
	Subject string
	Category string
	Description string
	UserId int
	User User
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Reply struct{
	ID int
	QuestionId int
	ReplyId int
	Description string
	UserId int
	User User
	CreatedAt time.Time
	UpdatedAt time.Time
}