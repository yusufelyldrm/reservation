package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

// DB holds the database connection pool.
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

// ConnectSQL creates database pool for Postgres.
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil

}

// testDB checks if database is alive.
func testDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil

}

// NewDatabase creates a new database for the application.
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pqx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
