package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/skit/skits"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var skitRepo = skits.SkitSQLRepository("db.sqlite3")

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
	http.Redirect(w, r, "/u/"+user.(string), http.StatusFound)
}

func userHomeHandler(w http.ResponseWriter, r *http.Request) {
	u, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	s, err := skitRepo.ListWithUser(u.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Render(w, r, "user_view", map[string]interface{}{
		"skit":     "",
		"children": s,
		"user":     u,
		"request":  r,
	})
}

var r = mux.NewRouter()

func main() {
	go h.run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)

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
