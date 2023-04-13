package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mr-keppy/go-forum/internal/config"
	"github.com/mr-keppy/go-forum/internal/driver"
	"github.com/mr-keppy/go-forum/internal/handlers"
	"github.com/mr-keppy/go-forum/internal/helpers"
	"github.com/mr-keppy/go-forum/internal/models"
	"github.com/mr-keppy/go-forum/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//this is main
func run() (*driver.DB, error){
	gob.Register(models.User{})

	inProduction := *flag.Bool("production",true, "Application is in production")
	useCache:= *flag.Bool("cache",true, "use template cache")
	dbName:= flag.String("dbname","forum","DB Name")
	dbHost:= flag.String("dbhost","localhost","DB ost")
	dbUser:= flag.String("dbuser","kishorpadmanabhan","DB User")
	dbPass:= flag.String("dbpass","","DB Pass")
	dbPort:= flag.String("dbport","5432","DB port")
	dbSSL:= flag.String("dbssl","disable","DB SSL")


	app.InProduction = *&inProduction

	infoLog = log.New(os.Stdout, "INFO\t",log.Ldate | log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout,"ERROR\t",log.Ldate | log.Ltime| log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()

	session.Lifetime = 24*time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to db
	log.Println("connect to db")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",*dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL))
	if err != nil {
		log.Fatal("Error while connecting db")
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = *&useCache
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	return db, nil
}

func main() {

	db, err:= run()

	if(err != nil){
		log.Fatal(err)
	}
	// what going to store in session

	defer db.SQL.Close()


	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
	// fmt.Println((fmt.Sprintf("Starting applicaiton on port #:%d", portNumber)))
	//_ = http.ListenAndServe(portNumber, nil)

}