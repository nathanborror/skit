package skits

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/skit/render"
	"github.com/nathanborror/skit/users"
	"io"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var repo = NewSqlSkitRepository("db.sqlite3")
var userRepo = users.NewSqlUserRepository("db.sqlite3")

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

	// TODO: Get child skits

	render.RenderTemplate(w, "skit_view", map[string]interface{}{
		"skit":  s,
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
		http.Redirect(w, r, "/user/signin/", http.StatusFound)
		return
	}

	hash := r.FormValue("hash")
	parent := r.FormValue("parent")
	user := session.Values["hash"].(string)
	text := r.FormValue("text")

	if hash == "" {
		m := md5.New()
		io.WriteString(m, text)
		hash = fmt.Sprintf("%x", m.Sum(nil))
	}

	s := &Skit{Hash: hash, Parent: parent, User: user, Text: text}
	err := repo.Save(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/skit/view/"+hash, http.StatusFound)
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
