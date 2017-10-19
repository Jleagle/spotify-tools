package main

import (
	"net/http"

	"github.com/Jleagle/spotify-tools/session"
)

func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	//client := auth.NewClient(tok)

	session.Write(w, r, session.LastPage, "shuffle")

	returnTemplate(w, r, "shuffle", templateVars{}, nil)
	return
}
