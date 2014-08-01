package items

import (
	"time"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
)

var repo = ItemSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	hash := r.FormValue("hash")
	parent := r.FormValue("parent")
	root := r.FormValue("root")
	text := r.FormValue("text")
	color := r.FormValue("color")
	created := time.Now()

	if hash == "" {
		hash = GenerateItemHash(text)
	}

	if root == "" {
		root = hash
	}

	item, err := repo.Load(hash)
	if err == nil {
		created = item.Created
	}

	i := &Item{Hash: hash, Parent: parent, Root: root, User: user.Hash, Text: text, Color: color, Created: created}
	err = repo.Save(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/i/"+hash, http.StatusFound)
}

// DeleteHandler deletes a item
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	u, err := auth.GetAuthenticatedUser(r)

	s, err := repo.Load(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.User != u.Hash {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = repo.Delete(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
