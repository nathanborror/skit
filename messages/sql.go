package messages

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlMessageRepository struct {
	dbmap *modl.DbMap
}

// MessageSQLRepository returns a new sqlMessageRepository or panics if it cannot
func MessageSQLRepository(filename string) MessageRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlMessageRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Message{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlMessageRepository) Load(hash string) (*Message, error) {
	obj := []*Message{}
	err := r.dbmap.Select(&obj, "SELECT * FROM message WHERE hash=?", hash)
	if len(obj) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(obj))
	}
	return obj[0], err
}

func (r *sqlMessageRepository) Save(message *Message) error {
	n, err := r.dbmap.Update(message)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(message)
	}
	return err
}

func (r *sqlMessageRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM message WHERE hash=?", hash)
	return err
}

func (r *sqlMessageRepository) List(parent string) ([]*Message, error) {
	obj := []*Message{}
	err := r.dbmap.Select(&obj, "SELECT * FROM message WHERE parent = ? ORDER BY created ASC", parent)
	return obj, err
}

func (r *sqlMessageRepository) ListForRoot(root string) ([]*Message, error) {
	obj := []*Message{}
	err := r.dbmap.Select(&obj, "SELECT * FROM message WHERE root = ? ORDER BY created ASC", root)
	return obj, err
}
