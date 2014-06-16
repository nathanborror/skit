package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/skit/skits"
	"github.com/nathanborror/skit/users"
	"github.com/nathanborror/skit/render"
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

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		render.RenderJSON(w, map[string]interface{}{
			"skits": s,
		})
		return
	}

	render.RenderTemplate(w, "home", map[string]interface{}{
		"session": session.Values["hash"],
		"skits": s,
		"connections": h.connections,
	})
}

func main() {
	go h.run()

	r := mux.NewRouter()

	// Users
	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/signin/", users.SigninViewHandler)
	u.HandleFunc("/signout/", users.SignoutViewHandler)
	u.HandleFunc("/register/", users.RegisterViewHandler)

	// Skit
	s := r.PathPrefix("/skit").Subrouter()
	s.HandleFunc("/view/{hash:[a-zA-Z0-9-]+}", skits.ViewHandler)
	s.HandleFunc("/edit/{hash:[a-zA-Z0-9-]+}", skits.EditHandler)
	s.HandleFunc("/delete/{hash:[a-zA-Z0-9-]+}", skits.DeleteHandler)
	s.HandleFunc("/save/", skits.SaveHandler)
	s.HandleFunc("/new/", skits.NewHandler)

	r.HandleFunc("/ws", socketHandler)
	r.HandleFunc("/", homeHandler)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
