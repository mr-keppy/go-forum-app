package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/mr-keppy/go-forum/internal/helpers"
)

func WriteToConsole(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Println("Hit the page")
		next.ServeHTTP(w,r)
	})
}

func NoSurf(next http.Handler) http.Handler{
	csrfHandler:= nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !helpers.IsAuthenticated(r){
			session.Put(r.Context(),"error","Login first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w,r)
	})
}