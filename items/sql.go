package items

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlItemRepository struct {
	dbmap *modl.DbMap
}

// ItemSQLRepository returns a new sqlItemRepository or panics if it cannot
func ItemSQLRepository(filename string) ItemRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlItemRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Item{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlItemRepository) Load(hash string) (*Item, error) {
	obj := []*Item{}
	err := r.dbmap.Select(&obj, "SELECT * FROM item WHERE hash=?", hash)
	if len(obj) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(obj))
	}
	return obj[0], err
}

func (r *sqlItemRepository) Save(item *Item) error {
	n, err := r.dbmap.Update(item)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(item)
	}
	return err
}

func (r *sqlItemRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM item WHERE hash=?", hash)
	return err
}

func (r *sqlItemRepository) List(limit int) ([]*Item, error) {
	obj := []*Item{}
	err := r.dbmap.Select(&obj, "SELECT * FROM item WHERE parent = '' ORDER BY created DESC LIMIT ?", limit)
	return obj, err
}

func (r *sqlItemRepository) ListWithParent(parent string) ([]*Item, error) {
	obj := []*Item{}
	err := r.dbmap.Select(&obj, "SELECT * FROM item WHERE parent = ? ORDER BY created DESC", parent)
	return obj, err
}

func (r *sqlItemRepository) ListWithUser(user string) ([]*Item, error) {
	obj := []*Item{}
	err := r.dbmap.Select(&obj, "SELECT DISTINCT * FROM item WHERE hash IN (SELECT DISTINCT root FROM item WHERE user = ?) ORDER BY created DESC", user)
	return obj, err
}

func (r *sqlItemRepository) ListParents(hash string) ([]*Item, error) {
	obj := []*Item{}
	err := r.dbmap.Select(&obj, "SELECT * FROM item WHERE root = ? ORDER BY created DESC", hash)
	return obj, err
}

func (r *sqlItemRepository) Archive(item *Item) error {
	item.IsArchived = true
	_, err := r.dbmap.Update(item)
	return err
}

func (r *sqlItemRepository) UnArchive(item *Item) error {
	item.IsArchived = false
	_, err := r.dbmap.Update(item)
	return err
}
