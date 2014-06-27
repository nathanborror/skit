package skits

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
)

var repo = SkitSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// ViewHandler displays a skit
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	// Load the skit
	s, err := repo.Load(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load parent skits
	p, err := repo.ListParents(s.Root)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := repo.ListWithParent(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Render(w, r, "skit_view", map[string]interface{}{
		"request":  r,
		"skit":     s,
		"children": c,
		"parents":  p,
	})
}

// EditHandler edits a skit
func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	s, err := repo.Load(hash)
	if err != nil {
		s = &Skit{Hash: hash}
	}
	render.RenderTemplate(w, "skit_form", map[string]interface{}{
		"request": r,
		"skit":    s,
	})
}

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

	err := repo.Delete(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
