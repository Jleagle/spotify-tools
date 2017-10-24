package spotify

import (
	"net/http"
	"os"

	"github.com/Jleagle/spotify-tools/session"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var Auth = spotify.NewAuthenticator(getDomain()+"/login-callback", getScopes()...)

func getDomain() (domain string) {
	return "http://localhost:8084"
}

func getScopes() []string {
	return []string{
		spotify.ScopeUserReadEmail,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic,
	}
}
func GetAuthenticator() (auth spotify.Authenticator) {
	Auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	return Auth
}

func GetClient(r *http.Request) (client spotify.Client) {

	token := &oauth2.Token{
		AccessToken: session.Read(r, session.Token),
	}

	return spotify.Authenticator{}.NewClient(token)
}

func GetOptions(limit int, offset int) (opt *spotify.Options) {

	opt = &spotify.Options{}

	opt.Limit = new(int)
	opt.Offset = new(int)

	*opt.Limit = limit
	*opt.Offset = offset

	return opt
}
