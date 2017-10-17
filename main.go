package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/Jleagle/canihave/logger"
	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/login", loginHandler)
	r.Get("/logout", logoutHandler)
	r.Get("/login-callback", loginCallbackHandler)
	r.Get("/shuffle", shuffleHandler)

	log.Fatal(http.ListenAndServe(":8084", r))
}

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) {

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		logger.Err("No caller information")
	}
	folder := path.Dir(file)

	// Load templates needed
	always := []string{folder + "/templates/header.html", folder + "/templates/footer.html", folder + "/templates/" + page + ".html"}

	t, err := template.New("t").ParseFiles(always...)
	if err != nil {
		logger.ErrExit(err.Error())
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		logger.ErrExit(err.Error())
	}
}

type templateVars struct {
	LoggedIn bool
	Flashes  []interface{}
}
