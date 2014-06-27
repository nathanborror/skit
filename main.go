package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/hubspoke"
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
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	http.Redirect(w, r, "/u/"+user.Hash, http.StatusFound)
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
	go hubspoke.H.Run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)

	// Skit
	r.HandleFunc("/s/save", auth.LoginRequired(skits.SaveHandler))
	r.HandleFunc("/s/{hash:[a-zA-Z0-9-]+}/edit", auth.LoginRequired(skits.EditHandler))
	r.HandleFunc("/s/{hash:[a-zA-Z0-9-]+}/delete", auth.LoginRequired(skits.DeleteHandler))
	r.HandleFunc("/s/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(skits.ViewHandler))

	r.HandleFunc("/ws", hubspoke.SpokeHandler)
	r.HandleFunc("/u/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(userHomeHandler))
	r.HandleFunc("/", auth.LoginRequired(homeHandler))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
