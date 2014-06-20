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

// GeneratePasswordHash returns a hashed password
func GeneratePasswordHash(password string) (hash string) {
	hasher := md5.New()
	io.WriteString(hasher, password)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// GenerateUserHash returns a hash that represents a unique user ID
func GenerateUserHash(s string) (hash string) {
	hasher := fnv.New32a()
	io.WriteString(hasher, s)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Authenticate authenticates and returns a user object
func Authenticate(email string, password string, w http.ResponseWriter, r *http.Request) (user *User) {
	hash := GeneratePasswordHash(password)
	u, err := repo.LoadWithPassword(email, hash)
	if err != nil {
		return nil
	}

	// Update session
	if (w != nil && r != nil) {
		session, _ := store.Get(r, "authenticated-user")
		session.Values["hash"] = u.Hash
		session.Save(r, w)
	}

	return u
}

// SigninViewHandler signs a user in
func SigninViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		u := Authenticate(email, password, w, r)
		if u == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
		}

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
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" && password == "" {
			http.Redirect(w, r, "/register", http.StatusFound)
		}

		hash := GenerateUserHash(email)
		passwordHash := GeneratePasswordHash(password)

		// If user already exists, sign them in and send them to '/'
		u := Authenticate(email, password, w, r)
		if u != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		u = &User{Email: email, Password: passwordHash, Hash: hash, Name: name}
		err := repo.Save(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Auth user and redirect them to '/'
		u = Authenticate(email, password, w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	render.RenderTemplate(w, "user_register", nil)
}
