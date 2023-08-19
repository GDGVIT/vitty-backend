package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OauthConf oauth2.Config
var JWTSecret string

func InitializeAuth(jwtSecret string) {
	JWTSecret = jwtSecret
}

func InitializeGoogleOauth(googleClientId string, googleClientSecret string, googleRedirectUri string) {
	OauthConf.ClientID = googleClientId
	OauthConf.ClientSecret = googleClientSecret
	OauthConf.Scopes = []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"}
	OauthConf.Endpoint = google.Endpoint
}

func InitializeFirebaseAuth() {
}
