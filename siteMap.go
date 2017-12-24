package main

import (
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

func siteMapHandler(w http.ResponseWriter, r *http.Request) {

	sm := stm.NewSitemap()
	sm.SetDefaultHost("https://spotifyhelper.com/")
	sm.SetCompress(true)
	sm.Create()

	sm.Add(stm.URL{"loc": "/", "changefreq": "daily", "mobile": true})
	sm.Add(stm.URL{"loc": "/shuffle", "changefreq": "daily", "mobile": true})
	sm.Add(stm.URL{"loc": "/duplicates", "changefreq": "daily", "mobile": true})
	sm.Add(stm.URL{"loc": "/random", "changefreq": "daily", "mobile": true})
	sm.Add(stm.URL{"loc": "/top", "changefreq": "daily", "mobile": true})

	w.Write(sm.XMLContent())
	return
}
