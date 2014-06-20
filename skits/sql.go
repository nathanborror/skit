package skits

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlSkitRepository struct {
	dbmap *modl.DbMap
}

// NewSqlSkitRepository returns a new sqlSkitRepository or panics if it cannot
func NewSqlSkitRepository(filename string) SkitRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlSkitRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Skit{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlSkitRepository) Load(hash string) (*Skit, error) {
	obj := []*Skit{}
	err := r.dbmap.Select(&obj, "SELECT * FROM skit WHERE hash=?", hash)
	if len(obj) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(obj))
	}
	return obj[0], err
}

func (r *sqlSkitRepository) Save(skit *Skit) error {
	n, err := r.dbmap.Update(skit)
	if err != nil {
		panic(err)
		return err
	}
	if n == 0 {
		err = r.dbmap.Insert(skit)
	}
	return err
}

func (r *sqlSkitRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM skit WHERE hash=?", hash)
	return err
}

func (r *sqlSkitRepository) List(limit int) ([]*Skit, error) {
	obj := []*Skit{}
	err := r.dbmap.Select(&obj, "SELECT * FROM skit WHERE parent = '' ORDER BY modified DESC LIMIT ?", limit)
	return obj, err
}

func (r *sqlSkitRepository) ListWithParent(parent string) ([]*Skit, error) {
	obj := []*Skit{}
	err := r.dbmap.Select(&obj, "SELECT * FROM skit WHERE parent = ? ORDER BY modified DESC", parent)
	return obj, err
}

func (r *sqlSkitRepository) ListWithUser(user string) ([]*Skit, error) {
	obj := []*Skit{}
	err := r.dbmap.Select(&obj, "SELECT DISTINCT * FROM skit WHERE hash IN (SELECT DISTINCT root FROM skit WHERE user = ?)", user)
	return obj, err
}
