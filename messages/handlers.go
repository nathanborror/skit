package messages

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
)

var repo = MessageSQLRepository("db.sqlite3")
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

	if hash == "" {
		hash = GenerateMessageHash(text)
	}

	if root == "" {
		root = hash
	}

	m := &Message{Hash: hash, Parent: parent, Root: root, User: user.Key, Text: text}
	err = repo.Save(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/i/"+root, http.StatusFound)
}

// DeleteHandler deletes a item
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	u, err := auth.GetAuthenticatedUser(r)

	m, err := repo.Load(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m.User != u.Key {
		render.Render(w, r, "error", map[string]interface{}{
			"error":   "You can only delete items you created.",
			"request": r,
		})
		return
	}

	err = repo.Delete(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
