package spotify

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
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

	expiry := session.Read(r, session.TokenExpiry)

	i, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		fmt.Println("Converting expiry")
	}

	token := &oauth2.Token{
		AccessToken:  session.Read(r, session.TokenToken),
		TokenType:    session.Read(r, session.TokenType),
		RefreshToken: session.Read(r, session.TokenRefresh),
		Expiry:       time.Unix(int64(i), 0),
	}

	return GetAuthenticator(r).NewClient(token)
}

func GetOptions(r *http.Request, limit int, offset int) (opt *spotify.Options) {

	opt = &spotify.Options{}
	opt.Country = new(string)
	opt.Limit = new(int)
	opt.Offset = new(int)

	*opt.Country = session.Read(r, session.UserCountry)
	*opt.Limit = limit
	*opt.Offset = offset

	return opt
}

// Loops through pagination to get every playlist
func CurrentUsersPlaylists(r *http.Request) (playlists []spotify.SimplePlaylist, err error) {

	client := GetClient(r)

	var totalTracks = 1
	var page = 0

	for len(playlists) < totalTracks {

		options := GetOptions(r, PlayslistsLimit, page*PlayslistsLimit)
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
