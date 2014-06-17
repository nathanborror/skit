package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/skit/render"
	"github.com/nathanborror/skit/skits"
	"github.com/nathanborror/skit/users"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var skitRepo = skits.NewSqlSkitRepository("db.sqlite3")

func init() {
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "authenticated-user")

	s, err := skitRepo.List(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Render(w, r, "home", map[string]interface{}{
		"session":     session.Values["hash"],
		"skit":        "",
		"children":    s,
		// "connections": h.connections,
	})
}

func userHomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	session, _ := store.Get(r, "authenticated-user")

	s, err := skitRepo.ListWithUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Render(w, r, "home", map[string]interface{}{
		"session":     session.Values["hash"],
		"skit":        "",
		"children":    s,
		// "connections": h.connections,
	})
}

func main() {
	go h.run()

	r := mux.NewRouter()

	// Users
	u := r.PathPrefix("/u").Subrouter()
	u.HandleFunc("/signin", users.SigninViewHandler)
	u.HandleFunc("/signout", users.SignoutViewHandler)
	u.HandleFunc("/register", users.RegisterViewHandler)

	// Skit
	s := r.PathPrefix("/s").Subrouter()
	s.HandleFunc("/{hash:[a-zA-Z0-9-]+}", skits.ViewHandler)
	s.HandleFunc("/{hash:[a-zA-Z0-9-]+}/edit", skits.EditHandler)
	s.HandleFunc("/{hash:[a-zA-Z0-9-]+}/delete", skits.DeleteHandler)
	s.HandleFunc("/save", skits.SaveHandler)
	s.HandleFunc("/new", skits.NewHandler)

	r.HandleFunc("/ws", socketHandler)
	r.HandleFunc("/u/{user:[a-zA-Z0-9-]+}", userHomeHandler)
	r.HandleFunc("/", homeHandler)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
