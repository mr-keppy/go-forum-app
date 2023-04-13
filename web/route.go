package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mr-keppy/go-forum/internal/config"
	"github.com/mr-keppy/go-forum/internal/handlers"
)


func routes(app *config.AppConfig) http.Handler{
	/*mux := pat.New()
	mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
	*/

	mux:=chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about",handlers.Repo.About)
	mux.Get("/user/register",handlers.Repo.Register)
	mux.Get("/ask-question",handlers.Repo.AskQuestion)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	fileServer:= http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static",fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		

	})

	return mux
}