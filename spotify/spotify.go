package spotify

import (
	"net/http"
	"os"

	"github.com/Jleagle/spotify-tools/session"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var Auth = spotify.NewAuthenticator("http://localhost:8084/login-callback", spotify.ScopeUserReadPrivate)

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
