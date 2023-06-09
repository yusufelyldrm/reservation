package dbrepo

import (
	"database/sql"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig //pointer because we want to change the value of AppConfig
	DB  *sql.DB           //
}

type testDBRepo struct {
	App *config.AppConfig //pointer because we want to change the value of AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
