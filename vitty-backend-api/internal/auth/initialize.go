package auth

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OauthConf oauth2.Config
var JWTSecret string

func InitializeAuth(ouathCallbackUrl string, jwtSecret string) {
	file, err := os.ReadFile("./credentials/oauth2-credentials.json")

	if err != nil {
		panic(err)
	}

	var creds interface{}
	err = json.Unmarshal(file, &creds)

	if err != nil {
		panic(err)
	}

	JWTSecret = jwtSecret
	OauthConf.ClientID = creds.(map[string]interface{})["web"].(map[string]interface{})["client_id"].(string)
	OauthConf.ClientSecret = creds.(map[string]interface{})["web"].(map[string]interface{})["client_secret"].(string)
	OauthConf.Scopes = []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"}
	OauthConf.Endpoint = google.Endpoint
}
