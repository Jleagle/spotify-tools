package spotify

import (
	"net/http"
	"os"

	"github.com/Jleagle/spotify-tools/session"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func GetAuthenticator() (auth spotify.Authenticator) {

	scopes := []string{
		spotify.ScopeUserReadEmail,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic,
	}

	auth = spotify.NewAuthenticator("http://localhost:8084"+"/login-callback", scopes...)
	auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	return auth
}

func GetClient(r *http.Request) (client spotify.Client) {

	token := &oauth2.Token{
		AccessToken: session.Read(r, session.Token),
	}

	return GetAuthenticator().NewClient(token)
}

func GetOptions(limit int, offset int) (opt *spotify.Options) {

	opt = &spotify.Options{}

	opt.Limit = new(int)
	opt.Offset = new(int)

	*opt.Limit = limit
	*opt.Offset = offset

	return opt
}
