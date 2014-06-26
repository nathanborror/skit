package skits

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var repo = SkitSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// ViewHandler displays a skit
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	// get logged in user
	session, _ := store.Get(r, "authenticated-user")
	user := session.Values["hash"]
	if user == nil {
		user = ""
	} else {
		user = user.(string)
	}

	// Load the skit
	s, err := repo.Load(hash)
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
		"session":  user,
		"skit":     s,
		"children": c,
	})
}

// NewHandler creates a new skit
func NewHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "skit_form", nil)
}

// EditHandler edits a skit
func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	b, err := repo.Load(hash)
	if err != nil {
		b = &Skit{Hash: hash}
	}
	render.RenderTemplate(w, "skit_form", b)
}

// SaveHandler saves a skit
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "authenticated-user")
	if session.Values["hash"] == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	hash := r.FormValue("hash")
	parent := r.FormValue("parent")
	root := r.FormValue("root")
	user := session.Values["hash"].(string)
	text := r.FormValue("text")

	if hash == "" {
		hash = GenerateSkitHash(text)
	}

	if root == "" {
		root = hash
	}

	s := &Skit{Hash: hash, Parent: parent, Root: root, User: user, Text: text}
	err := repo.Save(s)
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
