package oauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	loginPath         = "/auth/google/login"
	callbackPath      = "/auth/google/callback"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateTracker string
	letters           = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func generateOauthStateTracker() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	oauthStateTracker = base64.URLEncoding.EncodeToString([]byte(string(b)))
	return oauthStateTracker
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func StartLocalhostServer() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost:8000%s", callbackPath),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	r := chi.NewRouter()
	r.Get(loginPath, googleLoginHandler)
	r.Get(callbackPath, callbackHandler)

	log.Printf("Click here http://localhost:8000%s to authenticate login with your Google Account", loginPath)
	if err := http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
	log.Println("Server closed!")
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := generateOauthStateTracker()
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateTracker {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// save credentials
	fmt.Fprintf(w, "UserInfo: %s\n", data)

	// TODO: close after some time, don't leave server hanging
}
