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
	Track           *spotify.FullTrack
	Tracks          []spotify.FullTrack
	Album           *spotify.FullAlbum
	Albums          *spotify.SimpleAlbumPage
	Artist          *spotify.FullArtist
	Playlist        *spotify.FullPlaylist
	User            *spotify.User
	UserPlaylists   *spotify.SimplePlaylistPage
	AudioAnalysis   *spotify.AudioAnalysis
	AudioFeatures   *spotify.AudioFeatures
	FullArtistPage  *spotify.FullArtistPage
	FullTrackPage   *spotify.FullTrackPage
}
