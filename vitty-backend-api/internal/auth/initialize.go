package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var OauthConf oauth2.Config
var JWTSecret string
var FirebaseApp *firebase.App

func InitializeAuth(jwtSecret string) {
	JWTSecret = jwtSecret
}

func InitializeGoogleOauth(googleClientId string, googleClientSecret string, googleRedirectUri string) {
	OauthConf.ClientID = googleClientId
	OauthConf.ClientSecret = googleClientSecret
	OauthConf.Scopes = []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"}
	OauthConf.Endpoint = google.Endpoint
}

func InitializeFirebaseApp() {
	opt := option.WithCredentialsFile("./credentials/firebase-creds.json")
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	FirebaseApp = firebaseApp
}
