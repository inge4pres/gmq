package gmq

import (
	"database/sql"
	"encoding/base64"
	_ "github.com/go-sql-driver/mysql"
)

type DbQueue struct {
	Name, Vendor, Dsn string
	conn              *sql.DB
}

func openConn(db *DbQueue) (err error) {
	db.conn, err = sql.Open(db.Vendor, db.Dsn)
	return err
}

func (db *DbQueue) Push(o []byte) error {
	if err := openConn(db); err != nil {
		return err
	}
	defer db.conn.Close()
	_, err := db.conn.Exec("INSERT INTO "+db.Name+" VALUES(NULL,?,?)", base64.StdEncoding.EncodeToString(o), false)
	return err
}

func (db *DbQueue) Pop() ([]byte, error) {
	if err := openConn(db); err != nil {
		return nil, err
	}
	defer db.conn.Close()
	var decode string
	var id int64
	if err := db.conn.QueryRow("SELECT MIN(id) FROM " + db.Name).Scan(&id); err != nil {
		return nil, err
	}
	if id > 0 {
		if err := db.conn.QueryRow("SELECT message FROM "+db.Name+" WHERE id = ?", id).Scan(&decode); err != nil {
			return nil, err
		}
		if _, err := db.conn.Exec("UPDATE "+db.Name+" SET processed = 1 WHERE id = ?", id); err != nil {
			return nil, err
		}
	}
	db.sync()
	return base64.StdEncoding.DecodeString(decode)
}

func (db *DbQueue) sync() {
	db.conn.Exec("DELETE FROM " + db.Name + " WHERE processed = 1")
}
