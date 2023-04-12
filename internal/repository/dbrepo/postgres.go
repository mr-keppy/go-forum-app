package dbrepo

import (
	"context"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"

	"github.com/mr-keppy/go-forum/internal/models"
)

func (m *postgreDBRepo) AllUsers() bool {
	return true
}

func (m *postgreDBRepo) GetUserByID(id int) (models.User, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	select id, first_name, last_name, email, password, access_level from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err:= row.Scan(
		&u.ID, u.FirstName, u.LastName, u.Email, u.Password, u.AccessLevel,
	)
	if err != nil{
		return u, err
	}else{
		return u, nil
	}
}

//update user 
func (m *postgreDBRepo) UpdateUser(u models.User) ( error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	update users set  first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5 where id = $6`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil{
		return err
	}else{
		return nil
	}
}

func (m *postgreDBRepo) Authenticate(email, password string) (int, string, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var id int
	var hashedPassword string

	row:= m.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err:= row.Scan(&id, &hashedPassword)

	if err!=nil{
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword{
		return 0, "", errors.New("incorrect username/password")
	}else if err!=nil{
		return 0, "",err
	}

	return id, hashedPassword, nil
}