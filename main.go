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
	"github.com/nathanborror/skit/items"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var itemRepo = items.ItemSQLRepository("db.sqlite3")
var authRepo = auth.AuthSQLRepository("db.sqlite3")

func init() {
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	u, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	http.Redirect(w, r, "/u/"+u.Hash, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	// Load user
	u, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Load item
	i, _ := itemRepo.Load(hash)

	var c []*items.Item

	// If a item exists then load it's children, otherwise load all the
	// items a user is involved in.
	if i != nil {
		c, err = itemRepo.ListWithParent(i.Hash)
		check(err, w)
	} else {
		c, err = itemRepo.ListWithUser(u.Hash)
		check(err, w)
	}

	render.Render(w, r, "home", map[string]interface{}{
		"item":    i,
		"items":   c,
		"user":    u,
		"request": r,
	})
}

func peopleHandler(w http.ResponseWriter, r *http.Request) {
	u, err := authRepo.List(10)
	check(err, w)

	render.Render(w, r, "user_list", map[string]interface{}{
		"users":   u,
		"request": r,
	})
}

func personHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	u, err := authRepo.Load(hash)
	check(err, w)

	i, err := itemRepo.ListWithUser(u.Hash)
	check(err, w)

	render.Render(w, r, "user_view", map[string]interface{}{
		"request": r,
		"user":    u,
		"items":   i,
	})
}

var r = mux.NewRouter()

func main() {
	go hubspoke.Hub.Run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)

	// Item
	r.HandleFunc("/i/save", auth.LoginRequired(items.SaveHandler))
	r.HandleFunc("/i/{hash:[a-zA-Z0-9-]+}/delete", auth.LoginRequired(items.DeleteHandler))
	r.HandleFunc("/i/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(homeHandler))

	r.HandleFunc("/ws", hubspoke.SpokeHandler)
	r.HandleFunc("/u/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(personHandler))
	r.HandleFunc("/u", auth.LoginRequired(peopleHandler))
	r.HandleFunc("/", auth.LoginRequired(rootHandler))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
