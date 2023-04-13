package handlers

import (
	"log"
	"net/http"

	"github.com/mr-keppy/go-forum/internal/config"
	"github.com/mr-keppy/go-forum/internal/driver"
	"github.com/mr-keppy/go-forum/internal/forms"
	"github.com/mr-keppy/go-forum/internal/models"
	"github.com/mr-keppy/go-forum/internal/render"
	"github.com/mr-keppy/go-forum/internal/repository"
	"github.com/mr-keppy/go-forum/internal/repository/dbrepo"
)

var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// Create new Repo
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {

	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL,a),
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

	if err!=nil{
		log.Println(err)
	}

	email:= r.Form.Get("email")
	password:= r.Form.Get("password")

	form := forms.New(r.PostForm)

	form.Required("email","password")

	if !form.Valid(){
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
	}

	id, _, err := m.DB.Authenticate(email,password)
	if err!=nil{
		log.Println(err)
		m.App.Session.Put(r.Context(),"error","invalid user input")
		http.Redirect(w, r, "/user/login",http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(),"user_id",id)
	m.App.Session.Put(r.Context(),"flash","login successfully")
	http.Redirect(w, r, "/",http.StatusSeeOther)

}
//logout
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
