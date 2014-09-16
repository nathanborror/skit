package items

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
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
	isarchived := false

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

	i := &Item{Hash: hash, Parent: parent, Root: root, User: user.Hash, Text: text, Color: color, Created: created, IsArchived: isarchived}
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

	i, err := repo.Load(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if i.User != u.Hash {
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

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	i, err := repo.Load(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if i.IsArchived {
		err = repo.UnArchive(i)
	} else {
		err = repo.Archive(i)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
