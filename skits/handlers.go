package skits

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
)

var repo = SkitSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a skit
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

	if hash == "" {
		hash = GenerateSkitHash(text)
	}

	if root == "" {
		root = hash
	}

	s := &Skit{Hash: hash, Parent: parent, Root: root, User: user.Hash, Text: text}
	err = repo.Save(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/s/"+hash, http.StatusFound)
}

// DeleteHandler deletes a skit
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
