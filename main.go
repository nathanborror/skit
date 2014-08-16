package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/hubspoke"
	"github.com/nathanborror/gommon/markdown"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/skit/items"
	"github.com/nathanborror/skit/messages"
)

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))
var itemRepo = items.ItemSQLRepository("db.sqlite3")
var messageRepo = messages.MessageSQLRepository("db.sqlite3")
var authRepo = auth.AuthSQLRepository("db.sqlite3")

func init() {
	_ = render.RegisterTemplateFunction("markdown", markdown.Markdown)
	_ = render.RegisterTemplateFunction("isQuestion", isQuestion)
	_ = render.RegisterTemplateFunction("slice", slice)
	_ = render.RegisterTemplateFunction("isLight", isLight)

	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func isQuestion(text string) bool {
	if len(text) <= 0 {
		return false
	}
	if string(text[len(text)-1]) == "?" {
		return true
	}
	return false
}

func isLight(text string) bool {
	if len(text) <= 0 {
		return false
	}
	rgb := strings.Split(text, ",")
	r, _ := strconv.Atoi(rgb[0])
	g, _ := strconv.Atoi(rgb[1])
	b, _ := strconv.Atoi(rgb[2])
	value := r + g + b
	return value > 580
}

func slice(args ...interface{}) string {
	if len(args) == 2 {
		text := args[1].(string)
		valueString := args[0].(string)
		values := strings.Split(valueString, ":")

		if len(values) > 1 {
			val1, _ := strconv.Atoi(values[0])
			val2, _ := strconv.Atoi(values[1])

			if val2 == -1 {
				result := text[val1 : len(text)-1]
				return string(result)
			}

			result := text[val1:val2]
			return string(result)
		}

		val, _ := strconv.Atoi(values[0])
		return string(text[val])
	}
	return args[0].(string)
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
	var m []*messages.Message

	// If a item exists then load it's children, otherwise load all the
	// items a user is involved in.
	if i != nil {
		c, err = itemRepo.ListWithParent(i.Hash)
		check(err, w)

		m, _ = messageRepo.ListForRoot(i.Root)
	} else {
		c, err = itemRepo.ListWithUser(u.Hash)
		check(err, w)
	}

	render.Render(w, r, "home", map[string]interface{}{
		"request":  r,
		"user":     u,
		"item":     i,
		"items":    c,
		"messages": m,
	})
}

func peopleHandler(w http.ResponseWriter, r *http.Request) {
	u, err := authRepo.List(10)
	check(err, w)

	render.Render(w, r, "users", map[string]interface{}{
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

	render.Render(w, r, "home", map[string]interface{}{
		"request": r,
		"user":    u,
		"items":   i,
	})
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	root := vars["root"]

	u, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	i, err := itemRepo.Load(root)
	check(err, w)

	m, err := messageRepo.ListForRoot(root)
	check(err, w)

	render.Render(w, r, "message_view", map[string]interface{}{
		"request":  r,
		"user":     u,
		"messages": m,
		"item":     i,
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

	// Message
	r.HandleFunc("/m/save", auth.LoginRequired(messages.SaveHandler))
	r.HandleFunc("/m/{hash:[a-zA-Z0-9-]+}/delete", auth.LoginRequired(messages.DeleteHandler))
	r.HandleFunc("/m/{root:[a-zA-Z0-9-]+}", auth.LoginRequired(messagesHandler))

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
