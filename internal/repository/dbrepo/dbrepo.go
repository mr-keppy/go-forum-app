package dbrepo

import (
	"database/sql"

	"github.com/mr-keppy/go-forum/internal/config"
	"github.com/mr-keppy/go-forum/internal/repository"
)

type postgreDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgreDBRepo{
		App: a,
		DB: conn,
	}
}