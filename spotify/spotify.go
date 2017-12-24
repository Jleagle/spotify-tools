package spotify

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Jleagle/go-helpers/rollbar"
	"github.com/Jleagle/spotifyhelper/session"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	PlayslistsLimit = 50
	TracksLimit     = 100
	MaxSearchOffset = 100000
	MaxArtistAlbums = 50
)

func GetAuthenticator(r *http.Request) (auth spotify.Authenticator) {

	scopes := []string{
		spotify.ScopeUserReadEmail,
		spotify.ScopeUserReadPrivate, // Need for country
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopeUserTopRead,
		spotify.ScopeUserReadRecentlyPlayed,
	}

	host := r.Host
	if strings.Contains(host, "8084") {
		host = "http://" + host
	} else {
		host = "https://" + host
	}

	auth = spotify.NewAuthenticator(host+"/login-callback", scopes...)
	auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	return auth
}

func GetClient(r *http.Request) (client spotify.Client) {

	expiry, err := session.Read(r, session.TokenExpiry)
	if err != nil {
		rollbar.ErrorCritical(err)
	}

	i, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		rollbar.ErrorError(err)
	}

	// Make a token
	token := &oauth2.Token{
		Expiry: time.Unix(int64(i), 0),
	}

	token.AccessToken, err = session.Read(r, session.TokenToken)
	if err != nil {
		rollbar.ErrorError(err)
	}

	token.TokenType, err = session.Read(r, session.TokenType)
	if err != nil {
		rollbar.ErrorError(err)
	}

	token.RefreshToken, err = session.Read(r, session.TokenRefresh)
	if err != nil {
		rollbar.ErrorError(err)
	}

	return GetAuthenticator(r).NewClient(token)
}

func GetOptions(r *http.Request, limit int, offset int, TimeRange string) (opt *spotify.Options) {

	opt = &spotify.Options{}
	opt.Country = new(string)
	opt.Limit = new(int)
	opt.Offset = new(int)
	opt.Timerange = new(string)

	// Get country
	country, err := session.Read(r, session.UserCountry)
	if err != nil {
		rollbar.ErrorError(err)
	}

	*opt.Country = country
	*opt.Limit = limit
	*opt.Offset = offset
	*opt.Timerange = TimeRange

	return opt
}

// Loops through pagination to get every playlist
func CurrentUsersPlaylists(r *http.Request) (playlists []spotify.SimplePlaylist, err error) {

	client := GetClient(r)

	var totalTracks = 1
	var page = 0

	for len(playlists) < totalTracks {

		options := GetOptions(r, PlayslistsLimit, page*PlayslistsLimit, "")
		response, err := client.CurrentUsersPlaylistsOpt(options)
		if err != nil {
			return playlists, err
		}
		totalTracks = response.Total
		page++

		playlists = append(playlists, response.Playlists...)
	}

	return playlists, err
}

func GetPlaylistTracks(r *http.Request) {

}
