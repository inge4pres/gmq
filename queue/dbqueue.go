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

func (db DbQueue) GetLength() (int, error) {
	var ret int
	if err := openConn(&db); err != nil {
		return -1, err
	}
	defer db.conn.Close()
	err := db.conn.QueryRow("SELECT MAX(id) FROM ?", db.Name).Scan(ret)
	return ret, err
}

func (db DbQueue) Create(name string) (QueueInterface, error) {
	dbq := DbQueue{Name: name}
	if err := openConn(&db); err != nil {
		return nil, err
	}
	defer db.conn.Close()
	sql := "DROP TABLE IF EXISTS " + name + ";"
	if _, err := db.conn.Exec(sql); err != nil {
		return nil, err
	}
	sql = "CREATE TABLE " + name + " ( id serial PRIMARY KEY NOT NULL AUTO_INCREMENT, message MEDIUMTEXT, processed BOOL DEFAULT FALSE, UNIQUE INDEX proc_idx USING BTREE (id, processed));"
	if _, err := db.conn.Exec(sql); err != nil {
		return nil, err
	}

	return dbq, nil
}

func (db DbQueue) Push(o []byte) error {
	if err := openConn(&db); err != nil {
		return err
	}
	defer db.conn.Close()
	_, err := db.conn.Exec("INSERT INTO "+db.Name+" VALUES(NULL,?,?)", base64.StdEncoding.EncodeToString(o), false)
	return err
}

func (db DbQueue) Pop() ([]byte, error) {
	if err := openConn(&db); err != nil {
		return nil, err
	}
	defer db.conn.Close()
	var decode string
	var id int64
	if err := db.conn.QueryRow("SELECT MIN(id) FROM " + db.Name + " WHERE processed = 0").Scan(&id); err != nil {
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

func (db DbQueue) sync() {
	db.conn.Exec("DELETE FROM " + db.Name + " WHERE processed = 1")
}
