package structs

import "github.com/zmb3/spotify"

type TemplateVars struct {
	LoggedIn        bool
	LoggedInID      string
	Flashes         []interface{}
	ErrorCode       string
	ErrorMessage    string
	Playlists       []spotify.SimplePlaylist
	Highlight       string
	SearchPlaylists *spotify.SearchResult
	SearchArtists   *spotify.SearchResult
	SearchAlbums    *spotify.SearchResult
	SearchTracks    *spotify.SearchResult
	Tracks          []spotify.FullTrack
	Track           *spotify.FullTrack
	Album           *spotify.FullAlbum
	Artist          *spotify.FullArtist
	Playlist        *spotify.FullPlaylist
	User            *spotify.User
}
