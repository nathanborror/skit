package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/skit/render"
	"github.com/nathanborror/skit/skits"
	"github.com/nathanborror/skit/users"
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
	user := session.Values["hash"]
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	s, err := skitRepo.ListWithUser(user.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Render(w, r, "home", map[string]interface{}{
		"session":  user.(string),
		"skit":     "",
		"children": s,
		// "connections": h.connections,
	})
}

func userHomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	session, _ := store.Get(r, "authenticated-user")

	s, err := skitRepo.ListWithUser(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Render(w, r, "home", map[string]interface{}{
		"session":  session.Values["hash"],
		"skit":     "",
		"children": s,
		// "connections": h.connections,
	})
}

var r = mux.NewRouter()

func main() {
	go h.run()

	// Users
	r.HandleFunc("/signin", users.SigninViewHandler)
	r.HandleFunc("/signout", users.SignoutViewHandler)
	r.HandleFunc("/register", users.RegisterViewHandler)

	// Skit
	s := r.PathPrefix("/s").Subrouter()
	s.HandleFunc("/save", skits.SaveHandler)
	s.HandleFunc("/new", skits.NewHandler)
	s.HandleFunc("/{hash:[a-zA-Z0-9-]+}/edit", skits.EditHandler)
	s.HandleFunc("/{hash:[a-zA-Z0-9-]+}/delete", skits.DeleteHandler)
	s.HandleFunc("/{hash:[a-zA-Z0-9-]+}", skits.ViewHandler)

	r.HandleFunc("/ws", socketHandler)
	r.HandleFunc("/u/{hash:[a-zA-Z0-9-]+}", userHomeHandler)
	r.HandleFunc("/", homeHandler)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
