package gmq

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_MYSQL    = "mysql"
	DB_POSTGRES = "postgres"
)

type DbQueue struct {
	Name, Dsn string
	conn      *sql.DB
}
