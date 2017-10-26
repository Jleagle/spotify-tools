package structs

import "github.com/zmb3/spotify"

type TemplateVars struct {
	LoggedIn     bool
	LoggedInID   string
	Flashes      []interface{}
	ErrorCode    string
	ErrorMessage string
	Playlists    []spotify.SimplePlaylist // Shuffle
	Highlight    string
}
