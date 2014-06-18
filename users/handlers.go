package users

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/skit/render"
	"hash/fnv"
	"io"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var repo = NewSqlUserRepository("db.sqlite3")

// SigninViewHandler signs a user in
func SigninViewHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "authenticated-user")

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Generate password hash
		hasher := md5.New()
		io.WriteString(hasher, password)
		passwordHash := fmt.Sprintf("%x", hasher.Sum(nil))

		u, err := repo.LoadWithPassword(email, string(passwordHash))
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		session.Values["hash"] = u.Hash
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}

	render.RenderTemplate(w, "user_signin", nil)
}

// SignoutViewHandler signs a user out
func SignoutViewHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "authenticated-user")
	session.Values["hash"] = nil
	session.Save(r, w)
	render.RenderTemplate(w, "user_signout", nil)
}

// RegisterViewHandler registers a new user
func RegisterViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		name := r.FormValue("name")

		if email == "" && password == "" {
			http.Redirect(w, r, "/register", http.StatusFound)
		}

		// Generate password hash
		hasher := md5.New()
		io.WriteString(hasher, password)
		passwordHash := fmt.Sprintf("%x", hasher.Sum(nil))

		// Generate hash
		hasher = fnv.New32a()
		io.WriteString(hasher, email)
		hash := fmt.Sprintf("%x", hasher.Sum(nil))

		u := &User{Email: email, Password: string(passwordHash), Hash: hash, Name: name}
		err := repo.Save(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
	}

	render.RenderTemplate(w, "user_register", nil)
}
