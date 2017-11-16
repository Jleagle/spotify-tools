package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi"
	"github.com/zmb3/spotify"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {

	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/logout", logoutHandler)
	r.Get("/login", loginHandler)
	r.Get("/login-callback", loginCallbackHandler)
	r.Get("/shuffle", shuffleHandler)
	r.Get("/shuffle/{playlist}/{new}", shuffleActionHandler)
	r.Get("/random", randomHandler)
	r.Get("/random/{type}", randomHandler)
	r.Get("/duplicates", duplicatesHandler)
	r.Get("/duplicates/{playlist}/{new}", duplicatesActionHandler)
	r.Get("/artist/{artist}", artistHandler)
	r.Get("/album/{album}", albumHandler)
	r.Get("/track/{track}", trackHandler)
	r.Get("/user/{user}", userHandler)
	r.Get("/user/{user}/playlist/{playlist}", playlistHandler)

	// Assets
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	log.Fatal(http.ListenAndServe(":8084", r))
}

func returnTemplate(w http.ResponseWriter, r *http.Request, page string, pageData structs.TemplateVars, err error) {

	// todo, log errors

	pageData.LoggedIn = session.IsLoggedIn(r)
	pageData.Flashes = session.GetFlashes(w, r)
	pageData.Highlight = r.URL.Query().Get("highlight")
	pageData.LoggedInID = session.Read(r, session.UserID)

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		//logger.Err("No caller information")
	}
	folder := path.Dir(file)

	// Load templates needed
	always := []string{
		folder + "/templates/header.html",
		folder + "/templates/footer.html",
		folder + "/templates/" + page + ".html",
		folder + "/templates/includes/album.html",
	}

	t, err := template.New("t").Funcs(getTemplateFuncMap()).ParseFiles(always...)
	if err != nil {
		//logger.ErrExit(err.Error())
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		//logger.ErrExit(err.Error())
	}
}

func returnLoggedOutTemplate(w http.ResponseWriter, r *http.Request, err error) {

	vars := structs.TemplateVars{}
	vars.ErrorMessage = "Please login"

	returnTemplate(w, r, "error", vars, err)
}

func getTemplateFuncMap() map[string]interface{} {
	return template.FuncMap{
		"join":  func(a []string) string { return strings.Join(a, ", ") },
		"title": func(a string) string { return strings.Title(a) },
		"comma": func(a uint) string { return humanize.Comma(int64(a)) },
		"bool": func(a bool) string {
			if a == true {
				return "Yes"
			} else {
				return "No"
			}
		},
		"artists": func(a []spotify.SimpleArtist) template.HTML {
			var artists []string
			for _, v := range a {
				artists = append(artists, "<a href=\"/artist/"+string(v.ID)+"\">"+v.Name+"</a>")
			}
			return template.HTML(strings.Join(artists, ", "))
		},
		"genres": func(a []string) string {
			var genres []string
			for _, v := range a {
				genres = append(genres, strings.Title(v))
			}
			return strings.Join(genres, ", ")
		},
		"seconds": func(inSeconds int) string {
			inSeconds = inSeconds / 1000
			minutes := inSeconds / 60
			seconds := inSeconds % 60
			return fmt.Sprintf("%vm %vs", minutes, seconds)
		},
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		//logger.ErrExit("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
