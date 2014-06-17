package render

import (
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"net/http"
	"regexp"
	"text/template"
)

var funcMap = template.FuncMap{
	"markdown": markDowner,
	"initials": initials,
}

// RenderTemplate renders a given template along with any data passed
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates := template.New("").Funcs(funcMap)
	_, err := templates.ParseGlob("templates/*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderJSON returns marshalled JSON
func RenderJSON(w http.ResponseWriter, data interface{}) {
	obj, _ := json.MarshalIndent(data, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(obj)
}

func markDowner(args ...interface{}) string {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return string(s)
}

func initials(args ...interface{}) string {
	s := fmt.Sprintf("%s", args...)
	re := regexp.MustCompile("[^A-Z]")
	return re.ReplaceAllString(s, "")
}
