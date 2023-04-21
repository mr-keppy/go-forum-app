package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mr-keppy/go-forum/internal/config"
	"github.com/mr-keppy/go-forum/internal/driver"
	"github.com/mr-keppy/go-forum/internal/forms"
	"github.com/mr-keppy/go-forum/internal/helpers"
	"github.com/mr-keppy/go-forum/internal/models"
	"github.com/mr-keppy/go-forum/internal/render"
	"github.com/mr-keppy/go-forum/internal/repository"
	"github.com/mr-keppy/go-forum/internal/repository/dbrepo"
)

var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Create new Repo
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {

	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

// login screen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// login screen validation
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)

	form.Required("email", "password")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid user input")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "login successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// logout
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "user/login", http.StatusSeeOther)
}

// Home page is the home handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// register page is the register handler
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// ask-question page is the ask-question handler
func (m *Repository) AskQuestion(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "ask-question.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

	fmt.Println("AskQuestion")

}

// Post AskQuestion page
func (m *Repository) PostAskQuestion(w http.ResponseWriter, r *http.Request) {

	fmt.Println("PostAskQuestion")

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var question models.Question

	question.Subject = r.Form.Get("subject")
	question.Description = r.Form.Get("description")
	question.Category = r.Form.Get("category")
	question.UserId = 1

	fmt.Println(question)

	form := forms.New(r.PostForm)

	log.Println(question)
	form.Required("subject", "description", "category")
	form.MinLength("subject", 5, r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["question"] = question

		render.Template(w, r, "ask-question.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.CreateQuestion(question)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "question", question)

	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}

// view question
func( m *Repository) ViewQuestion(w http.ResponseWriter, r *http.Request){
	//questionId,  _ := strconv.Atoi(r.URL.Query().Get("id"))
	questionId, _:= strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/view-question/"))
	var question models.Question
	
	fmt.Printf("questionid",questionId)

	question, err := m.DB.GetQuestionByID(questionId)

	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	
	data := make(map[string]interface{})
	data["question"] = question
	
	render.Template(w, r, "view-questions.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// edit question
func( m *Repository) EditQuestion(w http.ResponseWriter, r *http.Request){
	questionId,  _ := strconv.Atoi(r.URL.Query().Get("id"))
	
	var question models.Question

	fmt.Printf("questionid",questionId)

	question, err := m.DB.GetQuestionByID(questionId)

	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	
	data := make(map[string]interface{})
	data["question"] = question
	
	render.Template(w, r, "view-questions.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// edit post question
func( m *Repository) PostEditQuestion(w http.ResponseWriter, r *http.Request){
	questionId,  _ := strconv.Atoi(r.URL.Query().Get("id"))
	var question models.Question

	question, err := m.DB.GetQuestionByID(questionId)

	if err!=nil{
		helpers.ServerError(w, err)
		return
	}
	
	data := make(map[string]interface{})
	data["question"] = question
	
	render.Template(w, r, "view-questions.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// About this about handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, World"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// show all questions
func (m *Repository) AllQuestions(w http.ResponseWriter, r *http.Request) {

	questions, err := m.DB.GetQuestions()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["questions"] = questions
	render.Template(w, r, "questions.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
