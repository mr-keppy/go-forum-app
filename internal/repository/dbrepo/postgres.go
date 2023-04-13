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

//create user
func (m *postgreDBRepo) CreateUser(u models.User) ( error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	INSERT INTO public.users
		(first_name, last_name, email, "password", access_level, created_at, updated_at)
		VALUES($1, $2, $3, $4, 1, $5, $6);
		`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.AccessLevel,time.Now(), time.Now())
	if err != nil{
		return err
	}else{
		return nil
	}
}

//get user
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

//authenticate
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

//create question
func (m *postgreDBRepo) CreateQuestion(u models.Question) ( error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	INSERT INTO public.questions
		(subject, categoryId, description, userId, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6);
		`

	_, err := m.DB.ExecContext(ctx, query, u.Subject, u.Category, u.Description, u.UserId, time.Now(), time.Now())
	if err != nil{
		return err
	}else{
		return nil
	}
}

//get question by Id
func (m *postgreDBRepo) GetQuestionByID(id int) (models.Question, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	select id, subject, categoryId, description, userId, created_at, updated_at from questions where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.Question

	err:= row.Scan(
		&u.ID, &u.Subject, &u.Category, &u.Description, &u.UserId, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil{
		return u, err
	}else{
		return u, nil
	}
}

//get all question by UserId
func (m *postgreDBRepo) GetQuestionsByUserID(userId int) ([]models.Question, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var questions []models.Question

	query := `
	select id, subject, categoryId, description, userId, created_at, updated_at from questions where userId = $1`

	rows, err := m.DB.QueryContext(ctx, query, userId)

	if err != nil {
		return questions, err
	}

	for rows.Next(){
		var question models.Question
		err:= rows.Scan(
			&question.ID,
			&question.Subject,
			&question.Category,
			&question.Description,
			&question.UserId,
			&question.CreatedAt,
			&question.UpdatedAt,
		)

		if err != nil {
			return questions, err
		}
		questions = append(questions, question)
	}

	if  err = rows.Err(); err!=nil {
		return questions, err
		
	}

	return questions, nil
}

//get all question
func (m *postgreDBRepo) GetQuestions() ([]models.Question, error){

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var questions []models.Question

	query := `
	select id, subject, categoryId, description, userId, created_at, updated_at from questions`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return questions, err
	}

	for rows.Next(){
		var question models.Question
		err:= rows.Scan(
			&question.ID,
			&question.Subject,
			&question.Category,
			&question.Description,
			&question.UserId,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		
		if err != nil {
			return questions, err
		}
		questions = append(questions, question)
	}

	if  err = rows.Err(); err!=nil {
		return questions, err
		
	}

	return questions, nil
}